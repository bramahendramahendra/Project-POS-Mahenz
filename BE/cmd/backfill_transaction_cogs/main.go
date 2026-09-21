// Skrip one-off: koreksi HPP/modal historis pada transaction_items.purchase_price.
//
// LATAR BELAKANG (lihat docs/RANCANGAN_KOREKSI_HPP_HISTORIS.md):
// Sebagian transaksi lama menyimpan transaction_items.purchase_price (modal per
// satuan JUAL) dengan level satuan yang SALAH -- memakai harga satuan BESAR
// (mis. per Slop) untuk penjualan satuan KECIL (mis. per Pack/Batang) tanpa
// dikali faktor konversi. Akibatnya HPP membengkak 10-100x -> Laba Rugi tampak
// rugi padahal tidak. Bug pencatatan sudah diperbaiki utk transaksi baru; ini
// membetulkan DATA LAMA.
//
// METODE (akurat, terbukti di data dev):
//
//		modal_benar = anchor_cost_sekarang * ResolvePackageFactor(paket satuan jual)
//	  - anchor_cost_sekarang = product_packages.purchase_price paket anchor (is_default=1).
//	    Master paket SEKARANG sudah bersih/rapi -> sumber paling akurat.
//	  - ResolvePackageFactor menelusuri rantai ref_package_id PENUH ke anchor (sama
//	    persis dengan yang dipakai kasir/pembelian). WAJIB multi-level, bukan 1-level.
//
// KENAPA TIDAK last-cost dari purchase_items: purchase_items sendiri KOTOR
// (conversion_qty tak konsisten), jadi tidak bisa dipercaya sbg acuan.
//
// KEAMANAN:
//   - HANYA menyentuh transaction_items milik transaksi status='completed'
//     (transaksi void di-skip; laporan COGS memang exclude void).
//   - purchase_price write-once di kode aplikasi -> aman ditimpa historis.
//   - Baris MERAGUKAN (paket tak ketemu / anchor 0 / faktor error / modal>=harga jual)
//     TIDAK dikoreksi otomatis -> dicatat sbg "MANUAL" untuk ditinjau lewat menu
//     Rekonsiliasi Modal (admin). Bukan blocker.
//   - Setiap perubahan dicatat ke tabel audit transaction_items_cogs_audit
//     (nilai lama disimpan sebelum ditimpa) -> bisa ditelusuri & di-rollback.
//
// MODE:
//
//	default          -> DRY RUN (tidak menulis apa pun; hanya laporan)
//	flag --apply     -> benar-benar menulis (UPDATE + audit) dalam transaksi
//
// Jalankan dari folder BE:
//
//	go run ./cmd/backfill_transaction_cogs            (dry-run)
//	go run ./cmd/backfill_transaction_cogs --apply    (eksekusi)
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"pos_api/cmd/internal/migrationdb"
	product_model "pos_api/domain/product/model"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// txItem = satu baris transaction_items dari transaksi completed.
type txItem struct {
	ID              int
	TransactionID   int
	TransactionCode string
	ProductID       *int
	ProductName     string
	Quantity        float64
	Price           float64 // harga jual per satuan
	PurchasePrice   float64 // modal tersimpan (mungkin salah)
	UnitID          *int    // = package_id paket satuan jual
}

// hasil klasifikasi per baris
type classified struct {
	item    txItem
	newCost float64
	factor  float64
	anchor  float64
	reason  string // kosong = OK; selain itu = alasan MANUAL
}

const selectItemsQuery = `
	SELECT ti.id, ti.transaction_id, t.transaction_code, ti.product_id, ti.product_name,
	       ti.quantity, ti.price, ti.purchase_price, ti.unit_id
	FROM transaction_items ti
	JOIN transactions t ON t.id = ti.transaction_id
	WHERE t.status = 'completed'
	ORDER BY ti.product_id, ti.id`

// audit table -- dibuat kalau belum ada (idempotent). Menyimpan nilai lama sebelum ditimpa.
const createAuditTableQuery = `
	CREATE TABLE IF NOT EXISTS transaction_items_cogs_audit (
		id                  INT AUTO_INCREMENT PRIMARY KEY,
		transaction_item_id INT           NOT NULL,
		transaction_id      INT           NOT NULL,
		product_id          INT           NULL,
		old_purchase_price  DECIMAL(15,2) NOT NULL,
		new_purchase_price  DECIMAL(15,2) NOT NULL,
		anchor_cost         DECIMAL(15,2) NOT NULL,
		factor              DECIMAL(18,8) NOT NULL,
		note                VARCHAR(255)  NOT NULL DEFAULT '',
		created_at          DATETIME      DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_cogs_audit_item (transaction_item_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

const insertAuditQuery = `
	INSERT INTO transaction_items_cogs_audit
	  (transaction_item_id, transaction_id, product_id, old_purchase_price, new_purchase_price, anchor_cost, factor, note)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

