package repo

import (
	"encoding/json"
	"fmt"
	"time"

	dto_sync "pos_api/domain/sync/dto"
	"pos_api/domain/transaction/dto"
	"pos_api/domain/transaction/model"
	request_helper "pos_api/helper/request"

	"gorm.io/gorm"
)

const (
	getPackageByIDQuery           = `SELECT pp.conversion_qty, COALESCE(u.name, '') AS unit_name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.id = ? LIMIT 1`
	generateTransactionCodeQuery  = `SELECT COUNT(*) FROM transactions WHERE DATE(transaction_date) = CURDATE() AND device_source = ?`
	createTransactionQuery        = `INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, subtotal, discount, tax, total_amount, payment_method, payment_amount, change_amount, customer_id, is_credit, status, device_source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createTransactionItemQuery    = `INSERT INTO transaction_items (transaction_id, product_id, product_name, quantity, unit, price, subtotal, discount_item, conversion_qty, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	updateProductStockQuery       = `UPDATE products SET stock = stock - ?, updated_at = NOW() WHERE id = ? AND (stock - reserved_qty) >= ?`
	createStockMutationQuery      = `INSERT INTO stock_mutations (product_id, mutation_type, quantity, stock_before, stock_after, reference_type, reference_id, notes, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	voidTransactionQuery          = `UPDATE transactions SET status = 'void', updated_at = NOW() WHERE id = ?`
	restoreStockQuery             = `UPDATE products SET stock = stock + ?, updated_at = NOW() WHERE id = ?`
	createReceivableQuery         = `INSERT INTO receivables (transaction_id, customer_id, total_amount, remaining_amount, status) VALUES (?, ?, ?, ?, 'unpaid')`
	updateReceivableVoidQuery     = `UPDATE receivables SET status = 'void', updated_at = NOW() WHERE transaction_id = ?`
	getProductStockQuery          = `SELECT stock FROM products WHERE id = ? LIMIT 1`
	getTransactionItemsQuery      = `SELECT id, transaction_id, product_id, product_name, quantity, unit, price, subtotal, discount_item, conversion_qty, unit_id FROM transaction_items WHERE transaction_id = ?`
	getTransactionForVoidQuery    = `SELECT user_id, payment_method, total_amount FROM transactions WHERE id = ? LIMIT 1 FOR UPDATE`
	getTransactionByIDQuery       = `
		SELECT t.id, t.transaction_code, t.user_id, COALESCE(u.full_name, '') AS kasir_name,
		       t.shift_id, t.transaction_date,
		       t.subtotal, t.discount, t.tax, t.total_amount, t.payment_method,
		       t.payment_amount, t.change_amount, t.customer_id, COALESCE(c.name, '') AS customer_name,
		       t.is_credit, t.status, t.device_source
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN customers c ON c.id = t.customer_id
		WHERE t.id = ? LIMIT 1`
	getAllTransactionsBase = `
		SELECT t.id, t.transaction_code, t.user_id, COALESCE(u.full_name, '') AS kasir_name,
		       t.shift_id, t.transaction_date,
		       t.subtotal, t.discount, t.tax, t.total_amount, t.payment_method,
		       t.payment_amount, t.change_amount, t.customer_id, COALESCE(c.name, '') AS customer_name,
		       t.is_credit, t.status, t.device_source
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN customers c ON c.id = t.customer_id
		WHERE 1=1`
	countTransactionsBase = `SELECT COUNT(*) FROM transactions t WHERE 1=1`
)

func (r *transactionRepo) GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error) {
	var args, countArgs []interface{}
	conditions := ""

	if req.Status != "" {
		conditions += " AND t.status = ?"
		args = append(args, req.Status)
		countArgs = append(countArgs, req.Status)
	}
	if req.PaymentMethod != "" {
		conditions += " AND t.payment_method = ?"
		args = append(args, req.PaymentMethod)
		countArgs = append(countArgs, req.PaymentMethod)
	}
	if req.DateFrom != "" {
		conditions += " AND DATE(t.transaction_date) >= ?"
		args = append(args, req.DateFrom)
		countArgs = append(countArgs, req.DateFrom)
	}
	if req.DateTo != "" {
		conditions += " AND DATE(t.transaction_date) <= ?"
		args = append(args, req.DateTo)
		countArgs = append(countArgs, req.DateTo)
	}
	if req.UserID != nil {
		conditions += " AND t.user_id = ?"
		args = append(args, *req.UserID)
		countArgs = append(countArgs, *req.UserID)
	}
	if req.Search != "" {
		conditions += " AND (t.transaction_code LIKE ? OR c.name LIKE ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like)
		countArgs = append(countArgs, like, like)
	}

	var total int64
	if err := r.db.Raw(countTransactionsBase+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 20, 0)

	allowedSortFields := map[string]string{
		"transaction_date": "t.transaction_date",
		"total_amount":     "t.total_amount",
		"customer_name":    "c.name",
		"kasir_name":       "u.full_name",
		"payment_method":   "t.payment_method",
		"status":           "t.status",
	}
	const defaultOrder = " ORDER BY t.transaction_date DESC"
	query := getAllTransactionsBase + conditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, defaultOrder) +
		fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*dto.TransactionResponse
	for rows.Next() {
		var t dto.TransactionResponse
		if err := rows.Scan(
			&t.ID, &t.TransactionCode, &t.UserID, &t.KasirName, &t.ShiftID, &t.TransactionDate,
			&t.Subtotal, &t.Discount, &t.Tax, &t.TotalAmount, &t.PaymentMethod,
			&t.PaymentAmount, &t.ChangeAmount, &t.CustomerID, &t.CustomerName,
			&t.IsCredit, &t.Status, &t.DeviceSource,
		); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, &t)
	}
	if transactions == nil {
		transactions = []*dto.TransactionResponse{}
	}
	return transactions, total, nil
}

func (r *transactionRepo) GetByID(id int) (*dto.TransactionResponse, error) {
	rows, err := r.db.Raw(getTransactionByIDQuery, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	var t dto.TransactionResponse
	if err := rows.Scan(
		&t.ID, &t.TransactionCode, &t.UserID, &t.KasirName, &t.ShiftID, &t.TransactionDate,
		&t.Subtotal, &t.Discount, &t.Tax, &t.TotalAmount, &t.PaymentMethod,
		&t.PaymentAmount, &t.ChangeAmount, &t.CustomerID, &t.CustomerName,
		&t.IsCredit, &t.Status, &t.DeviceSource,
	); err != nil {
		return nil, err
	}
	rows.Close()

	items, err := r.GetItems(id)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		t.Items = append(t.Items, dto.TransactionItemResponse{
			ID:            item.ID,
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			Price:         item.Price,
			Subtotal:      item.Subtotal,
			DiscountItem:  item.DiscountItem,
			ConversionQty: item.ConversionQty,
			UnitID:        item.UnitID,
		})
	}
	return &t, nil
}

func (r *transactionRepo) Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	var resp dto.CreateTransactionResponse

	// 1. Generate transaction_code
	prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
	prefix, ok := prefixMap[req.DeviceSource]
	if !ok {
		prefix = "POS"
	}
	var count int
	if err := r.db.Raw(generateTransactionCodeQuery, req.DeviceSource).Scan(&count).Error; err != nil {
		return nil, err
	}
	code := fmt.Sprintf("%s-%s-%03d", prefix, time.Now().Format("20060102"), count+1)

	// 2. Simpan header transaksi
	if err := r.db.Exec(createTransactionQuery,
		code, userID, req.ShiftID, time.Now(),
		req.Subtotal, req.Discount, req.Tax, req.TotalAmount,
		req.PaymentMethod, req.PaymentAmount, req.ChangeAmount,
		req.CustomerID, req.IsCredit, "completed", req.DeviceSource,
	).Error; err != nil {
		return nil, err
	}

	var transactionID int
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error; err != nil {
		return nil, err
	}

	// 3. Loop items
	for _, item := range req.Items {
		// Resolve conversion_qty dan unit_name dari product_packages
		conversionQty := item.ConversionQty
		unitName := item.Unit
		if item.UnitID != nil && *item.UnitID > 0 {
			var pkgData struct {
				ConversionQty float64
				UnitName      string
			}
			if err := r.db.Raw(getPackageByIDQuery, *item.UnitID).Scan(&pkgData).Error; err == nil && pkgData.ConversionQty > 0 {
				conversionQty = pkgData.ConversionQty
				unitName = pkgData.UnitName
			}
		}
		if conversionQty <= 0 {
			conversionQty = 1
		}

		// Stok yang dikurangi = qty Ã— conversion_qty (konversi ke satuan dasar)
		stockDeduct := item.Quantity * conversionQty

		// Ambil stok sebelumnya
		var stockBefore float64
		if err := r.db.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
			return nil, err
		}

		// Kurangi stok (atomic dengan cek stok >= qty dalam satuan dasar)
		result := r.db.Exec(updateProductStockQuery, stockDeduct, item.ProductID, stockDeduct)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, fmt.Errorf("stok_insufficient:%s", item.ProductName)
		}

		// Simpan item dengan unit_name dari master dan conversion_qty yang benar
		if err := r.db.Exec(createTransactionItemQuery,
			transactionID, item.ProductID, item.ProductName,
			item.Quantity, unitName, item.Price, item.Subtotal,
			item.DiscountItem, conversionQty, item.UnitID,
		).Error; err != nil {
			return nil, err
		}

		// Catat mutasi stok (dalam satuan dasar)
		stockAfter := stockBefore - stockDeduct
		notes := fmt.Sprintf("Transaksi %s", code)
		if err := r.db.Exec(createStockMutationQuery,
			item.ProductID, "out", stockDeduct, stockBefore, stockAfter,
			"transaction", transactionID, notes, userID,
		).Error; err != nil {
			return nil, err
		}
	}

	// 4. Jika kredit â†’ buat piutang
	if req.IsCredit && req.CustomerID != nil {
		if err := r.db.Exec(createReceivableQuery,
			transactionID, *req.CustomerID, req.TotalAmount, req.TotalAmount,
		).Error; err != nil {
			return nil, err
		}
	}

	resp.ID = transactionID
	resp.TransactionCode = code
	resp.UserID = userID
	resp.ShiftID = req.ShiftID
	resp.TransactionDate = time.Now()
	resp.Subtotal = req.Subtotal
	resp.Discount = req.Discount
	resp.Tax = req.Tax
	resp.TotalAmount = req.TotalAmount
	resp.PaymentMethod = req.PaymentMethod
	resp.PaymentAmount = req.PaymentAmount
	resp.ChangeAmount = req.ChangeAmount
	resp.CustomerID = req.CustomerID
	resp.IsCredit = req.IsCredit
	resp.Status = "completed"
	resp.DeviceSource = req.DeviceSource

	for _, item := range req.Items {
		resp.Items = append(resp.Items, dto.TransactionItemResponse{
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			Price:         item.Price,
			Subtotal:      item.Subtotal,
			DiscountItem:  item.DiscountItem,
			ConversionQty: item.ConversionQty,
			UnitID:        item.UnitID,
		})
	}

	return &resp, nil
}

func (r *transactionRepo) Void(id, userID int) error {
	// 0. Kunci baris transaksi (FOR UPDATE) untuk mencegah race condition saat void bersamaan.
	var voidData struct {
		UserID        int
		PaymentMethod string
		TotalAmount   float64
	}
	if err := r.db.Raw(getTransactionForVoidQuery, id).Scan(&voidData).Error; err != nil {
		return err
	}

	// 1. Update status void
	if err := r.db.Exec(voidTransactionQuery, id).Error; err != nil {
		return err
	}

	// 2. Ambil semua items
	items, err := r.GetItems(id)
	if err != nil {
		return err
	}

	// 3. Kembalikan stok & catat mutasi void
	for _, item := range items {
		var stockBefore float64
		if err := r.db.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
			return err
		}

		convQty := item.ConversionQty
		if convQty <= 0 {
			convQty = 1
		}
		stockRestore := item.Quantity * convQty

		if err := r.db.Exec(restoreStockQuery, stockRestore, item.ProductID).Error; err != nil {
			return err
		}

		stockAfter := stockBefore + stockRestore
		notes := fmt.Sprintf("Void transaksi ID %d", id)
		if err := r.db.Exec(createStockMutationQuery,
			item.ProductID, "void", stockRestore, stockBefore, stockAfter,
			"transaction", id, notes, userID,
		).Error; err != nil {
			return err
		}
	}

	// 4. Jika ada piutang â†’ update status void
	if err := r.db.Exec(updateReceivableVoidQuery, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *transactionRepo) GetItems(transactionID int) ([]model.TransactionItem, error) {
	rows, err := r.db.Raw(getTransactionItemsQuery, transactionID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.TransactionItem
	for rows.Next() {
		var item model.TransactionItem
		if err := rows.Scan(
			&item.ID, &item.TransactionID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.Unit, &item.Price, &item.Subtotal,
			&item.DiscountItem, &item.ConversionQty, &item.UnitID,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// ApplySyncTransaction menerapkan transaksi offline secara atomik dengan SELECT FOR UPDATE.
// Jika stok produk mana pun tidak mencukupi, seluruh transaksi di-rollback dan error dikembalikan.
func (r *transactionRepo) ApplySyncTransaction(payload string, localID string) (int, error) {
	var tx dto_sync.SyncTransactionPayload
	if err := json.Unmarshal([]byte(payload), &tx); err != nil {
		return 0, fmt.Errorf("payload transaksi tidak valid: %w", err)
	}

	var serverID int

	err := r.db.Transaction(func(db *gorm.DB) error {
		// 1. Cek stok semua item dengan SELECT FOR UPDATE (lock baris product)
		for _, item := range tx.Items {
			var currentStock float64
			if err := db.Raw(`SELECT (stock - reserved_qty) FROM products WHERE id = ? FOR UPDATE`, item.ProductID).Scan(&currentStock).Error; err != nil {
				return fmt.Errorf("stok produk %d tidak ditemukan", item.ProductID)
			}
			if currentStock < item.Quantity {
				return fmt.Errorf("stok produk %d tidak mencukupi (%.2f tersedia, butuh %.2f)",
					item.ProductID, currentStock, item.Quantity)
			}
		}

		// 2. Generate kode transaksi
		prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
		prefix, ok := prefixMap[tx.DeviceSource]
		if !ok {
			prefix = "DSK"
		}
		var count int
		if err := db.Raw(generateTransactionCodeQuery, tx.DeviceSource).Scan(&count).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("%s-%s-%03d", prefix, time.Now().Format("20060102"), count+1)

		// 3. Insert header transaksi
		if err := db.Exec(createTransactionQuery,
			code, tx.UserID, tx.ShiftID, time.Now(),
			tx.Subtotal, tx.Discount, tx.Tax, tx.TotalAmount,
			tx.PaymentMethod, tx.PaymentAmount, tx.ChangeAmount,
			tx.CustomerID, tx.IsCredit, "completed", tx.DeviceSource,
		).Error; err != nil {
			return err
		}

		var transactionID int
		if err := db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error; err != nil {
			return err
		}

		// 4. Kurangi stok + insert item + catat mutasi SALE
		for _, item := range tx.Items {
			var stockBefore float64
			if err := db.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
				return err
			}

			if err := db.Exec(updateProductStockQuery, item.Quantity, item.ProductID, item.Quantity).Error; err != nil {
				return err
			}

			if err := db.Exec(createTransactionItemQuery,
				transactionID, item.ProductID, item.ProductName,
				item.Quantity, item.Unit, item.Price, item.Subtotal,
				item.DiscountItem, item.ConversionQty, item.UnitID,
			).Error; err != nil {
				return err
			}

			stockAfter := stockBefore - item.Quantity
			notes := fmt.Sprintf("Sync offline tx %s", localID)
			if err := db.Exec(createStockMutationQuery,
				item.ProductID, "out", item.Quantity, stockBefore, stockAfter,
				"transaction", transactionID, notes, tx.UserID,
			).Error; err != nil {
				return err
			}
		}

		// 5. Jika kredit â†’ buat piutang
		if tx.IsCredit && tx.CustomerID != nil {
			if err := db.Exec(createReceivableQuery, transactionID, *tx.CustomerID, tx.TotalAmount, tx.TotalAmount).Error; err != nil {
				return err
			}
		}

		serverID = transactionID
		return nil
	})

	return serverID, err
}

// ReturnStockForRejectSync mengembalikan stok setiap item transaksi yang ditolak
// dan mencatat mutasi REJECT_SYNC sebagai audit trail.
func (r *transactionRepo) ReturnStockForRejectSync(transactionID, resolvedBy int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		items, err := r.GetItems(transactionID)
		if err != nil {
			return err
		}

		for _, item := range items {
			var stockBefore float64
			if err := tx.Raw(getProductStockQuery, item.ProductID).Scan(&stockBefore).Error; err != nil {
				return err
			}

			if err := tx.Exec(restoreStockQuery, item.Quantity, item.ProductID).Error; err != nil {
				return err
			}

			stockAfter := stockBefore + item.Quantity
			notes := fmt.Sprintf("Reject sync konflik transaksi offline ID %d", transactionID)
			if err := tx.Exec(createStockMutationQuery,
				item.ProductID, "REJECT_SYNC", item.Quantity, stockBefore, stockAfter,
				"transaction", transactionID, notes, resolvedBy,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateFromSync menerapkan data desktop ke tabel transactions saat konflik di-approve.
// Hanya field yang aman di-overwrite; id/transaction_code/created_at tidak disentuh.
func (r *transactionRepo) UpdateFromSync(id int, data map[string]interface{}) error {
	allowed := []string{
		"subtotal", "discount", "tax", "total_amount",
		"payment_method", "payment_amount", "change_amount",
		"customer_id", "is_credit", "status",
	}
	updates := make(map[string]interface{}, len(allowed))
	for _, key := range allowed {
		if val, ok := data[key]; ok {
			updates[key] = val
		}
	}
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = "NOW()"
	return r.db.Table("transactions").Where("id = ?", id).Updates(updates).Error
}
