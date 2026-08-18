package repo

import (
	"fmt"

	model_product "pos_api/domain/product/model"

	"gorm.io/gorm"
)

// Fase 5 (docs/RENCANA_PERBAIKAN_STOK_PRESISI.md) -- helper baca-saja lintas
// domain (product, report, business_summary) untuk mengagregasi stok dari
// product_packages, bukan lagi bergantung langsung pada products.stock.
// Package-level function (bukan method), sama seperti ApplyStockDelta &
// ResolveDefaultPackageID di stock_delta_repo.go -- supaya domain lain
// (report_repo.go, business_summary_repo.go) bisa reuse tanpa duplikasi.

const getAllActivePackagesQuery = `
	SELECT pp.id, pp.product_id, pp.unit_id, COALESCE(u.name,'') AS unit_name, COALESCE(u.abbreviation,'') AS abbreviation,
	       COALESCE(pp.package_name,'') AS package_name, pp.ref_package_id, pp.qty, pp.ref_qty,
	       pp.purchase_price, pp.selling_price, pp.is_default, pp.stock, pp.reserved_qty, pp.is_active
	FROM product_packages pp
	JOIN units u ON u.id = pp.unit_id
	WHERE pp.is_active = 1`

// GetAllActivePackages memuat SEMUA baris product_packages aktif (lintas
// produk) dalam satu query -- dipakai jalur baca massal (list produk, laporan
// stok, dashboard) supaya tidak N+1 query per produk.
func GetAllActivePackages(db *gorm.DB) ([]*model_product.ProductPackage, error) {
	var rows []*model_product.ProductPackage
	if err := db.Raw(getAllActivePackagesQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ContinuousUnitMap -- versi exported dari continuousUnitMap (stock_delta_repo.go)
// supaya domain lain juga bisa pakai tanpa reimplementasi query yang sama.
func ContinuousUnitMap(db *gorm.DB) (map[int]bool, error) {
	return continuousUnitMap(db)
}

// BuildStockSummaries mengagregasi StockSummary (Fase 5) untuk SEMUA produk
// yang punya baris product_packages aktif. minStockByProduct wajib diisi
// pemanggil (dari products.min_stock, biasanya sudah dimuat di query utama
// pemanggil, jadi tidak perlu query ulang di sini).
//
// Produk yang gagal dianalisis (rantai bercabang/celah #10, tidak ada
// anchor, dst) tetap dimasukkan ke hasil TAPI dengan fallback: angka baris
// anchor (is_default=1) apa adanya, tanpa breakdown/agregasi penuh lintas
// level, dan IsLowStock selalu false (datanya sendiri diragukan -- celah #20,
// FE menampilkan badge needs_stock_review, bukan alert stok yang salah).
// Fase 8: fallback ini TIDAK LAGI baca products.stock/reserved_qty (kolom
// itu dihapus) -- diambil langsung dari baris product_packages yang sudah
// dimuat di atas, jadi tetap akurat untuk baris yang trusted (anchor) itu
// sendiri walau rantai penuhnya belum bisa dianalisis.
func BuildStockSummaries(db *gorm.DB, minStockByProduct map[int]float64) (map[int]*model_product.StockSummary, error) {
	packages, err := GetAllActivePackages(db)
	if err != nil {
		return nil, fmt.Errorf("muat product_packages aktif: %w", err)
	}
	continuous, err := ContinuousUnitMap(db)
	if err != nil {
		return nil, fmt.Errorf("muat units.is_continuous: %w", err)
	}

	grouped := make(map[int][]*model_product.ProductPackage, len(minStockByProduct))
	for _, p := range packages {
		grouped[p.ProductID] = append(grouped[p.ProductID], p)
	}

	result := make(map[int]*model_product.StockSummary, len(grouped))
	for productID, pkgs := range grouped {
		summary, sErr := model_product.ComputeStockSummary(pkgs, continuous, minStockByProduct[productID])
		if sErr != nil {
			var anchor *model_product.ProductPackage
			for _, p := range pkgs {
				if p.IsDefault {
					anchor = p
					break
				}
			}
			if anchor == nil {
				continue
			}
			result[productID] = &model_product.StockSummary{
				AnchorStock:    anchor.Stock,
				AnchorReserved: anchor.ReservedQty,
				IsLowStock:     false,
			}
			continue
		}
		result[productID] = summary
	}
	return result, nil
}
