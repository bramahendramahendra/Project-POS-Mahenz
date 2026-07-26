package repo

import (
	stderrors "errors"
	"fmt"
	"time"

	model_product "pos_api/domain/product/model"
	dto "pos_api/domain/supplier_purchase/dto"
	model "pos_api/domain/supplier_purchase/model"
	"pos_api/errors"
	request_helper "pos_api/helper/request"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	getPackagesByProductQuery           = `SELECT pp.id, pp.ref_package_id, pp.qty, pp.ref_qty, COALESCE(u.name, '') AS unit_name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.product_id = ?`
	generatePurchaseCodeQuery           = `SELECT COUNT(*) FROM purchases WHERE DATE(purchase_date) = ?`
	createPurchaseQuery                 = `INSERT INTO purchases (purchase_code, invoice_number, supplier_id, purchase_date, discount_amount, total_amount, payment_status, paid_amount, remaining_amount, user_id, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createPurchaseItemQuery             = `INSERT INTO purchase_items (purchase_id, product_id, quantity, unit, conversion_qty, purchase_price, subtotal) VALUES (?, ?, ?, ?, ?, ?, ?)`
	addStockQuery                       = `UPDATE products SET stock = stock + ?, updated_at = NOW() WHERE id = ?`
	payPurchaseQuery                    = `UPDATE purchases SET paid_amount = paid_amount + ?, remaining_amount = remaining_amount - ?, payment_status = CASE WHEN remaining_amount <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
	getPurchaseItemsQuery               = `SELECT pi.id, pi.product_id, COALESCE(p.name, '') as product_name, pi.quantity, pi.unit, COALESCE(pi.conversion_qty, 1) as conversion_qty, pi.purchase_price, pi.subtotal FROM purchase_items pi LEFT JOIN products p ON pi.product_id = p.id WHERE pi.purchase_id = ?`
	createPaymentQuery                  = `INSERT INTO purchase_payments (purchase_id, payment_date, amount, payment_method, notes, user_id) VALUES (?, ?, ?, ?, ?, ?)`
	getPaymentsQuery                    = `SELECT pp.id, pp.payment_date, pp.amount, COALESCE(pp.payment_method, '') as payment_method, COALESCE(pp.notes, '') as notes, COALESCE(u.full_name, '') as user_name, pp.created_at FROM purchase_payments pp LEFT JOIN users u ON pp.user_id = u.id WHERE pp.purchase_id = ? ORDER BY pp.created_at ASC`
	rollbackStockQuery                  = `UPDATE products SET stock = stock - ?, updated_at = NOW() WHERE id = ?`
	deleteStockMutationsQuery           = `DELETE FROM stock_mutations WHERE reference_type = 'purchase' AND reference_id = ?`
	deletePurchaseItemsQuery            = `DELETE FROM purchase_items WHERE purchase_id = ?`
	deletePurchaseQuery                 = `DELETE FROM purchases WHERE id = ?`
	getPurchaseByIDQuery                = `SELECT p.id, p.purchase_code, p.invoice_number, p.supplier_id, COALESCE(s.name, '') as supplier_name, p.purchase_date, p.discount_amount, p.total_amount, p.payment_status, p.paid_amount, p.remaining_amount, p.status, COALESCE(u.full_name, '') as user_name, p.notes FROM purchases p LEFT JOIN users u ON p.user_id = u.id LEFT JOIN suppliers s ON p.supplier_id = s.id WHERE p.id = ?`
	getRawPurchaseByIDQuery             = `SELECT id, purchase_code, invoice_number, supplier_id, purchase_date, discount_amount, total_amount, payment_status, paid_amount, remaining_amount, status, user_id, notes FROM purchases WHERE id = ?`
	getAllPurchasesBase                 = `SELECT p.id, p.purchase_code, p.invoice_number, p.supplier_id, COALESCE(s.name, '') as supplier_name, p.purchase_date, p.discount_amount, p.total_amount, p.payment_status, p.paid_amount, p.remaining_amount, p.status, COALESCE(u.full_name, '') as user_name, p.notes FROM purchases p LEFT JOIN users u ON p.user_id = u.id LEFT JOIN suppliers s ON p.supplier_id = s.id WHERE 1=1`
	countPurchasesBase                  = `SELECT COUNT(*) FROM purchases p WHERE 1=1`
	createStockMutationQuery            = `INSERT INTO stock_mutations (product_id, mutation_type, quantity, stock_before, stock_after, reference_type, reference_id, notes, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getProductStockQuery                = `SELECT stock FROM products WHERE id = ? LIMIT 1`
	validatePaymentMethodQuery          = `SELECT COUNT(*) FROM payment_methods WHERE code = ? AND is_active = 1`
	createExpiryBatchQuery              = `INSERT INTO product_expiry_batches (product_id, purchase_item_id, qty, expired_date) VALUES (?, ?, ?, ?)`
	getExpiryBatchesByPurchaseItemQuery = `SELECT qty, expired_date FROM product_expiry_batches WHERE purchase_item_id = ? ORDER BY expired_date ASC`

	getPurchaseForVoidQuery = `SELECT status FROM purchases WHERE id = ? LIMIT 1 FOR UPDATE`
	// voidPurchaseQuery: remaining_amount digugurkan ke 0 karena void artinya "tidak
	// ada lagi yang perlu ditagih/dibayar sejak sekarang". paid_amount & payment_status
	// SENGAJA tidak disentuh — itu fakta historis (mis. "PO ini sempat Lunas sebelum
	// dibatalkan"), tetap ditampilkan berdampingan dengan badge void di detail PO. Kalau
	// PO sudah pernah dibayar sebagian/lunas, uang yang sudah keluar TIDAK otomatis
	// di-refund oleh sistem ini — itu proses manual di luar sistem.
	voidPurchaseQuery         = `UPDATE purchases SET status = 'void', remaining_amount = 0, updated_at = NOW() WHERE id = ?`
	countReturnsByPurchaseQry = `SELECT COUNT(*) FROM supplier_returns WHERE purchase_id = ?`
	updatePurchaseTotalsQuery = `UPDATE purchases SET total_amount = ?, remaining_amount = ?, payment_status = CASE WHEN ? <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
)

