package repo

import (
	stderrors "errors"
	"fmt"

	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	dto "pos_api/domain/supplier_purchase/dto"
	model "pos_api/domain/supplier_purchase/model"
	"pos_api/errors"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	getPackagesByProductQuery           = `SELECT pp.id, pp.ref_package_id, pp.qty, pp.ref_qty, COALESCE(u.name, '') AS unit_name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.product_id = ?`
	generatePurchaseCodeQuery           = `SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(purchase_code, '-', -1) AS UNSIGNED)), 0) FROM purchases WHERE purchase_code LIKE CONCAT('PO-', ?, '-%')`
	createPurchaseQuery                 = `INSERT INTO purchases (purchase_code, invoice_number, supplier_id, purchase_date, discount_amount, total_amount, payment_status, paid_amount, remaining_amount, user_id, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createPurchaseItemQuery             = `INSERT INTO purchase_items (purchase_id, product_id, package_id, quantity, unit, conversion_qty, purchase_price, subtotal) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	createExpiryBatchQuery              = `INSERT INTO product_expiry_batches (product_id, purchase_item_id, package_id, qty, expired_date) VALUES (?, ?, ?, ?, ?)`
	getDefaultPackageIDQuery            = `SELECT id FROM product_packages WHERE product_id = ? AND is_default = 1 LIMIT 1`
	payPurchaseQuery                    = `UPDATE purchases SET paid_amount = paid_amount + ?, remaining_amount = remaining_amount - ?, payment_status = CASE WHEN remaining_amount <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
	getPurchaseItemsQuery               = `SELECT pi.id, pi.product_id, COALESCE(p.name, '') as product_name, pi.package_id, pi.quantity, pi.unit, COALESCE(pi.conversion_qty, 1) as conversion_qty, pi.purchase_price, pi.subtotal FROM purchase_items pi LEFT JOIN products p ON pi.product_id = p.id WHERE pi.purchase_id = ?`
	createPaymentQuery                  = `INSERT INTO purchase_payments (purchase_id, payment_date, amount, payment_method, notes, user_id) VALUES (?, ?, ?, ?, ?, ?)`
	getPaymentsQuery                    = `SELECT pp.id, pp.payment_date, pp.amount, COALESCE(pp.payment_method, '') as payment_method, COALESCE(pp.notes, '') as notes, COALESCE(u.full_name, '') as user_name, pp.created_at FROM purchase_payments pp LEFT JOIN users u ON pp.user_id = u.id WHERE pp.purchase_id = ? ORDER BY pp.created_at ASC`
	deleteStockMutationsQuery           = `DELETE FROM stock_mutations WHERE reference_type = 'purchase' AND reference_id = ?`
	deletePurchaseItemsQuery            = `DELETE FROM purchase_items WHERE purchase_id = ?`
	deletePurchaseQuery                 = `DELETE FROM purchases WHERE id = ?`
	getPurchaseByIDQuery                = `SELECT p.id, p.purchase_code, p.invoice_number, p.supplier_id, COALESCE(s.name, '') as supplier_name, p.purchase_date, p.discount_amount, p.total_amount, p.payment_status, p.paid_amount, p.remaining_amount, p.status, COALESCE(u.full_name, '') as user_name, p.notes FROM purchases p LEFT JOIN users u ON p.user_id = u.id LEFT JOIN suppliers s ON p.supplier_id = s.id WHERE p.id = ?`
	getRawPurchaseByIDQuery             = `SELECT id, purchase_code, invoice_number, supplier_id, purchase_date, discount_amount, total_amount, payment_status, paid_amount, remaining_amount, status, user_id, notes FROM purchases WHERE id = ?`
	getAllPurchasesBase                 = `SELECT p.id, p.purchase_code, p.invoice_number, p.supplier_id, COALESCE(s.name, '') as supplier_name, p.purchase_date, p.discount_amount, p.total_amount, p.payment_status, p.paid_amount, p.remaining_amount, p.status, COALESCE(u.full_name, '') as user_name, p.notes FROM purchases p LEFT JOIN users u ON p.user_id = u.id LEFT JOIN suppliers s ON p.supplier_id = s.id WHERE 1=1`
	countPurchasesBase                  = `SELECT COUNT(*) FROM purchases p WHERE 1=1`
	validatePaymentMethodQuery          = `SELECT COUNT(*) FROM payment_methods WHERE code = ? AND is_active = 1`
	getExpiryBatchesByPurchaseItemQuery = `SELECT qty, expired_date FROM product_expiry_batches WHERE purchase_item_id = ? ORDER BY expired_date ASC`

	getPurchaseForVoidQuery   = `SELECT status FROM purchases WHERE id = ? LIMIT 1 FOR UPDATE`
	voidPurchaseQuery         = `UPDATE purchases SET status = 'void', remaining_amount = 0, updated_at = NOW() WHERE id = ?`
	countReturnsByPurchaseQry = `SELECT COUNT(*) FROM supplier_returns WHERE purchase_id = ?`
	updatePurchaseTotalsQuery = `UPDATE purchases SET total_amount = ?, remaining_amount = ?, payment_status = CASE WHEN ? <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
)