const updatePurchasePriceQuery = `UPDATE transaction_items SET purchase_price = ? WHERE id = ?`

// Upsert baris MANUAL ke daftar tinjauan (menu Rekonsiliasi Modal).
// Idempotent: kalau baris (transaction_item_id) sudah ada & masih 'pending',
// perbarui snapshot-nya; kalau sudah 'resolved'/'skipped' JANGAN timpa statusnya.
const upsertReviewQuery = `
	INSERT INTO transaction_cost_review
	  (transaction_item_id, transaction_id, product_id, reason, old_purchase_price, suggested_price, sell_price, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')
	ON DUPLICATE KEY UPDATE
	  reason = IF(status = 'pending', VALUES(reason), reason),
	  old_purchase_price = IF(status = 'pending', VALUES(old_purchase_price), old_purchase_price),
	  suggested_price = IF(status = 'pending', VALUES(suggested_price), suggested_price),
	  sell_price = IF(status = 'pending', VALUES(sell_price), sell_price),
	  updated_at = NOW()`

func main() {
	apply := false
	for _, a := range os.Args[1:] {
		if a == "--apply" {
			apply = true
		}
	}

	dsn, err := migrationdb.ResolveDSN()
	if err != nil {
		log.Fatalf("resolve DSN: %v", err)
	}
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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

	// Muat semua paket per produk sekali, untuk resolve faktor & anchor cost.
	packagesByProduct, anchorByProduct, err := loadPackages(db)
	if err != nil {
		log.Fatalf("muat product_packages: %v", err)
	}

	var items []txItem
	if err := db.Raw(selectItemsQuery).Scan(&items).Error; err != nil {
		log.Fatalf("query transaction_items: %v", err)
	}

	// Klasifikasi tiap baris.
	var okList, manualList []classified
	var cogsOld, cogsNew float64
	for _, it := range items {
		cogsOld += it.Quantity * it.PurchasePrice
		c := classify(it, packagesByProduct, anchorByProduct)
		if c.reason == "" {
			okList = append(okList, c)
			cogsNew += it.Quantity * c.newCost
		} else {
			manualList = append(manualList, c)
			cogsNew += it.Quantity * it.PurchasePrice // baris manual: pertahankan nilai lama dulu
		}
	}

	// ---- Laporan ----
	fmt.Println("=== Koreksi HPP historis (transaction_items.purchase_price) ===")
	if apply {
		fmt.Println("MODE: APPLY (menulis perubahan + audit)")
	} else {
		fmt.Println("MODE: DRY RUN (tidak menulis apa pun)")
	}
	fmt.Printf("Total baris (completed)      : %d\n", len(items))
	fmt.Printf("  Akan dikoreksi otomatis    : %d\n", len(okList))
	fmt.Printf("  Perlu tinjauan MANUAL      : %d\n", len(manualList))
	fmt.Printf("Total COGS (lifetime) SEBELUM: %.2f\n", cogsOld)
	fmt.Printf("Total COGS (lifetime) SESUDAH: %.2f  (baris manual pakai nilai lama dulu)\n", cogsNew)
	fmt.Println()

	// Contoh baris yang akan dikoreksi (maks 25).
	fmt.Println("--- Contoh koreksi (maks 25) ---")
	shown := 0
	changed := 0
	for _, c := range okList {
		if c.item.PurchasePrice == c.newCost {
			continue // tidak berubah
		}
		changed++
		if shown < 25 {
			fmt.Printf(" [%s] %s: modal %.2f -> %.2f (anchor %.2f x faktor %.6f)\n",
				c.item.TransactionCode, c.item.ProductName,
				c.item.PurchasePrice, c.newCost, c.anchor, c.factor)
			shown++
		}
	}
	fmt.Printf("(baris OK yang nilainya BERUBAH: %d dari %d)\n\n", changed, len(okList))

	// Daftar MANUAL (semua, biar admin tahu).
	fmt.Println("--- Perlu tinjauan MANUAL (Rekonsiliasi Modal) ---")
	manualProducts := map[int]string{}
	for _, c := range manualList {
		pid := 0
		if c.item.ProductID != nil {
			pid = *c.item.ProductID
		}
		manualProducts[pid] = c.item.ProductName
	}
	for pid, name := range manualProducts {
		fmt.Printf(" - produk %d (%s)\n", pid, name)
	}
	if len(manualList) > 0 {
		fmt.Println("   contoh baris manual (maks 15):")
		for i, c := range manualList {
			if i >= 15 {
				break
			}
			fmt.Printf("   [%s] %s: modal tersimpan %.2f, hitung %.2f, harga jual %.2f -> %s\n",
				c.item.TransactionCode, c.item.ProductName,
				c.item.PurchasePrice, c.newCost, c.item.Price, c.reason)
		}
	}
	fmt.Println()

	if !apply {
		fmt.Println("DRY RUN selesai. Jalankan ulang dengan --apply untuk menulis perubahan.")
		fmt.Println("PASTIKAN sudah BACKUP database sebelum --apply.")
		return
	}

	// ---- APPLY ----
	if err := db.Exec(createAuditTableQuery).Error; err != nil {
		log.Fatalf("buat tabel audit: %v", err)
	}
	now := time.Now()
	var written int
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, c := range okList {
			if c.item.PurchasePrice == c.newCost {
				continue // tidak berubah, tidak perlu ditulis/di-audit
			}
			pid := any(nil)
			if c.item.ProductID != nil {
				pid = *c.item.ProductID
			}
			if err := tx.Exec(insertAuditQuery,
				c.item.ID, c.item.TransactionID, pid,
				c.item.PurchasePrice, c.newCost, c.anchor, c.factor,
				"koreksi HPP historis (anchor x faktor)",
			).Error; err != nil {
				return fmt.Errorf("audit item %d: %w", c.item.ID, err)
			}
			if err := tx.Exec(updatePurchasePriceQuery, c.newCost, c.item.ID).Error; err != nil {
				return fmt.Errorf("update item %d: %w", c.item.ID, err)
			}
			written++
		}

		// Daftarkan baris MANUAL ke tabel tinjauan (menu Rekonsiliasi Modal).
		// Idempotent: baris yang sudah 'resolved'/'skipped' tidak di-reset.
		for _, c := range manualList {
			pid := any(nil)
			if c.item.ProductID != nil {
				pid = *c.item.ProductID
			}
			if err := tx.Exec(upsertReviewQuery,
				c.item.ID, c.item.TransactionID, pid,
				c.reason, c.item.PurchasePrice, c.newCost, c.item.Price,
			).Error; err != nil {
				return fmt.Errorf("daftar tinjauan item %d: %w", c.item.ID, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("APPLY gagal (rollback): %v", err)
	}
	_ = now
	fmt.Printf("APPLY selesai. Baris diperbarui: %d. Audit tersimpan di transaction_items_cogs_audit.\n", written)
	fmt.Printf("Baris manual (%d) didaftarkan ke transaction_cost_review -- tangani lewat menu Rekonsiliasi Modal.\n", len(manualList))
}

// classify menghitung modal benar & menentukan apakah baris OK atau MANUAL.
func classify(it txItem, pkgByProduct map[int][]*product_model.ProductPackage, anchorByProduct map[int]float64) classified {
	c := classified{item: it}

	if it.ProductID == nil {
		c.reason = "product_id NULL"
		return c
	}
	pkgs, ok := pkgByProduct[*it.ProductID]
	if !ok || len(pkgs) == 0 {
		c.reason = "produk tidak punya product_packages"
		return c
	}
	anchor, ok := anchorByProduct[*it.ProductID]
	if !ok || anchor <= 0 {
		c.reason = "anchor cost 0/tidak ada"
		return c
	}
	if it.UnitID == nil || *it.UnitID == 0 {
		c.reason = "unit_id (paket jual) NULL"
		return c
	}
	factor, err := product_model.ResolvePackageFactor(pkgs, *it.UnitID)
	if err != nil || factor <= 0 {
		c.reason = "faktor tak bisa di-resolve (rantai bercabang/melingkar)"
		return c
	}

	c.anchor = anchor
	c.factor = factor
	c.newCost = round2(anchor * factor)

	// Cross-check: modal hasil hitung TIDAK boleh >= harga jual (indikator masih janggal).
	// Kecuali harga jual 0 (data aneh) -> juga manual.
	if it.Price <= 0 {
		c.reason = "harga jual 0"
		return c
	}
	if c.newCost >= it.Price {
		c.reason = "modal hasil >= harga jual"
		return c
	}
	return c // OK
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

// loadPackages memuat semua paket per produk + anchor cost per produk.
func loadPackages(db *gorm.DB) (map[int][]*product_model.ProductPackage, map[int]float64, error) {
	var rows []*product_model.ProductPackage
	// unit_name via join tidak wajib untuk faktor; ambil kolom yang dipakai ResolvePackageFactor.
	if err := db.Raw(`
		SELECT id, product_id, unit_id, package_name, ref_package_id, qty, ref_qty,
		       purchase_price, is_default, COALESCE(is_active,1) AS is_active
		FROM product_packages`).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}
	byProduct := make(map[int][]*product_model.ProductPackage)
	anchor := make(map[int]float64)
	for _, p := range rows {
		byProduct[p.ProductID] = append(byProduct[p.ProductID], p)
		if p.IsDefault {
			anchor[p.ProductID] = p.PurchasePrice
		}
	}
	return byProduct, anchor, nil
}