// insertExpiryBatches menyimpan rincian tanggal expired (opsional) untuk 1 baris item
// PO yang baru saja di-insert — dipanggil setelah ID baris purchase_items itu didapat.
// Sum qty-nya sudah divalidasi sama dengan qty item di service (validateExpiryBatches),
// di sini tinggal insert apa adanya.
func insertExpiryBatches(tx *gorm.DB, productID, purchaseItemID int, batches []dto.ExpiryBatchDraft) error {
	for _, b := range batches {
		if err := tx.Exec(createExpiryBatchQuery, productID, purchaseItemID, b.Qty, b.ExpiredDate).Error; err != nil {
			return err
		}
	}
	return nil
}

// resolveConversionQty menghitung ulang faktor konversi dari server berdasarkan PackageID
// yang dipilih user (menelusuri rantai ref_package_id sampai ke paket anchor), bukan
// percaya ConversionQty yang dikirim klien — supaya rasio yang sudah berubah di produk
// tidak bisa dieksploitasi lewat payload yang usang. ConversionQty klien cuma jadi
// fallback kalau PackageID tidak dikirim (kompatibilitas alur lama).
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

// calculateTotal menghitung subtotal & total dari daftar item pembelian.
// Diekstrak dari Create/Update supaya tidak duplikasi logic (dan dipakai juga oleh AddItems).
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
		conditions += " AND p.payment_status = ?"
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
			&item.ID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.Unit, &item.ConversionQty, &item.PurchasePrice, &item.Subtotal,
		); err != nil {
			rows.Close()
			return nil, err
		}
		items = append(items, item)
	}
	rows.Close()

	// Rincian expired diambil per item SETELAH rows di atas ditutup (bukan di dalam loop
	// yang sama) supaya tidak ada query bersarang di koneksi yang masih dipakai rows.Next().
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
	today := time.Now().Format("2006-01-02")
	var count int
	if err := r.db.Raw(generatePurchaseCodeQuery, today).Scan(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("PO-%s-%03d", time.Now().Format("20060102"), count+1), nil
}

