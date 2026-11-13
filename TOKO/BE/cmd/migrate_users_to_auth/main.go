// Skrip one-off: salin user POS Toko ke BE User (layanan akun, pos_auth_db) dan tautkan
// profilnya (kolom users.auth_user_id). Pengganti cmd/migrate_users_to_portal (Fase 3B).
//
// LATAR BELAKANG (docs/RANCANGAN_PORTAL_TOKO_COFFEE.md, bagian 5.0d dan keputusan 63-77):
// Akun (username, password, aktif, akses aplikasi) dikelola BE User. Profil Toko yang sudah ada
// TIDAK dipindah (barisnya dirujuk sesi, transaksi, shift, laci kas), melainkan akunnya DISALIN
// ke BE User lalu ditautkan, supaya kasir tetap login dengan password lama.
//
// YANG DILAKUKAN (per profil Toko yang belum tertaut):
//   - Hash password disalin apa adanya (Argon2id yang sama), jadi password tidak berubah.
//   - Role di BE User: role "admin" di Toko -> "admin"; SEMUA role lain -> "user".
//     (Role Toko sendiri tidak berubah: owner tetap owner, kasir tetap kasir.)
//   - Status aktif disalin. Akun baru tidak diminta ganti password.
//   - Akses aplikasi "toko" diberikan di BE User.
//   - users.auth_user_id di Toko diisi dengan id akun BE User.
//
// BILA USERNAME SUDAH ADA DI BE USER:
//   - password BE User masih bawaan seed (admin123/owner123) -> DIAMBIL ALIH: password, nama,
//     role, dan status disamakan dengan Toko, lalu ditautkan.
//   - password BE User sama dengan Toko -> langsung ditautkan.
//   - password berbeda dan bukan bawaan seed -> KONFLIK, dilewati dan dilaporkan. Tidak menebak:
//     bisa jadi dua orang berbeda. Bila memang orang yang sama, cukup login sekali ke Toko memakai
//     password BE User: profil ditautkan otomatis (keputusan 67).
//
// Profil Toko tanpa password lokal (dibuat otomatis saat login lewat BE User) selalu sudah
// tertaut; bila belum, dilaporkan TANPA PASSWORD dan dilewati.
//
// AMAN DIULANG: profil yang sudah tertaut dilewati; kegagalan di tengah jalan bisa dilanjutkan
// dengan menjalankan ulang.
//
// MODE:
//
//	default                          -> DRY RUN (tidak menulis apa pun; hanya laporan)
//	--apply --i-have-backup          -> benar-benar menulis
//
// BACKUP WAJIB sebelum --apply (kedua database). Contoh:
//
//	mysqldump -u USER -p pos_retail_db > backups/pos_retail_db_sebelum_salin_akun.sql
//	mysqldump -u USER -p pos_auth_db   > backups/pos_auth_db_sebelum_salin_akun.sql
//
// KONEKSI:
//   - Toko   : sama seperti skrip cmd lain (MIGRATION_DSN, atau config sesuai RELEASE_MODE).
//   - BE User: --auth-dsn / env AUTH_MIGRATION_DSN; bila kosong, memakai kredensial Toko
//     dengan nama database --auth-db (default pos_auth_db).
//
// Jalankan dari folder TOKO/BE:
//
//	go run ./cmd/migrate_users_to_auth                          (dry-run)
//	go run ./cmd/migrate_users_to_auth --apply --i-have-backup  (eksekusi)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"pos_api/cmd/internal/migrationdb"

	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Hash password bawaan seed (migrasi 002 Toko dan BE User): admin123 dan owner123.
// Dipakai untuk mengenali akun BE User yang belum pernah diganti passwordnya.
var seedHashes = map[string]string{
	"$argon2id$v=19$m=65536,t=3,p=4$TJxWkwEYoSuvSZ0JtwrfJw$ne9pID5QjltGeu0cNxiKYGCOXkUCRRjunQUYvttxCsM": "admin123",
	"$argon2id$v=19$m=65536,t=3,p=4$gcWNYu+y/FCRuCT6WS/+eg$1qAjwCU3HoK87s/CbDPHYo06W6B4mhoHF1SQmkigukY": "owner123",
}

