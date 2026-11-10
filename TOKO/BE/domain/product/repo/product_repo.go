package repo

import (
	"fmt"
	"sort"

	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const (
	// Fase 8: p.stock/p.reserved_qty (kolom lama) DIHAPUS dari sini -- semua
	// query di bawah cuma ambil field non-stok, angka stok final SELALU diisi
	// di lapisan Go lewat attachStockSummaries (agregasi product_packages,
	// Fase 5), termasuk fallback untuk produk yang gagal dianalisis (celah
	// #10/#20 -- lihat BuildStockSummaries di stock_read_repo.go).
	getAllProductsBase = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.needs_stock_review, COALESCE(p.stock_review_note, '') as stock_review_note,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE 1=1`

	getProductByIDQuery = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.needs_stock_review, COALESCE(p.stock_review_note, '') as stock_review_note, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.id = ? LIMIT 1`

	getProductByBarcodeQuery = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.needs_stock_review, COALESCE(p.stock_review_note, '') as stock_review_note, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.barcode = ? LIMIT 1`

	searchProductsQuery = `
		SELECT p.id, p.barcode, p.name, p.selling_price, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1 AND (p.name LIKE ? OR p.barcode LIKE ?)`

	// getAllActiveProductsForLowStockQuery -- perbandingan "stok menipis"
	// TIDAK dilakukan di SQL (anchor-vs-anchor, celah #12 -- salah alarm
	// begitu sisa < 1 unit anchor). Cuma ambil kandidat (semua produk aktif)
	// plus min_stock, perbandingan sebenarnya dilakukan di Go pakai
	// ComputeStockSummary (satuan terkecil) -- lihat GetLowStock().
	getAllActiveProductsForLowStockQuery = `
		SELECT p.id, p.name, p.min_stock, p.needs_stock_review, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1`

	getProductOptionsQuery      = `SELECT id, name FROM products WHERE is_active = 1 ORDER BY name`
	checkProductUsedQuery       = `SELECT COUNT(*) FROM transaction_items WHERE product_id = ?`
	checkProductPurchasedQuery  = `SELECT COUNT(*) FROM purchase_items WHERE product_id = ?`
	createProductQuery          = `INSERT INTO products (barcode, sku, name, category_id, purchase_price, selling_price, min_stock, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	insertAnchorPackageOnCreate = `INSERT INTO product_packages (product_id, unit_id, purchase_price, selling_price, is_default) VALUES (?, ?, ?, ?, 1)`
	getLastProductInsertIDQuery = `SELECT LAST_INSERT_ID()`
	updateProductQuery          = `UPDATE products SET barcode=?, sku=?, name=?, category_id=?, purchase_price=?, selling_price=?, min_stock=?, updated_at=? WHERE id=?`
	getAnchorPackageForUpdateQuery = `SELECT id, stock FROM product_packages WHERE product_id = ? AND is_default = 1 LIMIT 1`
	updateAnchorPackagePrice    = `UPDATE product_packages SET purchase_price=?, selling_price=?, updated_at=? WHERE product_id=? AND is_default=1`
	deleteProductQuery          = `DELETE FROM products WHERE id = ?`
	toggleProductStatusQuery    = `UPDATE products SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	markStockReviewedQuery      = `UPDATE products SET needs_stock_review = 0, stock_review_note = NULL, updated_at = ? WHERE id = ?`
	getAllProductsDefaultOrder  = ` ORDER BY p.name ASC`
	countProductsBase           = `SELECT COUNT(*) FROM products p WHERE 1=1`
)

// attachStockSummaries mengisi Stock/ReservedQty/IsLowStock tiap produk dari
// agregasi product_packages (Fase 5, celah #12) -- SATU-SATUNYA sumber sejak
// Fase 8 (kolom products.stock/reserved_qty sudah dihapus). Produk yang
// gagal dianalisis penuh (rantai bercabang dkk, celah #10/#20) tetap dapat
// entri di BuildStockSummaries lewat fallback baris anchor apa adanya --
// lihat komentar BuildStockSummaries di stock_read_repo.go.
func attachStockSummaries(db *gorm.DB, products []*model.Product) error {
	if len(products) == 0 {
		return nil
	}
	minStockByProduct := make(map[int]float64, len(products))
	for _, p := range products {
		minStockByProduct[p.ID] = p.MinStock
	}
	summaries, err := BuildStockSummaries(db, minStockByProduct)
	if err != nil {
		return err
	}
	for _, p := range products {
		if s, ok := summaries[p.ID]; ok {
			p.Stock = s.AnchorStock
			p.ReservedQty = s.AnchorReserved
			p.IsLowStock = s.IsLowStock
		}
	}
	return nil
}

func (r *productRepo) GetAll(req *dto.GetAllRequest) ([]*model.Product, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND (p.name LIKE ? OR p.barcode LIKE ?)`
		args = append(args, search, search)
	}
	if req.CategoryID != nil {
		conditions += ` AND p.category_id = ?`
		args = append(args, *req.CategoryID)
	}
	if req.IsActive != nil {
		conditions += ` AND p.is_active = ?`
		args = append(args, *req.IsActive)
	}

	allowedSortFields := map[string]string{
		"name":           "p.name",
		"selling_price":  "p.selling_price",
		"purchase_price": "p.purchase_price",
		"is_active":      "p.is_active",
	}
	orderClause := request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, getAllProductsDefaultOrder)

	// Fase 5/8 (celah #12): "stok menipis" DAN sortir berdasarkan stok TIDAK
	// BISA lagi dilakukan di SQL -- angka stok sekarang cuma ada lewat
	// agregasi product_packages (Fase 8: kolom products.stock sudah dihapus).
	// Kalau salah satu aktif, muat SEMUA kandidat yang cocok filter lain
	// (tanpa LIMIT/OFFSET), hitung di Go, baru filter+sortir+paginasi di
	// memori -- dataset produk cukup kecil (~ratusan) untuk ini aman.
	needsInMemory := req.LowStock || req.SortBy == "stock"
	if needsInMemory {
		var all []*model.Product
		query := getAllProductsBase + conditions + orderClause
		if err := r.db.Raw(query, args...).Scan(&all).Error; err != nil {
			return nil, 0, err
		}
		if err := attachStockSummaries(r.db, all); err != nil {
			return nil, 0, err
		}

		if req.SortBy == "stock" {
			sort.Slice(all, func(i, j int) bool {
				if req.SortOrder == "desc" {
					return all[i].Stock > all[j].Stock
				}
				return all[i].Stock < all[j].Stock
			})
		}

		filtered := all
		if req.LowStock {
			filtered = make([]*model.Product, 0, len(all))
			for _, p := range all {
				if p.IsLowStock {
					filtered = append(filtered, p)
				}
			}
		}

		_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
		total := int64(len(filtered))
		if offset >= len(filtered) {
			return []*model.Product{}, total, nil
		}
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		return filtered[offset:end], total, nil
	}

	var total int64
	if err := r.db.Raw(countProductsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	query := getAllProductsBase + conditions + orderClause + " LIMIT ? OFFSET ?"
	pagedArgs := append(append([]any{}, args...), limit, offset)

	var dataDB []*model.Product
	if err := r.db.Raw(query, pagedArgs...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	if err := attachStockSummaries(r.db, dataDB); err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *productRepo) GetOptions() ([]*model.ProductOption, error) {
	var dataDB []*model.ProductOption
	err := r.db.Raw(getProductOptionsQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productRepo) GetByID(id int) (*model.Product, error) {
	var dataDB model.Product
	err := r.db.Raw(getProductByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	if err := attachStockSummaries(r.db, []*model.Product{&dataDB}); err != nil {
		return nil, err
	}
	return &dataDB, nil
}

func (r *productRepo) GetByBarcode(barcode string) (*model.Product, error) {
	var dataDB model.Product
	err := r.db.Raw(getProductByBarcodeQuery, barcode).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	if err := attachStockSummaries(r.db, []*model.Product{&dataDB}); err != nil {
		return nil, err
	}
	return &dataDB, nil
}

// Search dipakai kasir (ProductSearch.tsx) -- dataset hasil selalu kecil
// (limit <= 50), jadi agregasi per-produk lewat product_packages (bukan lagi
// products.stock - reserved_qty langsung) di sini murah dilakukan per-row.
func (r *productRepo) Search(req *dto.SearchRequest) ([]*model.ProductSearchResult, error) {
	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	query := searchProductsQuery
	var args []any
	search := "%" + req.Q + "%"
	query += " LIMIT ?"
	args = append(args, search, search, limit)

	var dataDB []*model.ProductSearchResult
	err := r.db.Raw(query, args...).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}

	if len(dataDB) > 0 {
		productIDs := make([]int, len(dataDB))
		minStockByProduct := make(map[int]float64, len(dataDB))
		for i, v := range dataDB {
			productIDs[i] = v.ID
			minStockByProduct[v.ID] = v.MinStock
		}
		summaries, err := BuildStockSummaries(r.db, minStockByProduct)
		if err != nil {
			return nil, err
		}
		for _, v := range dataDB {
			if s, ok := summaries[v.ID]; ok {
				v.Stock = s.AnchorStock - s.AnchorReserved
			}
		}
	}
	return dataDB, nil
}

// GetLowStock -- Fase 5 (celah #12): dulu perbandingan (stock - reserved_qty)
// <= min_stock langsung di SQL, anchor-vs-anchor, salah alarm begitu sisa
// produk < 1 unit anchor penuh (mis. 0 Kardus + 20 Botol dari kapasitas 48).
// Sekarang perbandingan dilakukan di satuan TERKECIL lewat ComputeStockSummary.
func (r *productRepo) GetLowStock() ([]*model.LowStockProduct, error) {
	type candidateRow struct {
		ID               int     `gorm:"column:id"`
		Name             string  `gorm:"column:name"`
		MinStock         float64 `gorm:"column:min_stock"`
		NeedsStockReview bool    `gorm:"column:needs_stock_review"`
		UnitName         string  `gorm:"column:unit_name"`
	}
	var candidates []*candidateRow
	if err := r.db.Raw(getAllActiveProductsForLowStockQuery).Scan(&candidates).Error; err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return []*model.LowStockProduct{}, nil
	}

	minStockByProduct := make(map[int]float64, len(candidates))
	for _, c := range candidates {
		minStockByProduct[c.ID] = c.MinStock
	}
	summaries, err := BuildStockSummaries(r.db, minStockByProduct)
	if err != nil {
		return nil, err
	}

	dataDB := make([]*model.LowStockProduct, 0, len(candidates))
	for _, c := range candidates {
		s, ok := summaries[c.ID]
		if !ok {
			// Gagal dianalisis (rantai bercabang dkk) -- tidak dimasukkan ke
			// daftar stok menipis (datanya sendiri meragukan), badge
			// needs_stock_review di FE (Fase 6) yang memberi tahu admin.
			continue
		}
		if !s.IsLowStock {
			continue
		}
		dataDB = append(dataDB, &model.LowStockProduct{
			ID:       c.ID,
			Name:     c.Name,
			Stock:    s.AnchorStock - s.AnchorReserved,
			MinStock: c.MinStock,
			UnitName: c.UnitName,
		})
	}
	return dataDB, nil
}

func (r *productRepo) CountTransactionItems(productID int) (int, error) {
	var count int
	err := r.db.Raw(checkProductUsedQuery, productID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepo) CountPurchaseItems(productID int) (int, error) {
	var count int
	err := r.db.Raw(checkProductPurchasedQuery, productID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create menyimpan produk baru sekaligus paket anchor-nya (satuan pencatatan stok,
// permanen sejak dibuat) dan satuan lain (req.Packages, opsional) dalam SATU transaksi —
// supaya user bisa isi semua satuan produk langsung di form Tambah Produk, satu kali
// simpan, tanpa produk pernah ada dalam keadaan "setengah jadi" kalau ada yang gagal.
func (r *productRepo) Create(req *dto.CreateRequest) (int64, error) {
	var id int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(createProductQuery,
			req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
			req.SellingPrice, req.MinStock, req.UnitID,
		).Error; err != nil {
			return err
		}

		if err := tx.Raw(getLastProductInsertIDQuery).Scan(&id).Error; err != nil {
			return err
		}

		if err := tx.Exec(insertAnchorPackageOnCreate, id, req.UnitID, req.PurchasePrice, req.SellingPrice).Error; err != nil {
			return err
		}
		var anchorID int64
		if err := tx.Raw(getLastProductInsertIDQuery).Scan(&anchorID).Error; err != nil {
			return err
		}

		// Stok awal (kalau diisi admin) HARUS lewat ApplyStockDelta -- bukan
		// ditulis langsung -- supaya tercatat di stock_mutations & konsisten
		// dengan satu-satunya sumber kebenaran stok (product_packages), bukan
		// bergantung pada products.stock (kolom lama, dihapus Fase 8).
		if req.Stock > 0 {
			var userID *int
			if req.UserID > 0 {
				userID = &req.UserID
			}
			if _, err := ApplyStockDelta(tx, ApplyStockDeltaParams{
				ProductID:     int(id),
				PackageID:     int(anchorID),
				Quantity:      req.Stock,
				Direction:     model.StockIn,
				MutationType:  "adjustment",
				ReferenceType: "product_create",
				ReferenceID:   int(id),
				Notes:         "Stok awal saat produk dibuat",
				UserID:        userID,
			}); err != nil {
				return WrapStockError(err, int(id))
			}
		}

		// tempToReal: peta penanda sementara dari FE (temp_id) ke ID asli product_packages
		// yang baru dibuat. 0 selalu berarti paket anchor.
		tempToReal := map[int]int64{0: anchorID}
		for _, p := range req.Packages {
			refID, ok := tempToReal[p.RefTempID]
			if !ok {
				return fmt.Errorf("paket dengan temp_id %d merujuk paket yang belum dibuat (ref_temp_id %d)", p.TempID, p.RefTempID)
			}

			var pkgName *string
			if p.PackageName != "" {
				pkgName = &p.PackageName
			}
			if err := tx.Exec(insertProductPackageQuery,
				id, p.UnitID, pkgName, refID, p.Qty, p.RefQty, p.PurchasePrice, p.SellingPrice,
			).Error; err != nil {
				return err
			}
			var newID int64
			if err := tx.Raw(getLastProductInsertIDQuery).Scan(&newID).Error; err != nil {
				return err
			}
			tempToReal[p.TempID] = newID
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update mengubah data produk. PENTING (celah #1 / Aturan Operasional #5):
// req.Stock TIDAK LAGI ditulis langsung ke products.stock -- dulu ini jalan
// bypass yang bisa mengubah stok tanpa tercatat di stock_mutations sama
// sekali (dikonfirmasi lewat testing browser: field "Stok" di form edit bisa
// diedit bebas, 2 klik dari daftar produk, tanpa jejak audit). Sekarang
// selisih antara req.Stock dan stok anchor saat ini diterjemahkan jadi
// delta lewat ApplyStockDelta, supaya tercatat & tervalidasi (mis. tidak
// bisa dikurangi sampai minus) sama seperti jalur lain.
func (r *productRepo) Update(req *dto.UpdateRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time_helper.GetTimeNow()
		if err := tx.Exec(updateProductQuery,
			req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
			req.SellingPrice, req.MinStock, now, req.ID,
		).Error; err != nil {
			return err
		}

		if err := tx.Exec(updateAnchorPackagePrice, req.PurchasePrice, req.SellingPrice, now, req.ID).Error; err != nil {
			return err
		}

		var anchor struct {
			ID    int
			Stock float64
		}
		if err := tx.Raw(getAnchorPackageForUpdateQuery, req.ID).Scan(&anchor).Error; err != nil {
			return err
		}
		if anchor.ID == 0 {
			return fmt.Errorf("produk %d tidak punya baris satuan anchor", req.ID)
		}

		delta := req.Stock - anchor.Stock
		if delta == 0 {
			return nil
		}
		direction := model.StockIn
		qty := delta
		if delta < 0 {
			direction = model.StockOut
			qty = -delta
		}

		var userID *int
		if req.UserID > 0 {
			userID = &req.UserID
		}
		_, err := ApplyStockDelta(tx, ApplyStockDeltaParams{
			ProductID:     req.ID,
			PackageID:     anchor.ID,
			Quantity:      qty,
			Direction:     direction,
			MutationType:  "adjustment",
			ReferenceType: "product_edit",
			ReferenceID:   req.ID,
			Notes:         "Edit stok manual lewat form produk",
			UserID:        userID,
		})
		if err != nil {
			return WrapStockError(err, req.ID)
		}
		return nil
	})
}

func (r *productRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteProductQuery, req.ID).Error
	return err
}

func (r *productRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleProductStatusQuery, time_helper.GetTimeNow(), req.ID).Error
	return err
}

func (r *productRepo) MarkStockReviewed(id int) error {
	return r.db.Exec(markStockReviewedQuery, time_helper.GetTimeNow(), id).Error
}