func insertExpiryBatches(tx *gorm.DB, productID, purchaseItemID, packageID int, batches []dto.ExpiryBatchDraft) error {
	for _, b := range batches {
		if err := tx.Exec(createExpiryBatchQuery, productID, purchaseItemID, packageID, b.Qty, b.ExpiredDate).Error; err != nil {
			return err
		}
	}
	return nil
}

func resolveConversionQty(tx *gorm.DB, productID int, packageID *int, fallback float64) float64 {
	if fallback <= 0 {
		fallback = 1
	}
	if packageID == nil || *packageID <= 0 {
		return fallback
	}

	var pkgRows []*model_product.ProductPackage
	if err := tx.Raw(getPackagesByProductQuery, productID).Scan(&pkgRows).Error; err != nil {
		return fallback
	}

	factor, err := model_product.ResolvePackageFactor(pkgRows, *packageID)
	if err != nil || factor <= 0 {
		return fallback
	}
	return factor
}

func resolvePackageID(tx *gorm.DB, productID int, packageID *int) *int {
	if packageID != nil && *packageID > 0 {
		return packageID
	}
	var defaultID int
	if err := tx.Raw(getDefaultPackageIDQuery, productID).Scan(&defaultID).Error; err != nil || defaultID == 0 {
		return nil
	}
	return &defaultID
}

func calculateTotal(items []dto.PurchaseRequest, discountAmount float64) (subtotal float64, totalAmount float64) {
	for _, item := range items {
		subtotal += item.PurchasePrice * item.Quantity
	}
	totalAmount = subtotal - discountAmount
	if totalAmount < 0 {
		totalAmount = 0
	}
	return subtotal, totalAmount
}

func (r *purchaseRepo) GetAll(req *dto.GetAllRequest) ([]*model.PurchaseRow, int64, error) {
	var args []interface{}
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += " AND (p.invoice_number LIKE ? OR p.purchase_code LIKE ?)"
		args = append(args, search, search)
	}
	if req.StartDate != "" {
		conditions += " AND DATE(p.purchase_date) >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		conditions += " AND DATE(p.purchase_date) <= ?"
		args = append(args, req.EndDate)
	}
	if req.SupplierID != nil {
		conditions += " AND p.supplier_id = ?"
		args = append(args, *req.SupplierID)
	}
	if req.PaymentStatus != "" {
		conditions += " AND p.payment_status = ? AND p.status != 'void'"
		args = append(args, req.PaymentStatus)
	}

	var total int64
	if err := r.db.Raw(countPurchasesBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 20, 0)

	allowedSortFields := map[string]string{
		"purchase_date":  "p.purchase_date",
		"total_amount":   "p.total_amount",
		"supplier_name":  "s.name",
		"payment_status": "p.payment_status",
	}
	const defaultOrder = " ORDER BY p.purchase_date DESC, p.id DESC"

	query := getAllPurchasesBase + conditions + request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, defaultOrder) + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*model.PurchaseRow
	for rows.Next() {
		var item model.PurchaseRow
		if err := rows.Scan(
			&item.ID, &item.PurchaseCode, &item.InvoiceNumber, &item.SupplierID, &item.SupplierName,
			&item.PurchaseDate, &item.DiscountAmount, &item.TotalAmount, &item.PaymentStatus,
			&item.PaidAmount, &item.RemainingAmount, &item.Status, &item.UserName, &item.Notes,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
	}
	if items == nil {
		items = []*model.PurchaseRow{}
	}
	return items, total, nil
}

