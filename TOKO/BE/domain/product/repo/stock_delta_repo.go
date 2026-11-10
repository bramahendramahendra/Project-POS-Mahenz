package repo

import (
	"errors"
	"fmt"

	model_product "pos_api/domain/product/model"
	app_errors "pos_api/errors"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const (
	getPackagesForUpdateQuery = `
		SELECT pp.id, pp.product_id, pp.unit_id, COALESCE(u.name,'') AS unit_name, COALESCE(u.abbreviation,'') AS abbreviation,
		       COALESCE(pp.package_name,'') AS package_name, pp.ref_package_id, pp.qty, pp.ref_qty,
		       pp.purchase_price, pp.selling_price, pp.is_default, pp.stock, pp.reserved_qty, pp.is_active
		FROM product_packages pp
		JOIN units u ON u.id = pp.unit_id
		WHERE pp.product_id = ? AND pp.is_active = 1
		FOR UPDATE`

	getContinuousUnitsQuery = `SELECT id, is_continuous FROM units`

	getProductReviewFlagQuery = `SELECT needs_stock_review FROM products WHERE id = ? FOR UPDATE`

	updatePackageStockQuery = `UPDATE product_packages SET stock = ?, updated_at = ? WHERE id = ?`

	insertStockMutationQuery = `
		INSERT INTO stock_mutations
			(product_id, package_id, mutation_type, quantity, stock_before, stock_after, reference_type, reference_id, notes, user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// ApplyStockDeltaParams -- lihat docs/RENCANA_PERBAIKAN_STOK_PRESISI.md
// bagian "Aturan operasional" #5.
type ApplyStockDeltaParams struct {
	ProductID     int
	PackageID     int
	Quantity      float64
	Direction     model_product.StockDirection
	MutationType  string // 'in' | 'out' | 'adjustment' | 'void' | 'return' | 'void_purchase' | 'expired'
	ReferenceType string
	ReferenceID   int
	Notes         string
	UserID        *int
}

// ApplyStockDelta adalah fungsi terpusat SATU-SATUNYA yang boleh mengubah
// product_packages.stock (Aturan Operasional #5, Fase 2). Menerima tx yang
// SUDAH DIBUKA oleh pemanggil (bukan membuka transaksi sendiri) supaya
// pemanggil bisa membungkus operasi ini bersama langkah lain (insert
// transaction_items, dst) dalam satu unit atomik yang sama -- ini yang
// dipakai untuk menutup celah #18 (jalur Create/Void yang sebelumnya tidak
// dibungkus transaksi sama sekali) begitu Fase 4 memindahkan jalur-jalur itu
// ke fungsi ini.
//
// Mengunci SEMUA baris product_packages milik produk ini (FOR UPDATE)
// sebelum menghitung, bukan cuma baris level yang sedang ditransaksikan
// (Aturan Operasional #3) -- supaya breakdown antar level tetap konsisten
// kalau ada 2 transaksi berjalan bersamaan pada produk yang sama.
//
// PENTING (celah #21): Quantity di sini harus qty ASLI di satuan PackageID
// (integer untuk satuan diskrit) -- bukan angka conversion_qty yang sudah
// dihitung & disimpan sebelumnya di tempat lain. Faktor konversi selalu
// dihitung ulang fresh oleh ComputeStockDelta lewat ResolvePackageFactor,
// supaya void/pembalikan transaksi tidak pernah membaca balik nilai yang
// sudah berpotensi kepotong presisi.
func ApplyStockDelta(tx *gorm.DB, p ApplyStockDeltaParams) (*model_product.StockDeltaResult, error) {
	var needsReview bool
	if err := tx.Raw(getProductReviewFlagQuery, p.ProductID).Scan(&needsReview).Error; err != nil {
		return nil, fmt.Errorf("baca needs_stock_review: %w", err)
	}

	var packages []*model_product.ProductPackage
	if err := tx.Raw(getPackagesForUpdateQuery, p.ProductID).Scan(&packages).Error; err != nil {
		return nil, fmt.Errorf("kunci & baca product_packages: %w", err)
	}
	if len(packages) == 0 {
		return nil, errors.New("produk tidak punya baris product_packages aktif")
	}

	continuous, err := continuousUnitMap(tx)
	if err != nil {
		return nil, fmt.Errorf("baca units.is_continuous: %w", err)
	}

	result, err := model_product.ComputeStockDelta(model_product.StockDeltaInput{
		Packages:         packages,
		ContinuousUnitID: continuous,
		PackageID:        p.PackageID,
		Quantity:         p.Quantity,
		Direction:        p.Direction,
		NeedsStockReview: needsReview,
	})
	if err != nil {
		return nil, err
	}

	now := time_helper.GetTimeNow()
	for _, u := range result.Updates {
		if err := tx.Exec(updatePackageStockQuery, u.Stock, now, u.PackageID).Error; err != nil {
			return nil, fmt.Errorf("update stock package %d: %w", u.PackageID, err)
		}
	}

	if err := tx.Exec(insertStockMutationQuery,
		p.ProductID, p.PackageID, p.MutationType, result.MutationQtyAnchor,
		result.StockBeforeAnchor, result.StockAfterAnchor,
		p.ReferenceType, p.ReferenceID, p.Notes, p.UserID,
	).Error; err != nil {
		return nil, fmt.Errorf("insert stock_mutations: %w", err)
	}

	return result, nil
}

// ResolveDefaultPackageID -- helper lintas-domain (celah #15 pattern):
// kalau jalur tulis manapun tidak punya package_id eksplisit (mis. FE tidak
// mengirim), cari sendiri package anchor (is_default=true) produk itu.
// Dipakai purchase_repo.go, transaction_repo.go, dst -- daripada tiap domain
// reimplementasi query yang sama.
func ResolveDefaultPackageID(tx *gorm.DB, productID int) (int, bool) {
	var id int
	if err := tx.Raw(`SELECT id FROM product_packages WHERE product_id = ? AND is_default = 1 LIMIT 1`, productID).Scan(&id).Error; err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

// WrapStockError menerjemahkan error teknis (sentinel) dari ApplyStockDelta
// jadi *errors.BadRequestError yang dikenali middleware error handler global
// (supaya jadi HTTP 400 + pesan ramah, bukan 500 generik yang menyembunyikan
// pesan aslinya). Dipakai bersama lintas domain -- pola yang sama sebelumnya
// cuma ada di supplier_purchase/repo (wrapStockError, tidak exported, tidak
// bisa dipakai domain lain). Bug QA Fase A skenario 14: product_repo.go
// Update() TIDAK memanggil pembungkus apa pun, jadi ErrInsufficientStock
// (mis. turunkan stok manual sampai di bawah reserved_qty retur pending)
// bocor sebagai error mentah -> 500 Internal Server Error, bukan 400 dengan
// pesan jelas. Diperbaiki dengan helper terpusat ini, dipakai di Update().
func WrapStockError(err error, productID int) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, model_product.ErrInsufficientStock) {
		return &app_errors.BadRequestError{Message: fmt.Sprintf("Stok produk ID %d tidak mencukupi untuk perubahan ini (kemungkinan sebagian sedang ditahan retur/reservasi lain)", productID)}
	}
	if errors.Is(err, model_product.ErrNeedsStockReview) {
		return &app_errors.BadRequestError{Message: fmt.Sprintf("Produk ID %d ditandai perlu ditinjau manual (needs_stock_review), operasi stok diblokir sampai ditinjau admin", productID)}
	}
	if errors.Is(err, model_product.ErrBranchingChain) {
		return &app_errors.BadRequestError{Message: fmt.Sprintf("Produk ID %d punya struktur satuan bercabang, tidak bisa diproses otomatis -- hubungi admin", productID)}
	}
	return err
}

func continuousUnitMap(tx *gorm.DB) (map[int]bool, error) {
	type unitContinuousRow struct {
		ID           int  `gorm:"column:id"`
		IsContinuous bool `gorm:"column:is_continuous"`
	}
	var rows []unitContinuousRow
	if err := tx.Raw(getContinuousUnitsQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[int]bool, len(rows))
	for _, r := range rows {
		m[r.ID] = r.IsContinuous
	}
	return m, nil
}
