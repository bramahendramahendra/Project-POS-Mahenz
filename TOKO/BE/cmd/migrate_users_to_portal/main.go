// Skrip one-off: salin user POS Toko ke Portal dan tautkan akunnya (kolom users.portal_user_id).
//
// LATAR BELAKANG (lihat docs/RANCANGAN_PORTAL_TOKO_COFFEE.md, bagian 7):
// Login dipindah ke Portal. Akun Toko yang sudah ada TIDAK dipindah (barisnya dirujuk
// sesi, transaksi, shift, dan laci kas), melainkan DISALIN ke Portal lalu ditautkan.
//
// YANG DILAKUKAN (per user Toko yang belum tertaut):
//   - Hash password disalin apa adanya (Argon2id yang sama), jadi password tidak berubah.
//   - Role Portal: role "admin" di Toko -> "admin" di Portal; SEMUA role lain -> "user".
//   - Status aktif disalin. Akun Portal baru tidak diminta ganti password.
//   - Akses aplikasi "toko" diberikan di Portal.
//   - users.portal_user_id di Toko diisi dengan id akun Portal.
//
// BILA USERNAME SUDAH ADA DI PORTAL:
//   - password Portal masih bawaan seed (admin123/owner123) -> DIAMBIL ALIH: password, nama,
//     role, dan status disamakan dengan Toko, lalu ditautkan.
//   - password Portal sama dengan Toko -> langsung ditautkan.
//   - password Portal berbeda dan bukan bawaan seed -> KONFLIK, dilewati dan dilaporkan
//     (selesaikan manual, jangan menebak: bisa jadi dua orang berbeda).
//
// AMAN DIULANG: user yang sudah tertaut dilewati; kegagalan di tengah jalan bisa dilanjutkan
// dengan menjalankan ulang.
//
// MODE:
//
//	default                          -> DRY RUN (tidak menulis apa pun; hanya laporan)
//	--apply --i-have-backup          -> benar-benar menulis
//
// BACKUP WAJIB sebelum --apply (kedua database). Contoh:
//
//	mysqldump -u USER -p pos_retail_db > backups/pos_retail_db_sebelum_tautan_portal.sql
//	mysqldump -u USER -p pos_portal_db > backups/pos_portal_db_sebelum_tautan_portal.sql
//
// KONEKSI:
//   - Toko  : sama seperti skrip cmd lain (MIGRATION_DSN, atau config sesuai RELEASE_MODE).
//   - Portal: --portal-dsn / env PORTAL_MIGRATION_DSN; bila kosong, memakai kredensial Toko
//     dengan nama database --portal-db (default pos_portal_db).
//
// Jalankan dari folder TOKO/BE:
//
//	go run ./cmd/migrate_users_to_portal                          (dry-run)
//	go run ./cmd/migrate_users_to_portal --apply --i-have-backup  (eksekusi)
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

// Hash password bawaan seed (lihat migrasi 002 Toko dan Portal): admin123 dan owner123.
// Dipakai untuk mengenali akun Portal yang belum pernah diganti passwordnya.
var seedHashes = map[string]string{
	"$argon2id$v=19$m=65536,t=3,p=4$TJxWkwEYoSuvSZ0JtwrfJw$ne9pID5QjltGeu0cNxiKYGCOXkUCRRjunQUYvttxCsM": "admin123",
	"$argon2id$v=19$m=65536,t=3,p=4$gcWNYu+y/FCRuCT6WS/+eg$1qAjwCU3HoK87s/CbDPHYo06W6B4mhoHF1SQmkigukY": "owner123",
}

const appCode = "toko"

type tokoUser struct {
	ID           int
	Username     string
	Password     string
	FullName     string
	IsActive     bool
	PortalUserID *int
	RoleName     string
}

type portalUser struct {
	ID       int
	Username string
	Password string
	FullName string
	IsActive bool
	RoleName string
}

type action string

const (
	actCreate   action = "BUAT"
	actAdopt    action = "AMBIL ALIH (akun seed)"
	actLink     action = "TAUTKAN"
	actSkip     action = "SUDAH TERTAUT"
	actConflict action = "KONFLIK"
	actOrphan   action = "TAUTAN YATIM"
)

type plan struct {
	toko       tokoUser
	act        action
	portalRole string
	portalID   int // id akun Portal bila sudah ada
	note       string
}

func portalRoleFor(tokoRole string) string {
	if tokoRole == "admin" {
		return "admin"
	}
	return "user"
}

