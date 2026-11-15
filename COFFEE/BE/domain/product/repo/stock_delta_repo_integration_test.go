package repo

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	model_product "coffee_api/domain/product/model"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Skenario 4 (docs/RENCANA_PERBAIKAN_STOK_PRESISI.md, Fase 2): locking
// konsisten saat 2+ operasi berjalan bersamaan pada produk yang sama. Ini
// TIDAK bisa diuji lewat unit test murni (ComputeStockDelta tidak menyentuh
// DB sama sekali) -- butuh koneksi DB nyata untuk membuktikan FOR UPDATE
// benar-benar mencegah lost update. Test ini skip otomatis kalau DB dev
// tidak bisa dihubungi (mis. CI tanpa MySQL), tidak menggagalkan `go test`.
//
// Coffee: memakai DB dev pos_coffee_db (salinan Toko sebelumnya menunjuk pos_retail_db). Paket ini
// mengimpor config yang membaca ./.env saat init, jadi jalankan dari folder COFFEE/BE, mis.:
//
//	go test -c -o ../../testing/x/product_repo.test.exe ./domain/product/repo
//	../../testing/x/product_repo.test.exe -test.v
func TestApplyStockDelta_ConcurrentSales_NoLostUpdate(t *testing.T) {
	db := connectDevDBOrSkip(t)

	productID := setupConcurrencyTestProduct(t, db, 20) // 20 Botol lepas
	defer teardownConcurrencyTestProduct(t, db, productID)

	const workers = 20
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := db.Transaction(func(tx *gorm.DB) error {
				_, err := ApplyStockDelta(tx, ApplyStockDeltaParams{
					ProductID:     productID,
					PackageID:     botolPackageIDForTest,
					Quantity:      1,
					Direction:     model_product.StockOut,
					MutationType:  "out",
					ReferenceType: "test_concurrency",
					ReferenceID:   0,
					Notes:         "concurrency test",
				})
				return err
			})
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("salah satu goroutine gagal (harusnya semua berhasil, stok cukup untuk 20x jual 1): %v", err)
		}
	}

	var finalBotolStock float64
	if err := db.Raw("SELECT stock FROM product_packages WHERE id = ?", botolPackageIDForTest).Scan(&finalBotolStock).Error; err != nil {
		t.Fatalf("gagal baca stok akhir: %v", err)
	}
	if finalBotolStock != 0 {
		t.Fatalf("kalau locking bocor (lost update), stok akhir tidak akan tepat 0 -- dapat %v setelah 20x jual 1 Botol dari stok awal 20", finalBotolStock)
	}
}

// botolPackageIDForTest diisi oleh setupConcurrencyTestProduct.
var botolPackageIDForTest int

func connectDevDBOrSkip(t *testing.T) *gorm.DB {
	t.Helper()
	// Uji ini menulis ke database dev: hanya jalan bila diminta jelas (POS_DB_TEST=1), agar `go test ./...` aman.
	if os.Getenv("POS_DB_TEST") != "1" {
		t.Skip("skip: uji integrasi DB; jalankan dengan POS_DB_TEST=1")
	}
	dsn := "root:@tcp(127.0.0.1:3306)/pos_coffee_db?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Skipf("skip: DB dev tidak bisa dihubungi (%v) -- test ini butuh MySQL dev nyata untuk menguji FOR UPDATE locking", err)
	}
	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Skip("skip: DB dev tidak bisa di-ping")
	}
	return db
}

// setupConcurrencyTestProduct membuat produk uji sementara dengan rantai
// Kardus (anchor) -> Botol (24 Botol = 1 Kardus), stok awal cuma di Botol
// sejumlah botolStock, supaya hasil akhir gampang diverifikasi (harus tepat
// 0 kalau semua goroutine berhasil tanpa lost update).
func setupConcurrencyTestProduct(t *testing.T, db *gorm.DB, botolStock int) int {
	t.Helper()

	// Coffee: dua satuan apa saja dari master Unit DB dev (id satuan berbeda dengan Toko).
	var unitIDs []int
	if err := db.Raw(`SELECT id FROM units ORDER BY id LIMIT 2`).Scan(&unitIDs).Error; err != nil || len(unitIDs) < 2 {
		t.Skip("skip: butuh minimal 2 satuan di master Unit DB dev Coffee")
	}
	kardusUnit, botolUnit := unitIDs[0], unitIDs[1]

	// id dibaca berdasarkan nama/kunci, bukan LAST_INSERT_ID(): db di sini pool koneksi, jadi
	// LAST_INSERT_ID() bisa terbaca dari koneksi lain (id salah, teardown bisa salah sasaran).
	name := fmt.Sprintf("TEST STOCKDELTA CONCURRENCY %d", nowSuffix())
	res := db.Exec(`INSERT INTO products (name, unit_id, min_stock, is_active) VALUES (?, ?, 0, 1)`, name, kardusUnit)
	if res.Error != nil {
		t.Fatalf("gagal buat produk uji: %v", res.Error)
	}
	var productID int
	if err := db.Raw("SELECT id FROM products WHERE name = ?", name).Scan(&productID).Error; err != nil || productID == 0 {
		t.Fatalf("gagal ambil id produk uji: %v", err)
	}

	res = db.Exec(`INSERT INTO product_packages (product_id, unit_id, qty, ref_qty, ref_package_id, is_default, stock, reserved_qty, is_active) VALUES (?, ?, 1, NULL, NULL, 1, 0, 0, 1)`, productID, kardusUnit)
	if res.Error != nil {
		t.Fatalf("gagal buat package Kardus uji: %v", res.Error)
	}
	var kardusPackageID int
	if err := db.Raw("SELECT id FROM product_packages WHERE product_id = ? AND is_default = 1", productID).Scan(&kardusPackageID).Error; err != nil || kardusPackageID == 0 {
		t.Fatalf("gagal ambil id package Kardus: %v", err)
	}

	res = db.Exec(`INSERT INTO product_packages (product_id, unit_id, qty, ref_qty, ref_package_id, is_default, stock, reserved_qty, is_active) VALUES (?, ?, 24, 1, ?, 0, ?, 0, 1)`, productID, botolUnit, kardusPackageID, botolStock)
	if res.Error != nil {
		t.Fatalf("gagal buat package Botol uji: %v", res.Error)
	}
	if err := db.Raw("SELECT id FROM product_packages WHERE product_id = ? AND ref_package_id IS NOT NULL", productID).Scan(&botolPackageIDForTest).Error; err != nil || botolPackageIDForTest == 0 {
		t.Fatalf("gagal ambil id package Botol: %v", err)
	}

	return productID
}

func teardownConcurrencyTestProduct(t *testing.T, db *gorm.DB, productID int) {
	t.Helper()
	// Kemasan anak (merujuk kemasan lain) dihapus lebih dulu agar foreign key ref_package_id tidak menolak.
	for _, q := range []string{
		"DELETE FROM stock_mutations WHERE product_id = ?",
		"DELETE FROM product_packages WHERE product_id = ? AND ref_package_id IS NOT NULL",
		"DELETE FROM product_packages WHERE product_id = ?",
		"DELETE FROM products WHERE id = ?",
	} {
		if err := db.Exec(q, productID).Error; err != nil {
			t.Errorf("teardown gagal (%s): %v", q, err)
		}
	}
}

var nowSuffixCounter int

// nowSuffix menghasilkan angka unik ringan supaya nama produk uji tidak
// bentrok kalau test dijalankan berkali-kali berturut-turut.
func nowSuffix() int64 {
	nowSuffixCounter++
	return time.Now().UnixNano() + int64(nowSuffixCounter)
}
