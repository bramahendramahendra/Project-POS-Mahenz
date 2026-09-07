// Package migrationdb menyediakan resolusi koneksi database untuk skrip-skrip
// backfill/migrasi one-off di BE/cmd/*, TANPA menarik package "pos_api/config".
//
// Kenapa tidak pakai package config biasa? Package config punya init() yang
// memuat SELURUH konfigurasi aplikasi (termasuk JWT) dan akan panic kalau
// environment variable SECRETKEY tidak diset. Skrip backfill hanya butuh
// kredensial DATABASE, tidak butuh SecretKey — jadi memaksa set SECRETKEY cuma
// untuk menjalankan backfill itu merepotkan dan tidak relevan. Helper ini
// membaca hanya blok "Database" dari file config, tanpa validasi lain.
package migrationdb

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// configPathByMode memetakan RELEASE_MODE ke file config-nya, sama seperti
// mapping di package config utama.
var configPathByMode = map[string]string{
	"dev":   "./config/config_dev.json",
	"prod":  "./config/config_prod.json",
	"uat":   "./config/config_uat.json",
	"local": "./config/config_local.json",
	"bors":  "./config/config_bors_kost.json",
}

// ResolveDSN mengembalikan DSN MySQL untuk skrip migrasi, dengan urutan
// prioritas:
//
//  1. Env MIGRATION_DSN — override manual penuh. Berguna untuk menunjuk
//     database lain (mis. pos_retail_db_20260907 hasil migrasi) atau
//     kredensial berbeda tanpa menyentuh file config sama sekali. Contoh:
//     export MIGRATION_DSN='pos_user:pass@tcp(127.0.0.1:3306)/nama_db?charset=utf8&parseTime=True&loc=Local'
//
//  2. Kalau MIGRATION_DSN kosong: baca RELEASE_MODE dari ./.env, lalu ambil
//     HANYA blok Database dari ./config/config_<mode>.json. Kredensial yang
//     dipakai sama persis dengan backend, tapi tanpa memicu validasi SecretKey.
//
// Karena membaca ./.env dan ./config/*.json (path relatif), skrip pemanggil
// harus dijalankan dari folder BE, mis: go run ./cmd/backfill_stock_restore
func ResolveDSN() (string, error) {
	if v := os.Getenv("MIGRATION_DSN"); v != "" {
		return v, nil
	}

	mode, err := releaseMode()
	if err != nil {
		return "", err
	}

	cfgPath, ok := configPathByMode[mode]
	if !ok {
		return "", fmt.Errorf("RELEASE_MODE %q tidak dikenal (pilihan: dev/prod/uat/local/bors), atau set MIGRATION_DSN manual", mode)
	}

	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return "", fmt.Errorf("gagal baca config %s: %w (jalankan dari folder BE, atau set MIGRATION_DSN)", cfgPath, err)
	}

	user := v.GetString("Database.User")
	pass := v.GetString("Database.Password")
	host := v.GetString("Database.Host")
	port := v.GetString("Database.Port")
	name := v.GetString("Database.Database")

	if user == "" || host == "" || port == "" || name == "" {
		return "", fmt.Errorf("blok Database di %s tidak lengkap (User/Host/Port/Database wajib terisi)", cfgPath)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, pass, host, port, name)
	return dsn, nil
}

// releaseMode membaca RELEASE_MODE dari ./.env (mengikuti cara package config
// utama). Env variable RELEASE_MODE yang sudah diset di shell diprioritaskan.
func releaseMode() (string, error) {
	if v := os.Getenv("RELEASE_MODE"); v != "" {
		return v, nil
	}

	v := viper.New()
	v.SetConfigFile("./.env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return "", fmt.Errorf("gagal baca ./.env: %w (jalankan dari folder BE, set RELEASE_MODE, atau set MIGRATION_DSN)", err)
	}
	mode := v.GetString("RELEASE_MODE")
	if mode == "" {
		return "", fmt.Errorf("RELEASE_MODE tidak ditemukan di ./.env (atau set MIGRATION_DSN manual)")
	}
	return mode, nil
}