func main() {
	apply := flag.Bool("apply", false, "benar-benar menulis (default: dry-run)")
	backupOK := flag.Bool("i-have-backup", false, "konfirmasi bahwa backup kedua database sudah dibuat (wajib bersama --apply)")
	portalDB := flag.String("portal-db", "pos_portal_db", "nama database Portal (kredensial memakai milik Toko)")
	portalDSN := flag.String("portal-dsn", os.Getenv("PORTAL_MIGRATION_DSN"), "DSN lengkap database Portal (menimpa --portal-db)")
	flag.Parse()

	if *apply && !*backupOK {
		log.Fatal("--apply membutuhkan --i-have-backup. Buat backup pos_retail_db dan pos_portal_db dulu (lihat komentar di atas file ini).")
	}

	tokoDSN, err := migrationdb.ResolveDSN()
	if err != nil {
		log.Fatalf("DSN Toko: %v", err)
	}
	pDSN := *portalDSN
	if pDSN == "" {
		cfg, err := mysql.ParseDSN(tokoDSN)
		if err != nil {
			log.Fatalf("parse DSN Toko: %v", err)
		}
		cfg.DBName = *portalDB
		pDSN = cfg.FormatDSN()
	}

	toko := openDB(tokoDSN, "Toko")
	portal := openDB(pDSN, "Portal")

	if *apply {
		fmt.Println("MODE: APPLY (menulis ke database)")
	} else {
		fmt.Println("MODE: DRY-RUN (tidak menulis apa pun; tambahkan --apply --i-have-backup untuk eksekusi)")
	}

	plans, warnings := buildPlans(toko, portal)
	printReport(plans, warnings)

	if !*apply {
		fmt.Println("\nDry-run selesai. Tidak ada perubahan.")
		return
	}

	if err := execute(toko, portal, plans); err != nil {
		log.Fatalf("GAGAL: %v\n(Portal ditulis lebih dulu lalu Toko; jalankan ulang untuk melanjutkan, skrip aman diulang)", err)
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

func buildPlans(toko, portal *gorm.DB) ([]plan, []string) {
	var warnings []string

	// Aplikasi toko harus ada di Portal.
	var appID int
	if err := portal.Raw(`SELECT id FROM apps WHERE code = ? LIMIT 1`, appCode).Scan(&appID).Error; err != nil || appID == 0 {
		log.Fatalf("aplikasi %q belum ada di Portal (tabel apps). Daftarkan dulu lewat menu Kelola Aplikasi.", appCode)
	}

	var tUsers []tokoUser
	if err := toko.Raw(`SELECT u.id, u.username, u.password, u.full_name, u.is_active, u.portal_user_id, r.name AS role_name
		FROM users u INNER JOIN roles r ON r.id = u.role_id ORDER BY u.id`).Scan(&tUsers).Error; err != nil {
		log.Fatalf("baca user Toko: %v", err)
	}

	var pUsers []portalUser
	if err := portal.Raw(`SELECT u.id, u.username, u.password, u.full_name, u.is_active, r.name AS role_name
		FROM users u INNER JOIN roles r ON r.id = u.role_id ORDER BY u.id`).Scan(&pUsers).Error; err != nil {
		log.Fatalf("baca user Portal: %v", err)
	}

	byName := make(map[string]portalUser, len(pUsers))
	byID := make(map[int]portalUser, len(pUsers))
	for _, p := range pUsers {
		byName[strings.ToLower(p.Username)] = p
		byID[p.ID] = p
	}

	// Akun Portal yang sudah dipakai tautan Toko (untuk mendeteksi satu akun Portal diklaim dua kali).
	claimed := make(map[int]int) // portalID -> tokoID
	for _, t := range tUsers {
		if t.PortalUserID != nil {
			claimed[*t.PortalUserID] = t.ID
		}
	}

	var plans []plan
	for _, t := range tUsers {
		p := plan{toko: t, portalRole: portalRoleFor(t.RoleName)}

		switch {
		case t.PortalUserID != nil:
			if _, ok := byID[*t.PortalUserID]; ok {
				p.act, p.portalID = actSkip, *t.PortalUserID
			} else {
				p.act = actOrphan
				p.note = fmt.Sprintf("portal_user_id=%d tidak ada di Portal; periksa manual", *t.PortalUserID)
			}

		default:
			existing, found := byName[strings.ToLower(t.Username)]
			switch {
			case !found:
				p.act = actCreate

			case claimed[existing.ID] != 0 && claimed[existing.ID] != t.ID:
				p.act = actConflict
				p.note = fmt.Sprintf("akun Portal id=%d sudah ditautkan ke user Toko id=%d", existing.ID, claimed[existing.ID])

			case existing.Password == t.Password:
				p.act, p.portalID = actLink, existing.ID
				p.note = "password sama"

			case seedHashes[existing.Password] != "":
				p.act, p.portalID = actAdopt, existing.ID
				p.note = fmt.Sprintf("akun Portal masih memakai password bawaan (%s)", seedHashes[existing.Password])

			default:
				p.act = actConflict
				p.note = fmt.Sprintf("username sudah ada di Portal (id=%d) dengan password berbeda; selesaikan manual", existing.ID)
			}
		}
		plans = append(plans, p)
	}

	// Peringatan: akun Portal yang masih berpassword bawaan setelah migrasi.
	linkedPortalIDs := make(map[int]bool)
	for _, p := range plans {
		if p.portalID != 0 && (p.act == actLink || p.act == actAdopt || p.act == actSkip) {
			linkedPortalIDs[p.portalID] = true
		}
	}
	for _, pu := range pUsers {
		if seedHashes[pu.Password] != "" && !linkedPortalIDs[pu.ID] {
			warnings = append(warnings, fmt.Sprintf(
				"akun Portal %q (id=%d, role %s) masih memakai password bawaan %q dan tidak ditimpa migrasi. Ganti passwordnya atau nonaktifkan/hapus lewat menu Kelola User",
				pu.Username, pu.ID, pu.RoleName, seedHashes[pu.Password]))
		}
	}

	return plans, warnings
}

func printReport(plans []plan, warnings []string) {
	counts := map[action]int{}
	fmt.Printf("\n%-5s %-20s %-10s %-10s %-24s %s\n", "ID", "USERNAME", "ROLE TOKO", "-> PORTAL", "AKSI", "CATATAN")
	for _, p := range plans {
		counts[p.act]++
		fmt.Printf("%-5d %-20s %-10s %-10s %-24s %s\n", p.toko.ID, p.toko.Username, p.toko.RoleName, p.portalRole, p.act, p.note)
	}

	fmt.Println("\nRINGKASAN:")
	for _, a := range []action{actCreate, actAdopt, actLink, actSkip, actConflict, actOrphan} {
		fmt.Printf("  %-26s %d\n", a, counts[a])
	}

	if len(warnings) > 0 {
		fmt.Println("\nPERINGATAN:")
		for _, w := range warnings {
			fmt.Println("  - " + w)
		}
	}
}

func execute(toko, portal *gorm.DB, plans []plan) error {
	var appID int
	if err := portal.Raw(`SELECT id FROM apps WHERE code = ? LIMIT 1`, appCode).Scan(&appID).Error; err != nil || appID == 0 {
		return fmt.Errorf("aplikasi %q tidak ada di Portal", appCode)
	}

	roleIDs := map[string]int{}
	for _, name := range []string{"admin", "user"} {
		var id int
		if err := portal.Raw(`SELECT id FROM roles WHERE name = ? LIMIT 1`, name).Scan(&id).Error; err != nil || id == 0 {
			return fmt.Errorf("role %q tidak ada di Portal", name)
		}
		roleIDs[name] = id
	}

	// 1) Portal: buat / ambil alih akun + akses aplikasi, semuanya dalam satu transaksi.
	linkTargets := map[int]int{} // tokoID -> portalID
	err := portal.Transaction(func(tx *gorm.DB) error {
		for _, p := range plans {
			switch p.act {
			case actCreate:
				if err := tx.Exec(`INSERT INTO users (username, password, full_name, role_id, is_active, must_change_password)
					VALUES (?, ?, ?, ?, ?, 0)`, p.toko.Username, p.toko.Password, p.toko.FullName, roleIDs[p.portalRole], p.toko.IsActive).Error; err != nil {
					return fmt.Errorf("buat %s: %w", p.toko.Username, err)
				}
				var id int
				if err := tx.Raw(`SELECT id FROM users WHERE username = ? LIMIT 1`, p.toko.Username).Scan(&id).Error; err != nil || id == 0 {
					return fmt.Errorf("ambil id %s setelah dibuat: %v", p.toko.Username, err)
				}
				linkTargets[p.toko.ID] = id

			case actAdopt:
				if err := tx.Exec(`UPDATE users SET password = ?, full_name = ?, role_id = ?, is_active = ?, must_change_password = 0 WHERE id = ?`,
					p.toko.Password, p.toko.FullName, roleIDs[p.portalRole], p.toko.IsActive, p.portalID).Error; err != nil {
					return fmt.Errorf("ambil alih %s: %w", p.toko.Username, err)
				}
				// Sesi lama memuat role/password lama; paksa login ulang.
				if err := tx.Exec(`DELETE FROM sessions WHERE user_id = ?`, p.portalID).Error; err != nil {
					return err
				}
				linkTargets[p.toko.ID] = p.portalID

			case actLink:
				linkTargets[p.toko.ID] = p.portalID
			}
		}

		// Akses aplikasi toko untuk semua akun yang tertaut/dibuat.
		for _, portalID := range linkTargets {
			if err := tx.Exec(`INSERT IGNORE INTO user_apps (user_id, app_id) VALUES (?, ?)`, portalID, appID).Error; err != nil {
				return fmt.Errorf("beri akses aplikasi: %w", err)
			}
		}
		// User yang sudah tertaut sebelumnya juga dipastikan punya akses.
		for _, p := range plans {
			if p.act == actSkip {
				if err := tx.Exec(`INSERT IGNORE INTO user_apps (user_id, app_id) VALUES (?, ?)`, p.portalID, appID).Error; err != nil {
					return fmt.Errorf("beri akses aplikasi: %w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaksi Portal dibatalkan: %w", err)
	}

	// 2) Toko: isi portal_user_id.
	err = toko.Transaction(func(tx *gorm.DB) error {
		for tokoID, portalID := range linkTargets {
			if err := tx.Exec(`UPDATE users SET portal_user_id = ? WHERE id = ? AND portal_user_id IS NULL`, portalID, tokoID).Error; err != nil {
				return fmt.Errorf("tautkan user Toko id=%d: %w", tokoID, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaksi Toko dibatalkan: %w", err)
	}

	fmt.Printf("Ditulis: %d akun ditautkan.\n", len(linkTargets))
	return nil
}
