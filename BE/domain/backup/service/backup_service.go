package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pos_api/config"
	"pos_api/domain/backup/dto"
	"pos_api/errors"
)

const backupDir = "backups"

// dbConnArgs membangun argumen koneksi mysql/mysqldump. Flag -p SENGAJA dihilangkan
// sepenuhnya saat password kosong — "-p" tanpa nilai langsung setelahnya akan dibaca
// client mysql sebagai "prompt password secara interaktif" dan membuat proses hang
// menunggu stdin selamanya, bukan berarti "password kosong".
func dbConnArgs() []string {
	args := []string{
		"-h", config.Db.Host,
		"-P", config.Db.Port,
		"-u", config.Db.User,
	}
	if config.Db.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", config.Db.Password))
	}
	return append(args, config.Db.Database)
}

func (s *backupService) CreateBackup() (*dto.BackupInfo, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal membuat folder backup"}
	}

	filename := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(backupDir, filename)

	cmd := exec.Command("mysqldump", dbConnArgs()...)

	outFile, err := os.Create(backupPath)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal membuat file backup"}
	}

	cmd.Stdout = outFile
	runErr := cmd.Run()
	outFile.Close()

	if runErr != nil {
		// File harus ditutup dulu sebelum dihapus — di Windows, os.Remove akan gagal
		// (tanpa error yang terlihat di sini) selama file masih dalam keadaan terbuka.
		os.Remove(backupPath)
		return nil, &errors.InternalServerError{Message: "Gagal menjalankan mysqldump: " + runErr.Error()}
	}

	info, err := os.Stat(backupPath)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal membaca info file backup"}
	}

	return &dto.BackupInfo{
		Filename:  filename,
		Size:      formatFileSize(info.Size()),
		CreatedAt: info.ModTime(),
	}, nil
}

func (s *backupService) GetList() (*dto.BackupListResponse, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal membaca folder backup"}
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal membaca daftar backup"}
	}

	var files []dto.BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, dto.BackupInfo{
			Filename:  entry.Name(),
			Size:      formatFileSize(info.Size()),
			CreatedAt: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].CreatedAt.After(files[j].CreatedAt)
	})

	return &dto.BackupListResponse{Files: files}, nil
}

func (s *backupService) RestoreBackup(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return &errors.BadRequestError{Message: "Gagal membuka file restore"}
	}
	defer src.Close()

	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("restore_%d.sql", time.Now().UnixNano()))
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal membuat file sementara"}
	}

	buf := make([]byte, 32*1024)
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			if _, writeErr := tmpFile.Write(buf[:n]); writeErr != nil {
				tmpFile.Close()
				os.Remove(tmpPath)
				return &errors.InternalServerError{Message: "Gagal menyimpan file sementara"}
			}
		}
		if readErr != nil {
			break
		}
	}
	tmpFile.Close()
	defer os.Remove(tmpPath)

	sqlFile, err := os.Open(tmpPath)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal membaca file restore"}
	}
	defer sqlFile.Close()

	cmd := exec.Command("mysql", dbConnArgs()...)
	cmd.Stdin = sqlFile

	if err := cmd.Run(); err != nil {
		return &errors.InternalServerError{Message: "Gagal melakukan restore: " + err.Error()}
	}

	return nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