const appCode = "toko"

type tokoUser struct {
	ID         int
	Username   string
	Password   *string
	FullName   string
	IsActive   bool
	AuthUserID *int
	RoleName   string
}

type authUser struct {
	ID       int
	Username string
	Password string
	FullName string
	IsActive bool
	RoleName string
}

type action string

const (
	actCreate     action = "BUAT"
	actAdopt      action = "AMBIL ALIH (akun seed)"
	actLink       action = "TAUTKAN"
	actSkip       action = "SUDAH TERTAUT"
	actConflict   action = "KONFLIK"
	actOrphan     action = "TAUTAN YATIM"
	actNoPassword action = "TANPA PASSWORD"
)

type plan struct {
	toko     tokoUser
	act      action
	authRole string
	authID   int // id akun BE User bila sudah ada
	note     string
}

func authRoleFor(tokoRole string) string {
	if tokoRole == "admin" {
		return "admin"
	}
	return "user"
}

func main() {
	apply := flag.Bool("apply", false, "benar-benar menulis (default: dry-run)")
	backupOK := flag.Bool("i-have-backup", false, "konfirmasi bahwa backup kedua database sudah dibuat (wajib bersama --apply)")
	authDB := flag.String("auth-db", "pos_auth_db", "nama database BE User (kredensial memakai milik Toko)")
	authDSN := flag.String("auth-dsn", os.Getenv("AUTH_MIGRATION_DSN"), "DSN lengkap database BE User (menimpa --auth-db)")
	flag.Parse()

	if *apply && !*backupOK {
		log.Fatal("--apply membutuhkan --i-have-backup. Buat backup pos_retail_db dan pos_auth_db dulu (lihat komentar di atas file ini).")
	}

	tokoDSN, err := migrationdb.ResolveDSN()
	if err != nil {
		log.Fatalf("DSN Toko: %v", err)
	}
	aDSN := *authDSN
	if aDSN == "" {
		cfg, err := mysql.ParseDSN(tokoDSN)
		if err != nil {
			log.Fatalf("parse DSN Toko: %v", err)
		}
		cfg.DBName = *authDB
		aDSN = cfg.FormatDSN()
	}

	toko := openDB(tokoDSN, "Toko")
	auth := openDB(aDSN, "BE User")

	if *apply {
		fmt.Println("MODE: APPLY (menulis ke database)")
	} else {
		fmt.Println("MODE: DRY-RUN (tidak menulis apa pun; tambahkan --apply --i-have-backup untuk eksekusi)")
	}

	plans, warnings := buildPlans(toko, auth)
	printReport(plans, warnings)

	if !*apply {
		fmt.Println("\nDry-run selesai. Tidak ada perubahan.")
		return
	}

	if err := execute(toko, auth, plans); err != nil {
		log.Fatalf("GAGAL: %v\n(BE User ditulis lebih dulu lalu Toko; jalankan ulang untuk melanjutkan, skrip aman diulang)", err)
	}
	fmt.Println("\nSelesai.")
}

func openDB(dsn, label string) *gorm.DB {
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("koneksi %s: %v", label, err)
	}
	return db
}

