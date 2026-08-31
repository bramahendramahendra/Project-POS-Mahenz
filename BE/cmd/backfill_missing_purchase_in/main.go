// Skrip one-off: menambal baris mutasi 'in' yang HILANG untuk item pembelian
// (purchase_items) pada PO ber-status 'active' yang tidak pernah tercatat di
// stock_mutations.
//
// LATAR BELAKANG (lihat docs/PERBAIKAN_MUTASI_STOK_EDIT_PEMBELIAN.md &
// docs/ANALISIS_REKONSILIASI_STOK_32_PRODUK.md):
// Di data prod, sebagian purchase_items TIDAK punya baris 'in' pasangannya di
// ledger stock_mutations (barang yang sebenarnya dibeli, tapi jejak stok
// masuknya tidak pernah tercatat -- kemungkinan bug/alur lama). Akibatnya saat
// backfill_stock_restore me-replay riwayat, barang itu "tak terlihat"
// (hanya penjualan 'out' yang tercatat) -> stok minus / selisih -> produk
// ditandai needs_stock_review.
//
// Skrip ini MEMBANGKITKAN baris 'in' yang hilang itu lewat jalur RESMI
// product_repo.ApplyStockDelta (fungsi terpusat satu-satunya yang boleh
// mengubah product_packages.stock), memakai quantity & package_id ASLI dari
// purchase_items. Tidak ada risiko dobel-hitung: skrip HANYA menyentuh item
// yang benar-benar BELUM punya baris 'in' sama sekali (idempotent via NOT
// EXISTS) -- jadi aman dijalankan berulang.
//
// URUTAN dalam alur migrasi (lihat docs/MIGRASI_STOK_PROD_KE_SKEMA_BARU.md):
//  1. Restore prod -> pos_retail_db
//  2. CREATE TABLE products_stock_backup (SEBELUM BE jalan)
//  3. Jalankan BE sekali (migrasi 003-009, products.stock di-drop)
//  4. go run ./cmd/backfill_purchase_package_id   (isi purchase_items.package_id)
//  5. go run ./cmd/backfill_missing_purchase_in   <-- SKRIP INI
//  6. go run ./cmd/backfill_stock_restore         (hitung ulang stok final)
//
// Kenapa langkah 5 SEBELUM 6: skrip ini menyisipkan baris 'in' ke ledger.
// backfill_stock_restore (langkah 6) mereset product_packages.stock ke 0 lalu
// mereplay SELURUH riwayat -- termasuk baris 'in' baru dari skrip ini -- jadi
// stok final tetap dihitung oleh backfill_stock_restore, skrip ini cukup
// memastikan ledger-nya lengkap dulu. (Update stok yang dilakukan skrip ini
// akan ditimpa/dihitung ulang oleh langkah 6 -- itu memang disengaja.)
//
// Jalankan: go run ./cmd/backfill_missing_purchase_in
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	product_repo "pos_api/domain/product/repo"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// defaultDSN dipakai kalau env MIGRATION_DSN tidak diset (local WAMP, root
// tanpa password). Di production, set MIGRATION_DSN dulu -- sama seperti
// skrip backfill_stock_restore.
const defaultDSN = "root:@tcp(127.0.0.1:3306)/pos_retail_db?charset=utf8&parseTime=True&loc=Local"

func resolveDSN() string {
	if v := os.Getenv("MIGRATION_DSN"); v != "" {
		return v
	}
	return defaultDSN
}

// missingItem = satu baris purchase_items (PO active) yang belum punya baris
// 'in' pasangannya di stock_mutations.
type missingItem struct {
	PurchaseItemID int
	PurchaseID     int
	ProductID      int
	PackageID      int
	Quantity       float64
	PurchaseCode   string
	PurchaseDate   time.Time
}