func (r *purchaseRepo) GetByID(id int) (*model.PurchaseRow, error) {
	rows, err := r.db.Raw(getPurchaseByIDQuery, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var item model.PurchaseRow
	if err := rows.Scan(
		&item.ID, &item.PurchaseCode, &item.InvoiceNumber, &item.SupplierID, &item.SupplierName,
		&item.PurchaseDate, &item.DiscountAmount, &item.TotalAmount, &item.PaymentStatus,
		&item.PaidAmount, &item.RemainingAmount, &item.Status, &item.UserName, &item.Notes,
	); err != nil {
		return nil, err
	}
	rows.Close()

	modelItems, err := r.GetItems(id)
	if err != nil {
		return nil, err
	}
	item.Items = modelItems

	return &item, nil
}

func (r *purchaseRepo) GetRawByID(id int) (*model.Purchase, error) {
	var p model.Purchase
	result := r.db.Raw(getRawPurchaseByIDQuery, id).Scan(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &p, nil
}

func (r *purchaseRepo) GetItems(purchaseID int) ([]model.PurchaseItem, error) {
	rows, err := r.db.Raw(getPurchaseItemsQuery, purchaseID).Rows()
	if err != nil {
		return nil, err
	}

	var items []model.PurchaseItem
	for rows.Next() {
		var item model.PurchaseItem
		if err := rows.Scan(
			&item.ID, &item.ProductID, &item.ProductName, &item.PackageID,
			&item.Quantity, &item.Unit, &item.ConversionQty, &item.PurchasePrice, &item.Subtotal,
		); err != nil {
			rows.Close()
			return nil, err
		}
		items = append(items, item)
	}
	rows.Close()

	for i := range items {
		batches, err := r.getExpiryBatchesByPurchaseItem(items[i].ID)
		if err != nil {
			return nil, err
		}
		items[i].ExpiryBatches = batches
	}

	return items, nil
}

func (r *purchaseRepo) getExpiryBatchesByPurchaseItem(purchaseItemID int) ([]model.PurchaseItemExpiryBatch, error) {
	var batches []model.PurchaseItemExpiryBatch
	if err := r.db.Raw(getExpiryBatchesByPurchaseItemQuery, purchaseItemID).Scan(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *purchaseRepo) GetPayments(purchaseID int) ([]model.PurchasePayment, error) {
	rows, err := r.db.Raw(getPaymentsQuery, purchaseID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.PurchasePayment
	for rows.Next() {
		var item model.PurchasePayment
		if err := rows.Scan(&item.ID, &item.PaymentDate, &item.Amount, &item.PaymentMethod, &item.Notes, &item.UserName, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []model.PurchasePayment{}
	}
	return items, nil
}

func (r *purchaseRepo) GenerateCode() (string, error) {
	todayCode := time_helper.GetTimeNow().Format("20060102")
	var maxSuffix int
	if err := r.db.Raw(generatePurchaseCodeQuery, todayCode).Scan(&maxSuffix).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("PO-%s-%03d", todayCode, maxSuffix+1), nil
}

func isDuplicatePurchaseCodeError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return stderrors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *purchaseRepo) Create(req *dto.CreateRequest) (*model.PurchaseRow, error) {
	var purchaseID int

	const maxCodeRetries = 5
	var err error
	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		err = r.createOnce(req, &purchaseID)
		if err == nil || !isDuplicatePurchaseCodeError(err) {
			break
		}
	}

	if err != nil {
		return nil, err
	}
	return r.GetByID(purchaseID)
}

func (r *purchaseRepo) createOnce(req *dto.CreateRequest, purchaseID *int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		todayCode := time_helper.GetTimeNow().Format("20060102")
		var maxSuffix int
		if err := tx.Raw(generatePurchaseCodeQuery, todayCode).Scan(&maxSuffix).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("PO-%s-%03d", todayCode, maxSuffix+1)

		_, totalAmount := calculateTotal(req.Items, req.DiscountAmount)

		paymentStatus := req.PaymentStatus
		if paymentStatus == "" {
			paymentStatus = "unpaid"
		}
		paidAmount := req.PaidAmount
		switch paymentStatus {
		case "paid":
			paidAmount = totalAmount
		case "unpaid":
			paidAmount = 0
		}
		remainingAmount := totalAmount - paidAmount

		if err := tx.Exec(createPurchaseQuery,
			code, req.InvoiceNumber, req.SupplierID, req.PurchaseDate,
			req.DiscountAmount, totalAmount, paymentStatus, paidAmount, remainingAmount, req.UserID, req.Notes,
		).Error; err != nil {
			return err
		}

		if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(purchaseID).Error; err != nil {
			return err
		}

		if paidAmount > 0 {
			paymentDate := req.PurchaseDate
			if paymentDate == "" {
				paymentDate = time_helper.GetTimeNow().Format("2006-01-02")
			}
			if err := tx.Exec(createPaymentQuery,
				*purchaseID, paymentDate, paidAmount, req.PaymentMethod, req.Notes, req.UserID,
			).Error; err != nil {
				return err
			}
		}

		for _, item := range req.Items {
			subtotal := item.PurchasePrice * item.Quantity
			conversionQty := resolveConversionQty(tx, item.ProductID, item.PackageID, item.ConversionQty)
			resolvedPackageID := resolvePackageID(tx, item.ProductID, item.PackageID)

			if err := tx.Exec(createPurchaseItemQuery,
				*purchaseID, item.ProductID, resolvedPackageID,
				item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, *resolvedPackageID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			if resolvedPackageID == nil {
				return fmt.Errorf("produk ID %d tidak punya paket satuan (package_id) yang bisa dipakai buat update stok", item.ProductID)
			}
			notes := fmt.Sprintf("Purchase Order %s", code)
			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     *resolvedPackageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockIn,
				MutationType:  "in",
				ReferenceType: "purchase",
				ReferenceID:   *purchaseID,
				Notes:         notes,
				UserID:        &req.UserID,
			}); err != nil {
				return wrapStockError(err, item.ProductID)
			}
		}

		return nil
	})
}