func buildPlans(toko, auth *gorm.DB) ([]plan, []string) {
	var warnings []string

	var appID int
	if err := auth.Raw(`SELECT id FROM apps WHERE code = ? LIMIT 1`, appCode).Scan(&appID).Error; err != nil || appID == 0 {
		log.Fatalf("aplikasi %q belum ada di BE User (tabel apps).", appCode)
	}

	var tUsers []tokoUser
	if err := toko.Raw(`SELECT u.id, u.username, u.password, u.full_name, u.is_active, u.auth_user_id, r.name AS role_name
		FROM users u INNER JOIN roles r ON r.id = u.role_id ORDER BY u.id`).Scan(&tUsers).Error; err != nil {
		log.Fatalf("baca user Toko: %v", err)
	}

	var aUsers []authUser
	if err := auth.Raw(`SELECT u.id, u.username, u.password, u.full_name, u.is_active, r.name AS role_name
		FROM users u INNER JOIN roles r ON r.id = u.role_id ORDER BY u.id`).Scan(&aUsers).Error; err != nil {
		log.Fatalf("baca user BE User: %v", err)
	}

	byName := make(map[string]authUser, len(aUsers))
	byID := make(map[int]authUser, len(aUsers))
	for _, a := range aUsers {
		byName[strings.ToLower(a.Username)] = a
		byID[a.ID] = a
	}

	// Akun BE User yang sudah dipakai tautan Toko (mendeteksi satu akun diklaim dua profil).
	claimed := make(map[int]int) // authID -> tokoID
	for _, t := range tUsers {
		if t.AuthUserID != nil {
			claimed[*t.AuthUserID] = t.ID
		}
	}

	var plans []plan
	for _, t := range tUsers {
		p := plan{toko: t, authRole: authRoleFor(t.RoleName)}
		hash := ""
		if t.Password != nil {
			hash = *t.Password
		}

		switch {
		case t.AuthUserID != nil:
			if _, ok := byID[*t.AuthUserID]; ok {
				p.act, p.authID = actSkip, *t.AuthUserID
			} else {
				p.act = actOrphan
				p.note = fmt.Sprintf("auth_user_id=%d tidak ada di BE User; periksa manual", *t.AuthUserID)
			}

		case hash == "":
			p.act = actNoPassword
			p.note = "profil tanpa password lokal; tertaut otomatis saat orangnya login lewat BE User"

		default:
			existing, found := byName[strings.ToLower(t.Username)]
			switch {
			case !found:
				p.act = actCreate

			case claimed[existing.ID] != 0 && claimed[existing.ID] != t.ID:
				p.act = actConflict
				p.note = fmt.Sprintf("akun BE User id=%d sudah ditautkan ke profil Toko id=%d", existing.ID, claimed[existing.ID])

			case existing.Password == hash:
				p.act, p.authID = actLink, existing.ID
				p.note = "password sama"

			case seedHashes[existing.Password] != "":
				p.act, p.authID = actAdopt, existing.ID
				p.note = fmt.Sprintf("akun BE User masih memakai password bawaan (%s)", seedHashes[existing.Password])

			default:
				p.act = actConflict
				p.note = fmt.Sprintf("username sudah ada di BE User (id=%d) dengan password berbeda; bila orang yang sama, cukup login ke Toko dengan password BE User", existing.ID)
			}
		}
		plans = append(plans, p)
	}

	// Peringatan: akun BE User yang masih berpassword bawaan setelah migrasi.
	linked := make(map[int]bool)
	for _, p := range plans {
		if p.authID != 0 && (p.act == actLink || p.act == actAdopt || p.act == actSkip) {
			linked[p.authID] = true
		}
	}
	for _, a := range aUsers {
		if seedHashes[a.Password] != "" && !linked[a.ID] {
			warnings = append(warnings, fmt.Sprintf(
				"akun BE User %q (id=%d, role %s) masih memakai password bawaan %q dan tidak ditimpa migrasi. Ganti passwordnya atau nonaktifkan/hapus lewat Portal",
				a.Username, a.ID, a.RoleName, seedHashes[a.Password]))
		}
	}

	return plans, warnings
}

func printReport(plans []plan, warnings []string) {
	counts := map[action]int{}
	fmt.Printf("\n%-5s %-20s %-15s %-9s %-24s %s\n", "ID", "USERNAME", "ROLE TOKO", "-> AKUN", "AKSI", "CATATAN")
	for _, p := range plans {
		counts[p.act]++
		fmt.Printf("%-5d %-20s %-15s %-9s %-24s %s\n", p.toko.ID, p.toko.Username, p.toko.RoleName, p.authRole, p.act, p.note)
	}

	fmt.Println("\nRINGKASAN:")
	for _, a := range []action{actCreate, actAdopt, actLink, actSkip, actConflict, actOrphan, actNoPassword} {
		fmt.Printf("  %-26s %d\n", a, counts[a])
	}

	if len(warnings) > 0 {
		fmt.Println("\nPERINGATAN:")
		for _, w := range warnings {
			fmt.Println("  - " + w)
		}
	}
}