// Query inti: purchase_items dari PO active yang TIDAK punya baris 'in' di
// stock_mutations untuk (reference_id=purchase_id, product_id). Ini yang bikin
// skrip idempotent -- item yang sudah punya 'in' (mis. dari run sebelumnya
// atau dari jalur normal) tidak akan diproses lagi.
const missingQuery = `
	SELECT pi.id, pi.purchase_id, pi.product_id, pi.package_id, pi.quantity, p.purchase_code, p.purchase_date
	FROM purchase_items pi
	JOIN purchases p ON p.id = pi.purchase_id AND p.status = 'active'
	WHERE pi.package_id IS NOT NULL
	  AND NOT EXISTS (
		SELECT 1 FROM stock_mutations sm
		WHERE sm.reference_type = 'purchase'
		  AND sm.reference_id = pi.purchase_id
		  AND sm.mutation_type = 'in'
		  AND sm.product_id = pi.product_id
	  )
	ORDER BY pi.purchase_id, pi.id`

func main() {
	dryRun := false
	for _, a := range os.Args[1:] {
		if a == "--dry-run" || a == "-n" {
			dryRun = true
		}
	}

	db, err := gorm.Open(gmysql.Open(resolveDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql.DB: %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	var items []missingItem
	if err := db.Raw(missingQuery).Scan(&items).Error; err != nil {
		log.Fatalf("query item yang bolong 'in': %v", err)
	}

	fmt.Printf("=== Backfill mutasi 'in' yang hilang (purchase_items PO active) ===\n")
	fmt.Printf("Item purchase yang BELUM punya baris 'in' di ledger: %d\n", len(items))
	if dryRun {
		fmt.Println("(DRY RUN -- tidak menulis apa pun)")
	}
	fmt.Println()

	if len(items) == 0 {
		fmt.Println("Tidak ada yang perlu ditambal. Selesai.")
		return
	}

	var okCount, skipCount int
	for _, it := range items {
		line := fmt.Sprintf(" - PO %s (id %d) produk %d pkg %d qty %.4f",
			it.PurchaseCode, it.PurchaseID, it.ProductID, it.PackageID, it.Quantity)

		if dryRun {
			fmt.Println("[DRY]" + line)
			okCount++
			continue
		}

		// Tiap item dibungkus transaksi sendiri: satu item yang gagal (mis.
		// produk rantai satuannya bercabang -> ditolak ComputeStockDelta) tidak
		// membatalkan item lain yang sudah berhasil.
		err := db.Transaction(func(tx *gorm.DB) error {
			// user_id sengaja NULL (kolom nullable, ON DELETE SET NULL): mutasi
			// hasil backfill bukan aksi user tertentu, dan mengisi id yang tidak
			// ada di tabel users akan melanggar FK stock_mutations_ibfk_2.
			_, e := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     it.ProductID,
				PackageID:     it.PackageID,
				Quantity:      it.Quantity,
				Direction:     "in",
				MutationType:  "in",
				ReferenceType: "purchase",
				ReferenceID:   it.PurchaseID,
				Notes:         fmt.Sprintf("Backfill: mutasi 'in' yang hilang untuk PO %s (rekonstruksi dari purchase_items)", it.PurchaseCode),
				UserID:        nil,
			})
			if e != nil {
				return e
			}
			// PENTING (urutan replay): ApplyStockDelta set created_at = now(),
			// tapi backfill_stock_restore me-replay urut created_at. Kalau baris
			// 'in' ini created_at-nya "sekarang", ia akan diproses SETELAH
			// penjualan lama -> stok minus. Jadi backdate created_at ke tanggal
			// PO-nya supaya 'in' diproses SEBELUM 'out' (kronologi benar).
			return tx.Exec(
				`UPDATE stock_mutations SET created_at = ?
				 WHERE reference_type='purchase' AND reference_id = ?
				   AND product_id = ? AND mutation_type='in'
				   AND notes LIKE 'Backfill:%'`,
				it.PurchaseDate, it.PurchaseID, it.ProductID,
			).Error
		})
		if err != nil {
			fmt.Printf("[SKIP]%s -- %v\n", line, err)
			skipCount++
			continue
		}
		fmt.Printf("[OK]  %s\n", line)
		okCount++
	}

	fmt.Printf("\n=== Ringkasan ===\n")
	fmt.Printf("Berhasil disisipkan 'in' : %d\n", okCount)
	fmt.Printf("Dilewati (gagal/error)   : %d\n", skipCount)
	if skipCount > 0 {
		fmt.Println("\nItem yang dilewati kemungkinan produk dengan rantai satuan bercabang")
		fmt.Println("(celah #10) -- perlu dirapikan strukturnya dulu, lalu jalankan ulang.")
	}
}
