package service

import (
	"mime/multipart"

	"pos_api/domain/backup/dto"
)

type (
	BackupServiceInterface interface {
		CreateBackup() (*dto.BackupInfo, error)
		GetList() (*dto.BackupListResponse, error)
		RestoreBackup(file *multipart.FileHeader) error
	}

	backupService struct{}
)

func NewBackupService() *backupService {
	return &backupService{}
}