func execute(toko, auth *gorm.DB, plans []plan) error {
	var appID int
	if err := auth.Raw(`SELECT id FROM apps WHERE code = ? LIMIT 1`, appCode).Scan(&appID).Error; err != nil || appID == 0 {
		return fmt.Errorf("aplikasi %q tidak ada di BE User", appCode)
	}

	roleIDs := map[string]int{}
	for _, name := range []string{"admin", "user"} {
		var id int
		if err := auth.Raw(`SELECT id FROM roles WHERE name = ? LIMIT 1`, name).Scan(&id).Error; err != nil || id == 0 {
			return fmt.Errorf("role %q tidak ada di BE User", name)
		}
		roleIDs[name] = id
	}

	// 1) BE User: buat / ambil alih akun + akses aplikasi, semuanya dalam satu transaksi.
	linkTargets := map[int]int{} // tokoID -> authID
	err := auth.Transaction(func(tx *gorm.DB) error {
		for _, p := range plans {
			switch p.act {
			case actCreate:
				if err := tx.Exec(`INSERT INTO users (username, password, full_name, role_id, is_active, must_change_password)
					VALUES (?, ?, ?, ?, ?, 0)`, p.toko.Username, *p.toko.Password, p.toko.FullName, roleIDs[p.authRole], p.toko.IsActive).Error; err != nil {
					return fmt.Errorf("buat %s: %w", p.toko.Username, err)
				}
				var id int
				if err := tx.Raw(`SELECT id FROM users WHERE username = ? LIMIT 1`, p.toko.Username).Scan(&id).Error; err != nil || id == 0 {
					return fmt.Errorf("ambil id %s setelah dibuat: %v", p.toko.Username, err)
				}
				linkTargets[p.toko.ID] = id

			case actAdopt:
				// BE User tidak menyimpan sesi; sesi Portal lama akun ini ikut ditolak karena
				// Portal memeriksa akun ke BE User saat /auth/me dan refresh.
				if err := tx.Exec(`UPDATE users SET password = ?, full_name = ?, role_id = ?, is_active = ?, must_change_password = 0 WHERE id = ?`,
					*p.toko.Password, p.toko.FullName, roleIDs[p.authRole], p.toko.IsActive, p.authID).Error; err != nil {
					return fmt.Errorf("ambil alih %s: %w", p.toko.Username, err)
				}
				linkTargets[p.toko.ID] = p.authID

			case actLink:
				linkTargets[p.toko.ID] = p.authID
			}
		}

		// Akses aplikasi toko untuk semua akun yang tertaut/dibuat, termasuk yang sudah tertaut.
		ids := make([]int, 0, len(linkTargets))
		for _, authID := range linkTargets {
			ids = append(ids, authID)
		}
		for _, p := range plans {
			if p.act == actSkip {
				ids = append(ids, p.authID)
			}
		}
		for _, authID := range ids {
			if err := tx.Exec(`INSERT IGNORE INTO user_apps (user_id, app_id) VALUES (?, ?)`, authID, appID).Error; err != nil {
				return fmt.Errorf("beri akses aplikasi: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaksi BE User dibatalkan: %w", err)
	}

	// 2) Toko: isi auth_user_id.
	err = toko.Transaction(func(tx *gorm.DB) error {
		for tokoID, authID := range linkTargets {
			if err := tx.Exec(`UPDATE users SET auth_user_id = ? WHERE id = ? AND auth_user_id IS NULL`, authID, tokoID).Error; err != nil {
				return fmt.Errorf("tautkan profil Toko id=%d: %w", tokoID, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaksi Toko dibatalkan: %w", err)
	}

	fmt.Printf("Ditulis: %d profil ditautkan.\n", len(linkTargets))
	return nil
}