// wrapStockError menerjemahkan error teknis dari ApplyStockDelta jadi pesan
// yang enak dibaca user, tanpa membuang informasi aslinya.
func wrapStockError(err error, productID int) error {
	if stderrors.Is(err, model_product.ErrInsufficientStock) {
		return &errors.BadRequestError{Message: fmt.Sprintf("Stok produk ID %d tidak mencukupi untuk perubahan ini", productID)}
	}
	if stderrors.Is(err, model_product.ErrNeedsStockReview) {
		return &errors.BadRequestError{Message: fmt.Sprintf("Produk ID %d ditandai perlu ditinjau manual (needs_stock_review), operasi stok diblokir sampai ditinjau admin", productID)}
	}
	if stderrors.Is(err, model_product.ErrBranchingChain) {
		return &errors.BadRequestError{Message: fmt.Sprintf("Produk ID %d punya struktur satuan bercabang, tidak bisa diproses otomatis -- hubungi admin", productID)}
	}
	return err
}

func (r *purchaseRepo) Update(req *dto.UpdateRequest) (*model.PurchaseRow, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		oldItems, err := r.GetItems(req.ID)
		if err != nil {
			return err
		}

		_, totalAmount := calculateTotal(req.Items, req.DiscountAmount)

		paidAmount := req.PaidAmount
		if paidAmount < 0 {
			paidAmount = 0
		}
		remainingAmount := totalAmount - paidAmount
		if remainingAmount < 0 {
			return &errors.BadRequestError{Message: fmt.Sprintf(
				"Total baru (%.2f) tidak boleh lebih kecil dari yang sudah dibayar (%.2f)", totalAmount, paidAmount,
			)}
		}

		var paymentStatus string
		switch {
		case paidAmount <= 0:
			paymentStatus = "unpaid"
		case remainingAmount <= 0:
			paymentStatus = "paid"
		default:
			paymentStatus = "partial"
		}

		type stockKey struct {
			productID int
			packageID int
		}
		oldQty := make(map[stockKey]float64)
		for _, old := range oldItems {
			packageID := old.PackageID
			if packageID == nil {
				resolved := resolvePackageID(tx, old.ProductID, nil)
				packageID = resolved
			}
			if packageID == nil {
				return fmt.Errorf("produk ID %d (item lama) tidak punya paket satuan yang bisa dipakai buat balikkan stok", old.ProductID)
			}
			oldQty[stockKey{old.ProductID, *packageID}] += old.Quantity
		}

		if err := tx.Exec(deletePurchaseItemsQuery, req.ID).Error; err != nil {
			return err
		}

		paymentMethod := req.PaymentMethod
		if paymentMethod == "" {
			paymentMethod = "cash"
		}

		if err := tx.Exec(
			`UPDATE purchases SET invoice_number=?, supplier_id=?, purchase_date=?, discount_amount=?, total_amount=?, payment_status=?, paid_amount=?, payment_method=?, remaining_amount=?, notes=?, updated_at=NOW() WHERE id=?`,
			req.InvoiceNumber, req.SupplierID, req.PurchaseDate, req.DiscountAmount, totalAmount,
			paymentStatus, paidAmount, paymentMethod, remainingAmount, req.Notes, req.ID,
		).Error; err != nil {
			return err
		}

		newQty := make(map[stockKey]float64)
		for _, item := range req.Items {
			subtotal := item.PurchasePrice * item.Quantity
			conversionQty := resolveConversionQty(tx, item.ProductID, item.PackageID, item.ConversionQty)
			resolvedPackageID := resolvePackageID(tx, item.ProductID, item.PackageID)
			if resolvedPackageID == nil {
				return fmt.Errorf("produk ID %d tidak punya paket satuan (package_id) yang bisa dipakai buat update stok", item.ProductID)
			}
			if err := tx.Exec(createPurchaseItemQuery,
				req.ID, item.ProductID, resolvedPackageID,
				item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, *resolvedPackageID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			newQty[stockKey{item.ProductID, *resolvedPackageID}] += item.Quantity
		}

		// Terapkan hanya SELISIH bersih per (produk, paket), bukan reverse-semua-lalu-reapply-semua.
		// Kalau qty item lama & baru sama persis (item lain yang diubah/dihapus), delta = 0, stok
		// produk itu sama sekali tidak disentuh -- jadi tidak bisa gagal gara-gara stok riilnya
		// sudah berkurang oleh transaksi lain sejak PO ini dibuat.
		keys := make(map[stockKey]struct{}, len(oldQty)+len(newQty))
		for k := range oldQty {
			keys[k] = struct{}{}
		}
		for k := range newQty {
			keys[k] = struct{}{}
		}
		for k := range keys {
			delta := newQty[k] - oldQty[k]
			if delta == 0 {
				continue
			}
			direction := model_product.StockIn
			qty := delta
			notes := fmt.Sprintf("Edit PO ID %d -- penyesuaian stok (selisih item baru vs lama)", req.ID)
			if delta < 0 {
				direction = model_product.StockOut
				qty = -delta
				notes = fmt.Sprintf("Edit PO ID %d -- balikkan sebagian item lama (selisih)", req.ID)
			}
			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     k.productID,
				PackageID:     k.packageID,
				Quantity:      qty,
				Direction:     direction,
				MutationType:  "adjustment",
				ReferenceType: "purchase",
				ReferenceID:   req.ID,
				Notes:         notes,
				UserID:        &req.UserID,
			}); err != nil {
				return wrapStockError(err, k.productID)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetByID(req.ID)
}

// Delete menghapus PO permanen. TIDAK membalikkan stok di sini -- service
// layer (purchase_service.go Delete) mewajibkan PO sudah berstatus 'void'
// dulu sebelum boleh dihapus, dan Void() sudah membalikkan stok saat itu.
// Membalikkan lagi di sini akan jadi pengurangan dobel. (Sempat salah
// ditambahkan saat tinjauan Fase 4 poin 1, langsung dibatalkan setelah
// ketahuan ada guard "harus di-void dulu" di service layer.)
func (r *purchaseRepo) Delete(id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(deleteStockMutationsQuery, id).Error; err != nil {
			return err
		}

		if err := tx.Exec(deletePurchaseItemsQuery, id).Error; err != nil {
			return err
		}

		return tx.Exec(deletePurchaseQuery, id).Error
	})
}

