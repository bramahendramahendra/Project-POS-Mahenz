package repo

import (
	"fmt"

	dto "pos_api/domain/stock_reconciliation/dto"
	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	custom_errors "pos_api/errors"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

// ============================================================================
// Konsep: "stok baru" = jumlah product_packages.stock tiap level yang
// dikonversi ke satuan dasar (anchor). "stok lama" = products_stock_backup.stock
// (tabel cadangan yang dibuat saat migrasi; BISA JADI tidak ada di sebagian
// environment). Breakdown "dari mana stok baru dihitung" diambil dari ledger
// stock_mutations (SUM per mutation_type, dalam satuan dasar).
// ============================================================================

const (
	// stok baru per produk = SUM stok tiap paket dikonversi ke anchor.
	// Konversi dihitung di Go (ResolvePackageFactor) supaya konsisten dengan
	// jalur stok live, bukan reimplementasi SQL.
	baseUnitNameQuery = `
		SELECT u.name
		FROM product_packages pp
		JOIN units u ON u.id = pp.unit_id
		WHERE pp.product_id = ? AND pp.is_default = 1
		LIMIT 1`

	packagesQuery = `
		SELECT pp.id, pp.product_id, pp.unit_id, u.name AS unit_name,
		       COALESCE(u.abbreviation, '') AS abbreviation,
		       COALESCE(pp.package_name, '') AS package_name,
		       pp.ref_package_id, pp.qty, pp.ref_qty,
		       pp.purchase_price, pp.selling_price, pp.is_default,
		       pp.stock, pp.reserved_qty, pp.is_active,
		       COALESCE(u.is_continuous, 0) AS is_continuous
		FROM product_packages pp
		JOIN units u ON u.id = pp.unit_id
		WHERE pp.product_id = ? AND pp.is_active = 1`

	breakdownQuery = `
		SELECT mutation_type, COALESCE(SUM(quantity), 0) AS total
		FROM stock_mutations
		WHERE product_id = ?
		GROUP BY mutation_type`

	kartuStokQuery = `
		SELECT sm.id, sm.mutation_type, sm.quantity, sm.stock_before, sm.stock_after,
		       COALESCE(sm.reference_type, '') AS reference_type,
		       COALESCE(sm.reference_id, 0) AS reference_id,
		       COALESCE(sm.notes, '') AS notes,
		       COALESCE(u.full_name, '') AS user_name, sm.created_at
		FROM stock_mutations sm
		LEFT JOIN users u ON sm.user_id = u.id
		WHERE sm.product_id = ?
		ORDER BY sm.created_at DESC, sm.id DESC`

	updateReviewFlagQuery = `UPDATE products SET needs_stock_review = ?, stock_review_note = ?, updated_at = ? WHERE id = ?`
)

// backupTableExists cek apakah tabel products_stock_backup tersedia.
func (r *stockReconciliationRepo) backupTableExists() bool {
	var name string
	err := r.db.Raw(`SHOW TABLES LIKE 'products_stock_backup'`).Scan(&name).Error
	return err == nil && name != ""
}

// oldStockOf mengambil stok lama 1 produk dari products_stock_backup.
// Mengembalikan (stok, tersedia). tersedia=false kalau tabel/baris tidak ada.
func (r *stockReconciliationRepo) oldStockOf(productID int, backupExists bool) (float64, bool) {
	if !backupExists {
		return 0, false
	}
	var stock *float64
	if err := r.db.Raw(`SELECT stock FROM products_stock_backup WHERE id = ?`, productID).Scan(&stock).Error; err != nil {
		return 0, false
	}
	if stock == nil {
		return 0, false
	}
	return *stock, true
}

// newStockOf menghitung stok baru 1 produk = SUM(stok paket -> anchor).
func newStockOf(packages []*model_product.ProductPackage) (float64, error) {
	total := 0.0
	for _, p := range packages {
		factor, err := model_product.ResolvePackageFactor(packages, p.ID)
		if err != nil {
			return 0, err
		}
		total += p.Stock * factor
	}
	return total, nil
}

func (r *stockReconciliationRepo) loadPackages(productID int) ([]*model_product.ProductPackage, error) {
	var packages []*model_product.ProductPackage
	if err := r.db.Raw(packagesQuery, productID).Scan(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

// ---------------------------------------------------------------------------
// GetList
// ---------------------------------------------------------------------------
func (r *stockReconciliationRepo) GetList(req *dto.ListRequest) ([]*dto.ListItem, int64, error) {
	backupExists := r.backupTableExists()

	// Ambil daftar produk (dengan filter) — hitung stok baru per produk di Go
	// supaya konversi antar level konsisten dengan jalur live.
	base := `
		SELECT p.id, COALESCE(p.sku, '') AS product_code, p.name,
		       p.needs_stock_review, COALESCE(p.stock_review_note, '') AS stock_review_note
		FROM products p
		WHERE p.is_active = 1`
	countBase := `SELECT COUNT(*) FROM products p WHERE p.is_active = 1`

	var args []any
	var countArgs []any
	cond := ""
	if req.Search != "" {
		cond += ` AND (p.name LIKE ? OR p.sku LIKE ? OR p.barcode LIKE ?)`
		like := "%" + req.Search + "%"
		args = append(args, like, like, like)
		countArgs = append(countArgs, like, like, like)
	}
	if req.CategoryID != nil {
		cond += ` AND p.category_id = ?`
		args = append(args, *req.CategoryID)
		countArgs = append(countArgs, *req.CategoryID)
	}
	if req.OnlyReview {
		cond += ` AND p.needs_stock_review = 1`
	}

	var total int64
	if err := r.db.Raw(countBase+cond, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
	query := base + cond + ` ORDER BY p.needs_stock_review DESC, p.name ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	type row struct {
		ID               int
		ProductCode      string
		Name             string
		NeedsStockReview bool
		StockReviewNote  string
	}
	var rows []row
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*dto.ListItem, 0, len(rows))
	for _, rw := range rows {
		packages, err := r.loadPackages(rw.ID)
		if err != nil {
			return nil, 0, err
		}
		newStock, err := newStockOf(packages)
		if err != nil {
			// produk dengan struktur satuan bermasalah tetap ditampilkan,
			// stok baru dianggap 0 dan biarkan admin meninjau.
			newStock = 0
		}
		baseUnit := ""
		for _, p := range packages {
			if p.IsDefault {
				baseUnit = p.UnitName
				break
			}
		}
		oldStock, oldAvail := r.oldStockOf(rw.ID, backupExists)

		item := &dto.ListItem{
			ProductID:         rw.ID,
			ProductCode:       rw.ProductCode,
			ProductName:       rw.Name,
			BaseUnit:          baseUnit,
			OldStock:          oldStock,
			OldStockAvailable: oldAvail,
			NewStock:          newStock,
			NeedsStockReview:  rw.NeedsStockReview,
			StockReviewNote:   rw.StockReviewNote,
		}
		if oldAvail {
			item.Diff = newStock - oldStock
		}
		items = append(items, item)
	}

	return items, total, nil
}

// ---------------------------------------------------------------------------
// GetSummary
// ---------------------------------------------------------------------------
func (r *stockReconciliationRepo) GetSummary() (*dto.SummaryResponse, error) {
	out := &dto.SummaryResponse{BackupAvailable: r.backupTableExists()}

	if err := r.db.Raw(`SELECT COUNT(*) FROM products WHERE is_active = 1`).Scan(&out.TotalProducts).Error; err != nil {
		return nil, err
	}
	if err := r.db.Raw(`SELECT COUNT(*) FROM products WHERE is_active = 1 AND needs_stock_review = 1`).Scan(&out.NeedsReview).Error; err != nil {
		return nil, err
	}
	out.Matched = out.TotalProducts - out.NeedsReview
	return out, nil
}

// ---------------------------------------------------------------------------
// GetDetail
// ---------------------------------------------------------------------------
func (r *stockReconciliationRepo) GetDetail(productID int) (*dto.DetailResponse, error) {
	type head struct {
		ID               int
		ProductCode      string
		Name             string
		NeedsStockReview bool
		StockReviewNote  string
	}
	var h head
	err := r.db.Raw(`
		SELECT id, COALESCE(sku, '') AS product_code, name,
		       needs_stock_review, COALESCE(stock_review_note, '') AS stock_review_note
		FROM products WHERE id = ?`, productID).Scan(&h).Error
	if err != nil {
		return nil, err
	}
	if h.ID == 0 {
		return nil, &custom_errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	packages, err := r.loadPackages(productID)
	if err != nil {
		return nil, err
	}

	baseUnit := ""
	for _, p := range packages {
		if p.IsDefault {
			baseUnit = p.UnitName
			break
		}
	}

	newStock, err := newStockOf(packages)
	if err != nil {
		newStock = 0
	}

	backupExists := r.backupTableExists()
	oldStock, oldAvail := r.oldStockOf(productID, backupExists)

	resp := &dto.DetailResponse{
		ProductID:         h.ID,
		ProductCode:       h.ProductCode,
		ProductName:       h.Name,
		BaseUnit:          baseUnit,
		OldStock:          oldStock,
		OldStockAvailable: oldAvail,
		NewStock:          newStock,
		NeedsStockReview:  h.NeedsStockReview,
		StockReviewNote:   h.StockReviewNote,
	}
	if oldAvail {
		resp.Diff = newStock - oldStock
	}

	// breakdown per jenis mutasi (satuan dasar)
	resp.Breakdown, err = r.buildBreakdown(productID)
	if err != nil {
		return nil, err
	}

	// level satuan + faktor ke dasar + stok sekarang
	for _, p := range packages {
		factor, ferr := model_product.ResolvePackageFactor(packages, p.ID)
		if ferr != nil {
			factor = 0
		}
		resp.Packages = append(resp.Packages, dto.PackageLevel{
			PackageID:    p.ID,
			UnitName:     p.UnitName,
			PackageName:  p.PackageName,
			IsDefault:    p.IsDefault,
			FactorToBase: factor,
			CurrentStock: p.Stock,
		})
	}

	// kartu stok (riwayat mutasi)
	var kartu []dto.KartuStokRow
	if err := r.db.Raw(kartuStokQuery, productID).Scan(&kartu).Error; err != nil {
		return nil, err
	}
	resp.KartuStok = kartu

	return resp, nil
}

// buildBreakdown menjumlahkan stock_mutations per jenis (satuan dasar) dan
// menyusun rumus stok baru. Angka SUM(quantity) di ledger selalu positif per
// baris; tanda (+/-) ditentukan dari jenis mutasinya.
func (r *stockReconciliationRepo) buildBreakdown(productID int) (dto.Breakdown, error) {
	type row struct {
		MutationType string
		Total        float64
	}
	var rows []row
	if err := r.db.Raw(breakdownQuery, productID).Scan(&rows).Error; err != nil {
		return dto.Breakdown{}, err
	}

	var b dto.Breakdown
	for _, rw := range rows {
		switch rw.MutationType {
		case "in":
			b.PurchaseIn = rw.Total
		case "void_purchase":
			b.PurchaseVoid = rw.Total
		case "out":
			b.SaleOut = rw.Total
		case "void":
			b.SaleVoid = rw.Total
		case "return":
			b.SupplierReturn = rw.Total
		case "expired":
			b.Expired = rw.Total
		case "adjustment":
			b.Adjustment = rw.Total
		}
	}
	// rumus: +in -void_purchase -out +void -return -expired +adjustment
	b.ComputedNewStock = b.PurchaseIn - b.PurchaseVoid - b.SaleOut + b.SaleVoid - b.SupplierReturn - b.Expired + b.Adjustment
	return b, nil
}

// ---------------------------------------------------------------------------
// Adjust — koreksi manual stok per level satuan.
// ---------------------------------------------------------------------------
func (r *stockReconciliationRepo) Adjust(req *dto.AdjustRequest) (*dto.DetailResponse, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Muat level satuan produk untuk validasi kepemilikan & stok sekarang.
		var packages []*model_product.ProductPackage
		if err := tx.Raw(packagesQuery, req.ProductID).Scan(&packages).Error; err != nil {
			return err
		}
		if len(packages) == 0 {
			return &custom_errors.BadRequestError{Message: "Produk tidak punya level satuan aktif"}
		}
		byID := make(map[int]*model_product.ProductPackage, len(packages))
		for _, p := range packages {
			byID[p.ID] = p
		}

		// Produk yang perlu ditinjau memblokir ApplyStockDelta (ErrNeedsStockReview).
		// Koreksi manual ini justru untuk MENYELESAIKAN tinjauan, jadi flag
		// di-clear DULU sebelum menerapkan delta, lalu diisi catatan koreksi.
		now := time_helper.GetTimeNow()
		clearNote := fmt.Sprintf("dikoreksi manual oleh user %d", req.UserID)
		if req.Note != "" {
			clearNote = req.Note + " (" + clearNote + ")"
		}
		if err := tx.Exec(updateReviewFlagQuery, false, clearNote, now, req.ProductID).Error; err != nil {
			return err
		}

		for _, lvl := range req.Levels {
			p, ok := byID[lvl.PackageID]
			if !ok {
				return &custom_errors.BadRequestError{Message: fmt.Sprintf("Level satuan %d bukan milik produk ini", lvl.PackageID)}
			}

			diff := lvl.NewStock - p.Stock
			if diff == 0 {
				continue // tidak berubah, lewati
			}

			direction := model_product.StockIn
			qty := diff
			if diff < 0 {
				direction = model_product.StockOut
				qty = -diff
			}

			// qty di sini dalam satuan level itu sendiri (PackageID), sesuai
			// kontrak ApplyStockDelta (celah #21). Selisih stok level = selisih
			// jumlah satuan level tsb.
			userID := req.UserID
			_, derr := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     req.ProductID,
				PackageID:     lvl.PackageID,
				Quantity:      qty,
				Direction:     direction,
				MutationType:  "adjustment",
				ReferenceType: "stock_reconciliation",
				ReferenceID:   req.ProductID,
				Notes:         clearNote,
				UserID:        &userID,
			})
			if derr != nil {
				return product_repo.WrapStockError(derr, req.ProductID)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetDetail(req.ProductID)
}

// ---------------------------------------------------------------------------
// MarkReviewed — tandai selesai tanpa mengubah stok.
// ---------------------------------------------------------------------------
func (r *stockReconciliationRepo) MarkReviewed(req *dto.MarkReviewedRequest) (*dto.DetailResponse, error) {
	now := time_helper.GetTimeNow()
	note := fmt.Sprintf("ditandai sudah ditinjau oleh user %d", req.UserID)
	if req.Note != "" {
		note = req.Note + " (" + note + ")"
	}
	if err := r.db.Exec(updateReviewFlagQuery, false, note, now, req.ProductID).Error; err != nil {
		return nil, err
	}
	return r.GetDetail(req.ProductID)
}
