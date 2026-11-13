// Perintah server: buat kunci baru untuk sebuah aplikasi di BE User (X-App-Key).
//
// Dipakai saat pemasangan (kunci awal portal/toko/coffee) dan untuk mengganti kunci Portal,
// yang sengaja tidak bisa diputar lewat UI Portal (Portal akan langsung terputus).
//
// Kunci asli hanya dicetak SEKALI ke layar; yang disimpan di pos_auth_db hanya hash SHA-256.
// Setelah dijalankan, kunci lama aplikasi itu langsung tidak berlaku: pasang kunci baru di
// database aplikasinya (tabel auth_credentials) secepatnya.
//
// MODE:
//
//	go run ./cmd/set_app_key --app toko              -> hanya menampilkan status (tidak menulis)
//	go run ./cmd/set_app_key --app toko --confirm    -> membuat kunci baru
//
// Jalankan dari folder AUTH/BE. Koneksi: env MIGRATION_DSN, atau config sesuai RELEASE_MODE.
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"auth_api/cmd/internal/migrationdb"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type app struct {
	ID       int
	Code     string
	Name     string
	IsActive bool
	HasKey   bool
}

func main() {
	appCode := flag.String("app", "", "kode aplikasi (mis. portal, toko, coffee)")
	confirm := flag.Bool("confirm", false, "benar-benar membuat kunci baru (kunci lama langsung tidak berlaku)")
	flag.Parse()

	if *appCode == "" {
		log.Fatal("--app wajib diisi (mis. --app toko)")
	}

	dsn, err := migrationdb.ResolveDSN()
	if err != nil {
		log.Fatalf("DSN: %v", err)
	}
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("koneksi database: %v", err)
	}

	var a app
	if err := db.Raw(`SELECT id, code, name, is_active, (app_key_hash IS NOT NULL AND app_key_hash <> '') AS has_key FROM apps WHERE code = ? LIMIT 1`, *appCode).Scan(&a).Error; err != nil {
		log.Fatalf("baca aplikasi: %v", err)
	}
	if a.ID == 0 {
		log.Fatalf("aplikasi %q tidak ada di BE User", *appCode)
	}

	fmt.Printf("Aplikasi: %s (%s), aktif=%v, sudah punya kunci=%v\n", a.Name, a.Code, a.IsActive, a.HasKey)
	if !*confirm {
		fmt.Println("Tidak ada perubahan. Tambahkan --confirm untuk membuat kunci baru (kunci lama langsung tidak berlaku).")
		return
	}

	key, err := generateAppKey()
	if err != nil {
		log.Fatalf("buat kunci: %v", err)
	}
	sum := sha256.Sum256([]byte(key))
	if err := db.Exec(`UPDATE apps SET app_key_hash = ?, updated_at = NOW() WHERE id = ?`, hex.EncodeToString(sum[:]), a.ID).Error; err != nil {
		log.Fatalf("simpan kunci: %v", err)
	}

	fmt.Println("\nKunci baru (TAMPIL SEKALI, simpan sekarang):")
	fmt.Println(key)
	fmt.Println("\nPasang di database aplikasi " + a.Code + ":")
	fmt.Printf("  INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, '%s', '<KUNCI>')\n", a.Code)
	fmt.Println("  ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key);")
}

// generateAppKey: 32 byte acak dalam hex, berawalan "pak_" (sama dengan kunci dari UI Portal).
func generateAppKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "pak_" + hex.EncodeToString(b), nil
}