func (r *purchaseRepo) IsValidPaymentMethod(code string) (bool, error) {
	var count int
	if err := r.db.Raw(validatePaymentMethodQuery, code).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *purchaseRepo) Pay(req *dto.PayRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(payPurchaseQuery, req.Amount, req.Amount, req.ID).Error; err != nil {
			return err
		}

		paymentDate := req.PaymentDate
		if paymentDate == "" {
			paymentDate = time_helper.GetTimeNow().Format("2006-01-02")
		}
		return tx.Exec(createPaymentQuery, req.ID, paymentDate, req.Amount, req.PaymentMethod, req.Notes, req.UserID).Error
	})
}

func (r *purchaseRepo) CountReturnsByPurchaseID(purchaseID int) (int64, error) {
	var count int64
	if err := r.db.Raw(countReturnsByPurchaseQry, purchaseID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *purchaseRepo) Void(id int, userID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lockData struct {
			Status string
		}
		if err := tx.Raw(getPurchaseForVoidQuery, id).Scan(&lockData).Error; err != nil {
			return err
		}

		if lockData.Status == "void" {
			return &errors.BadRequestError{Message: "PO sudah di-void"}
		}

		if err := tx.Exec(voidPurchaseQuery, id).Error; err != nil {
			return err
		}

		rows, err := tx.Raw(getPurchaseItemsQuery, id).Rows()
		if err != nil {
			return err
		}
		var items []model.PurchaseItem
		for rows.Next() {
			var item model.PurchaseItem
			if err := rows.Scan(
				&item.ID, &item.ProductID, &item.ProductName, &item.PackageID,
				&item.Quantity, &item.Unit, &item.ConversionQty, &item.PurchasePrice, &item.Subtotal,
			); err != nil {
				rows.Close()
				return err
			}
			items = append(items, item)
		}
		rows.Close()

		notes := fmt.Sprintf("Void purchase order ID %d", id)
		for _, item := range items {
			packageID := item.PackageID
			if packageID == nil {
				resolved := resolvePackageID(tx, item.ProductID, nil)
				packageID = resolved
			}
			if packageID == nil {
				return fmt.Errorf("produk ID %d tidak punya paket satuan yang bisa dipakai buat void", item.ProductID)
			}
			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     *packageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockOut,
				MutationType:  "void_purchase",
				ReferenceType: "purchase",
				ReferenceID:   id,
				Notes:         notes,
				UserID:        &userID,
			}); err != nil {
				return wrapStockError(err, item.ProductID)
			}
		}

		return nil
	})
}

