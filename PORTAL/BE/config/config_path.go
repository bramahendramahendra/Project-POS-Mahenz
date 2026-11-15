package config

import (
	"os"
	"path/filepath"
)

// resolvePath mengembalikan p apa adanya bila file itu ada di folder kerja (server, go run dari folder BE:
// perilaku tidak berubah). Bila tidak ada (mis. go test berjalan di folder paket seperti helper/), file dicari
// relatif terhadap folder induk terdekat yang memuat .env, agar uji paket tidak panik saat config dimuat.
func resolvePath(p string) string {
	if _, err := os.Stat(p); err == nil {
		return p
	}
	if base, ok := findBaseDir(); ok {
		return filepath.Join(base, p)
	}
	return p
}

// findBaseDir mencari folder (folder kerja lalu induknya, paling jauh 6 tingkat) yang memuat .env.
func findBaseDir() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