// isDuplicatePurchaseCodeError mendeteksi MySQL error 1062 (duplicate entry) pada
// kolom purchase_code — bisa terjadi kalau 2 request Create() berjalan bersamaan
// dan sama-sama menghitung count+1 yang sama sebelum salah satunya sempat commit
// (generatePurchaseCodeQuery pakai COUNT(*), bukan sequence atomik). Dulu error ini
// lolos mentah sebagai 500 Internal Server Error ke klien; sekarang di-retry otomatis.
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
		today := time.Now().Format("2006-01-02")
		var count int
		if err := tx.Raw(generatePurchaseCodeQuery, today).Scan(&count).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("PO-%s-%03d", time.Now().Format("20060102"), count+1)

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
				paymentDate = time.Now().Format("2006-01-02")
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
			stockAdd := item.Quantity * conversionQty

			if err := tx.Exec(createPurchaseItemQuery,
				*purchaseID, item.ProductID,
				item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			var stockBefore float64
			if err := tx.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
				return err
			}

			if err := tx.Exec(addStockQuery, stockAdd, item.ProductID).Error; err != nil {
				return err
			}

			stockAfter := stockBefore + stockAdd
			notes := fmt.Sprintf("Purchase Order %s", code)
			if err := tx.Exec(createStockMutationQuery,
				item.ProductID, "in", stockAdd, stockBefore, stockAfter,
				"purchase", *purchaseID, notes, req.UserID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update mengedit PO — item/qty/harga/diskon boleh diubah bebas terlepas dari status
// pembayaran (lihat purchase_service.go Update()). Urutan kerja SENGAJA: hitung &
// validasi semuanya dulu (total baru vs paid_amount, dan keamanan stok tiap produk
// yang dihapus/diganti dari daftar) SEBELUM satupun query mutasi (rollback stok,
// hapus item lama, insert item baru) dieksekusi — supaya kalau validasi gagal,
// transaksi di-abort tanpa mengubah apapun sama sekali, bukan "setengah ke-update".
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

		// payment_status akhir TIDAK dipercaya dari req.PaymentStatus — dihitung ulang murni
		// dari paidAmount vs totalAmount, supaya PO "Lunas" yang totalnya naik otomatis
		// balik jadi "Bayar Sebagian", dan tidak ada celah payment_status tidak sinkron
		// dengan angka aslinya.
		var paymentStatus string
		switch {
		case paidAmount <= 0:
			paymentStatus = "unpaid"
		case remainingAmount <= 0:
			paymentStatus = "paid"
		default:
			paymentStatus = "partial"
		}

		// Validasi stok: untuk tiap item LAMA yang produknya sudah tidak ada lagi di
		// daftar item BARU (dihapus atau diganti ke produk lain) — bukan yang cuma
		// qty-nya diubah — pastikan rollback qty pembelian asli tidak bikin stok
		// produk itu jadi minus. Minus berarti sebagian stok dari pembelian ini sudah
		// terjual/keluar lewat jalur lain (mis. Kasir), jadi tidak bisa "ditarik balik".
		newProductIDs := make(map[int]bool, len(req.Items))
		for _, item := range req.Items {
			newProductIDs[item.ProductID] = true
		}
		for _, old := range oldItems {
			if newProductIDs[old.ProductID] {
				continue
			}
			convQty := old.ConversionQty
			if convQty <= 0 {
				convQty = 1
			}
			originalStockQty := old.Quantity * convQty

			var currentStock float64
			if err := tx.Raw(getProductStockQuery, old.ProductID).Scan(&currentStock).Error; err != nil {
				return err
			}
			if currentStock-originalStockQty < 0 {
				return &errors.BadRequestError{Message: fmt.Sprintf(
					"Produk '%s' tidak bisa diganti/dihapus dari pembelian ini karena sebagian stoknya sudah terjual (sisa stok %.3f, butuh rollback %.3f)",
					old.ProductName, currentStock, originalStockQty,
				)}
			}
		}

		// Validasi lolos semua — baru sekarang boleh mulai mutasi.
		for _, item := range oldItems {
			convQty := item.ConversionQty
			if convQty <= 0 {
				convQty = 1
			}
			if err := tx.Exec(rollbackStockQuery, item.Quantity*convQty, item.ProductID).Error; err != nil {
				return err
			}
		}

		if err := tx.Exec(deletePurchaseItemsQuery, req.ID).Error; err != nil {
			return err
		}

		// payment_method punya FK constraint ke payment_methods(code) — kolom ini NOT NULL
		// DEFAULT 'cash' di skema, tapi UPDATE eksplisit tidak otomatis jatuh ke default itu
		// seperti INSERT yang mengosongkan kolom. Default manual di sini supaya konsisten
		// dengan Create() dan tidak pernah kirim string kosong ke kolom ber-FK.
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

		for _, item := range req.Items {
			subtotal := item.PurchasePrice * item.Quantity
			conversionQty := resolveConversionQty(tx, item.ProductID, item.PackageID, item.ConversionQty)
			if err := tx.Exec(createPurchaseItemQuery,
				req.ID, item.ProductID,
				item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			if err := tx.Exec(addStockQuery, item.Quantity*conversionQty, item.ProductID).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetByID(req.ID)
}

// Delete membersihkan PO secara permanen — TIDAK melakukan rollback stok karena
// service.Delete menjamin PO yang sampai ke sini sudah berstatus "void", yang
// artinya stoknya sudah pernah dikembalikan oleh Void(). Rollback kedua di sini
// akan mengurangi stok dua kali untuk item yang sama.
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
			paymentDate = time.Now().Format("2006-01-02")
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

// Void mengunci baris PO (FOR UPDATE), menandai status void beserta menggugurkan
// remaining_amount ke 0, lalu mengembalikan stok tiap item ke produk dan mencatat
// mutasi stok baru (tidak menghapus histori lama), meniru pola transaction_repo.go
// Void. Boleh dipanggil untuk PO di status pembayaran manapun (unpaid/partial/paid) —
// lihat komentar di voidPurchaseQuery soal kenapa paid_amount/payment_status sengaja
// tidak diubah.
func (r *purchaseRepo) Void(id int, userID int) error {
	var lockData struct {
		Status string
	}
	if err := r.db.Raw(getPurchaseForVoidQuery, id).Scan(&lockData).Error; err != nil {
		return err
	}

	// Re-cek di dalam row lock supaya aman dari race condition (mis. dua klik void
	// bersamaan): request kedua baru bisa lanjut setelah request pertama commit,
	// dan pada titik itu status sudah 'void' sehingga ditolak di sini.
	if lockData.Status == "void" {
		return &errors.BadRequestError{Message: "PO sudah di-void"}
	}

	if err := r.db.Exec(voidPurchaseQuery, id).Error; err != nil {
		return err
	}

	items, err := r.GetItems(id)
	if err != nil {
		return err
	}

	for _, item := range items {
		convQty := item.ConversionQty
		if convQty <= 0 {
			convQty = 1
		}
		stockRestore := item.Quantity * convQty

		var stockBefore float64
		if err := r.db.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
			return err
		}

		if err := r.db.Exec(rollbackStockQuery, stockRestore, item.ProductID).Error; err != nil {
			return err
		}

		stockAfter := stockBefore - stockRestore
		notes := fmt.Sprintf("Void purchase order ID %d", id)
		if err := r.db.Exec(createStockMutationQuery,
			item.ProductID, "void_purchase", stockRestore, stockBefore, stockAfter,
			"purchase", id, notes, userID,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

// AddItems menambahkan item baru ke PO yang sudah ada pembayaran, tanpa mengubah
// item lama. Total dihitung ulang (total lama + subtotal item baru) dan payment_status
// disesuaikan otomatis (paid_amount tetap).
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
			stockAdd := item.Quantity * conversionQty

			if err := tx.Exec(createPurchaseItemQuery,
				req.ID, item.ProductID, item.Quantity, item.Unit, conversionQty, item.PurchasePrice, subtotal,
			).Error; err != nil {
				return err
			}

			if len(item.ExpiryBatches) > 0 {
				var purchaseItemID int
				if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&purchaseItemID).Error; err != nil {
					return err
				}
				if err := insertExpiryBatches(tx, item.ProductID, purchaseItemID, item.ExpiryBatches); err != nil {
					return err
				}
			}

			var stockBefore float64
			if err := tx.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
				return err
			}

			if err := tx.Exec(addStockQuery, stockAdd, item.ProductID).Error; err != nil {
				return err
			}

			stockAfter := stockBefore + stockAdd
			notes := fmt.Sprintf("Tambah item PO ID %d", req.ID)
			if err := tx.Exec(createStockMutationQuery,
				item.ProductID, "in", stockAdd, stockBefore, stockAfter,
				"purchase", req.ID, notes, req.UserID,
			).Error; err != nil {
				return err
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