func (r *purchaseRepo) AddItems(req *dto.AddItemsRequest) (*model.PurchaseRow, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var current struct {
			TotalAmount float64
			PaidAmount  float64
		}
		if err := tx.Raw(`SELECT total_amount, paid_amount FROM purchases WHERE id = ? LIMIT 1 FOR UPDATE`, req.ID).Scan(&current).Error; err != nil {
			return err
		}

		_, newItemsTotal := calculateTotal(req.Items, 0)

		for _, item := range req.Items {
			subtotal := item.PurchasePrice * item.Quantity
			conversionQty := resolveConversionQty(tx, item.ProductID, item.PackageID, item.ConversionQty)
			resolvedPackageID := resolvePackageID(tx, item.ProductID, item.PackageID)

			if err := tx.Exec(createPurchaseItemQuery,
				req.ID, item.ProductID, resolvedPackageID, item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, *resolvedPackageID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			if resolvedPackageID == nil {
				return fmt.Errorf("produk ID %d tidak punya paket satuan (package_id) yang bisa dipakai buat update stok", item.ProductID)
			}
			notes := fmt.Sprintf("Tambah item PO ID %d", req.ID)
			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     *resolvedPackageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockIn,
				MutationType:  "in",
				ReferenceType: "purchase",
				ReferenceID:   req.ID,
				Notes:         notes,
				UserID:        &req.UserID,
			}); err != nil {
				return wrapStockError(err, item.ProductID)
			}
		}

		newTotal := current.TotalAmount + newItemsTotal
		newRemaining := newTotal - current.PaidAmount

		if err := tx.Exec(updatePurchaseTotalsQuery, newTotal, newRemaining, newRemaining, req.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetByID(req.ID)
}
