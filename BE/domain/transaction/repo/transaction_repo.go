package repo

import (
	"encoding/json"
	stderrors "errors"
	"fmt"

	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	dto_sync "pos_api/domain/sync/dto"
	"pos_api/domain/transaction/dto"
	"pos_api/domain/transaction/model"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"
	"pos_api/pkg/syncmap"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	getPackagesByProductQuery    = `SELECT pp.id, pp.ref_package_id, pp.qty, pp.ref_qty, pp.purchase_price, COALESCE(u.name, '') AS unit_name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.product_id = ?`
	generateTransactionCodeQuery = `SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(transaction_code, '-', -1) AS UNSIGNED)), 0) FROM transactions WHERE DATE(transaction_date) = ? AND device_source = ?`
	createTransactionQuery       = `INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, subtotal, discount, tax, total_amount, payment_method, payment_amount, change_amount, balance_used, customer_id, is_credit, status, device_source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createTransactionItemQuery   = `INSERT INTO transaction_items (transaction_id, product_id, product_name, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	voidTransactionQuery         = `UPDATE transactions SET status = 'void', updated_at = ? WHERE id = ?`
	createReceivableQuery        = `INSERT INTO receivables (transaction_id, customer_id, total_amount, remaining_amount, status) VALUES (?, ?, ?, ?, 'unpaid')`
	updateReceivableVoidQuery    = `UPDATE receivables SET status = 'void', updated_at = ? WHERE transaction_id = ?`
	getProductPurchasePriceQuery = `SELECT purchase_price FROM products WHERE id = ? LIMIT 1`
	getTransactionItemsQuery     = `SELECT id, transaction_id, product_id, product_name, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id FROM transaction_items WHERE transaction_id = ?`
	getTransactionForVoidQuery   = `SELECT user_id, payment_method, total_amount FROM transactions WHERE id = ? LIMIT 1 FOR UPDATE`
	getTransactionByIDQuery      = `
		SELECT t.id, t.transaction_code, t.user_id, COALESCE(u.full_name, '') AS kasir_name,
		       t.shift_id, t.transaction_date,
		       t.subtotal, t.discount, t.tax, t.total_amount, t.payment_method,
		       t.payment_amount, t.change_amount, t.balance_used, t.customer_id, COALESCE(c.name, '') AS customer_name,
		       t.is_credit, t.status, t.device_source
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN customers c ON c.id = t.customer_id
		WHERE t.id = ? LIMIT 1`
	getAllTransactionsBase = `
		SELECT t.id, t.transaction_code, t.user_id, COALESCE(u.full_name, '') AS kasir_name,
		       t.shift_id, t.transaction_date,
		       t.subtotal, t.discount, t.tax, t.total_amount, t.payment_method,
		       t.payment_amount, t.change_amount, t.balance_used, t.customer_id, COALESCE(c.name, '') AS customer_name,
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
			&t.PaymentAmount, &t.ChangeAmount, &t.BalanceUsed, &t.CustomerID, &t.CustomerName,
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
		&t.PaymentAmount, &t.ChangeAmount, &t.BalanceUsed, &t.CustomerID, &t.CustomerName,
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

func isDuplicateTransactionCodeError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return stderrors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *transactionRepo) Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	const maxCodeRetries = 5
	var resp *dto.CreateTransactionResponse
	var err error
	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		resp, err = r.createOnce(req, userID)
		if err == nil || !isDuplicateTransactionCodeError(err) {
			break
		}
	}
	return resp, err
}

func (r *transactionRepo) createOnce(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	var resp dto.CreateTransactionResponse

	now := time_helper.GetTimeNow()

	prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
	prefix, ok := prefixMap[req.DeviceSource]
	if !ok {
		prefix = "POS"
	}

	// Retry-on-duplicate-key SENDIRIAN ternyata tidak cukup di sini: dites nyata
	// dengan 2 request BENAR-BENAR konkuren (2 goroutine terpisah), retry-nya jalan
	// tapi SELECT MAX() di setiap percobaan retry masih membaca angka yang sama
	// persis berkali-kali (rows:0 x5 dengan kode identik) -- goroutine kedua tidak
	// langsung melihat commit dari goroutine pertama meski secara teori seharusnya
	// sudah ter-commit duluan. Makanya generate kode + insert sekarang diserialkan
	// pakai MySQL named lock (GET_LOCK/RELEASE_LOCK, scoped per tanggal+device_source)
	// di SATU koneksi terkunci (r.db.Connection) -- ini menghilangkan race di
	// sumbernya, bukan cuma mencoba lagi setelah tabrakan. Retry-on-duplicate di
	// Create() tetap dipertahankan sebagai lapis kedua (mis. kalau GET_LOCK timeout).
	lockName := fmt.Sprintf("txcode:%s:%s", time_helper.ToSQLDate(now), req.DeviceSource)
	var code string
	var transactionID int
	if err := r.db.Connection(func(conn *gorm.DB) error {
		var locked int
		if err := conn.Raw(`SELECT GET_LOCK(?, 5)`, lockName).Scan(&locked).Error; err != nil {
			return err
		}
		if locked != 1 {
			return fmt.Errorf("gagal mendapatkan lock generate kode transaksi (timeout)")
		}
		defer conn.Exec(`SELECT RELEASE_LOCK(?)`, lockName)

		var count int
		if err := conn.Raw(generateTransactionCodeQuery, time_helper.ToSQLDate(now), req.DeviceSource).Scan(&count).Error; err != nil {
			return err
		}
		code = fmt.Sprintf("%s-%s-%03d", prefix, now.Format("20060102"), count+1)

		if err := conn.Exec(createTransactionQuery,
			code, userID, req.ShiftID, now,
			req.Subtotal, req.Discount, req.Tax, req.TotalAmount,
			req.PaymentMethod, req.PaymentAmount, req.ChangeAmount, req.BalanceUsed,
			req.CustomerID, req.IsCredit, "completed", req.DeviceSource,
		).Error; err != nil {
			return err
		}

		return conn.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error
	}); err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		// Resolusi package_id: pakai yang dikirim FE (item.UnitID -- nama
		// kolom lama, isinya sebenarnya package_id, lihat komentar di
		// dto/model), fallback ke paket anchor kalau kosong. Sama pola dgn
		// celah #15 di purchase_repo.go.
		packageID := 0
		if item.UnitID != nil && *item.UnitID > 0 {
			packageID = *item.UnitID
		} else if resolved, ok := product_repo.ResolveDefaultPackageID(r.db, item.ProductID); ok {
			packageID = resolved
		}
		if packageID == 0 {
			return nil, fmt.Errorf("produk %s tidak punya paket satuan yang bisa dipakai buat jual", item.ProductName)
		}

		// conversion_qty & unitName di sini murni utk kolom informasi/laporan
		// (purchase_price margin, nama satuan tampilan) -- BUKAN lagi acuan
		// pengurangan stok, itu tugas ApplyStockDelta di bawah (dihitung
		// fresh dari package_id, bukan baca balik nilai ini -- celah #21).
		conversionQty := item.ConversionQty
		unitName := item.Unit
		var packagePurchasePrice float64
		var pkgRows []*model_product.ProductPackage
		if err := r.db.Raw(getPackagesByProductQuery, item.ProductID).Scan(&pkgRows).Error; err == nil {
			if factor, factorErr := model_product.ResolvePackageFactor(pkgRows, packageID); factorErr == nil && factor > 0 {
				conversionQty = factor
				for _, p := range pkgRows {
					if p.ID == packageID {
						unitName = p.UnitName
						packagePurchasePrice = p.PurchasePrice
						break
					}
				}
			}
		}
		if conversionQty <= 0 {
			conversionQty = 1
		}

		// HPP pakai harga beli RIIL yang dicatat manual per paket (product_packages.
		// purchase_price), BUKAN sekadar harga anchor dikali faktor konversi -- toko
		// bisa beli satuan kecil (mis. Butir) dengan harga per-unit yang TIDAK
		// proporsional dari harga per-kg (mis. beli eceran dari sumber lain, lebih
		// mahal per satuan). Kalau field itu kosong/0 (data lama sebelum field ini
		// dipakai), fallback ke hitungan proporsional lama supaya tidak pecah.
		purchasePrice := packagePurchasePrice
		if purchasePrice <= 0 {
			var basePurchasePrice float64
			if err := r.db.Raw(getProductPurchasePriceQuery, item.ProductID).Scan(&basePurchasePrice).Error; err != nil {
				return nil, err
			}
			purchasePrice = basePurchasePrice * conversionQty
		}

		if err := r.db.Exec(createTransactionItemQuery,
			transactionID, item.ProductID, item.ProductName,
			item.Quantity, unitName, item.Price, purchasePrice, item.Subtotal,
			item.DiscountItem, conversionQty, packageID,
		).Error; err != nil {
			return nil, err
		}

		notes := fmt.Sprintf("Transaksi %s", code)
		if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
			ProductID:     item.ProductID,
			PackageID:     packageID,
			Quantity:      item.Quantity,
			Direction:     model_product.StockOut,
			MutationType:  "out",
			ReferenceType: "transaction",
			ReferenceID:   transactionID,
			Notes:         notes,
			UserID:        &userID,
		}); err != nil {
			// Bug QA Fase A skenario 15: dulu HANYA ErrInsufficientStock yang
			// diterjemahkan (lewat konvensi prefix string "stok_insufficient:"
			// yang ditangkap transaction_service.go) -- ErrNeedsStockReview &
			// ErrBranchingChain bocor sebagai error mentah -> 500 Internal
			// Server Error, bukan 400 dengan pesan jelas. product_repo.WrapStockError
			// menangani ketiganya sekaligus, jadi sudah *errors.BadRequestError
			// yang dikenali (service tinggal pass-through, lihat perbaikan di sana).
			return nil, product_repo.WrapStockError(err, item.ProductID)
		}
	}

	if req.IsCredit && req.CustomerID != nil {
		receivableAmount := req.TotalAmount - req.BalanceUsed
		if receivableAmount > 0 {
			if err := r.db.Exec(createReceivableQuery,
				transactionID, *req.CustomerID, receivableAmount, receivableAmount,
			).Error; err != nil {
				return nil, err
			}
		}
	}

	resp.ID = transactionID
	resp.TransactionCode = code
	resp.UserID = userID
	resp.ShiftID = req.ShiftID
	resp.TransactionDate = now
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
	now := time_helper.GetTimeNow()

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
	if err := r.db.Exec(voidTransactionQuery, now, id).Error; err != nil {
		return err
	}

	// 2. Ambil semua items
	items, err := r.GetItems(id)
	if err != nil {
		return err
	}

	// 3. Kembalikan stok & catat mutasi void -- lewat ApplyStockDelta,
	// dihitung fresh dari package_id+quantity asli (celah #21), TIDAK baca
	// balik item.ConversionQty yang tersimpan (berpotensi sudah kepotong).
	notes := fmt.Sprintf("Void transaksi ID %d", id)
	for _, item := range items {
		packageID := 0
		if item.UnitID != nil && *item.UnitID > 0 {
			packageID = *item.UnitID
		} else if resolved, ok := product_repo.ResolveDefaultPackageID(r.db, item.ProductID); ok {
			packageID = resolved
		}
		if packageID == 0 {
			return fmt.Errorf("produk %s tidak punya paket satuan yang bisa dipakai buat void", item.ProductName)
		}

		if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
			ProductID:     item.ProductID,
			PackageID:     packageID,
			Quantity:      item.Quantity,
			Direction:     model_product.StockIn,
			MutationType:  "void",
			ReferenceType: "transaction",
			ReferenceID:   id,
			Notes:         notes,
			UserID:        &userID,
		}); err != nil {
			return product_repo.WrapStockError(err, item.ProductID)
		}
	}

	// 4. Jika ada piutang â†’ update status void
	if err := r.db.Exec(updateReceivableVoidQuery, now, id).Error; err != nil {
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
			&item.Quantity, &item.Unit, &item.Price, &item.PurchasePrice, &item.Subtotal,
			&item.DiscountItem, &item.ConversionQty, &item.UnitID,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *transactionRepo) ApplySyncTransaction(payload string, deviceID string, localID string, cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface) (int, error) {
	var tx dto_sync.SyncTransactionPayload
	if err := json.Unmarshal([]byte(payload), &tx); err != nil {
		return 0, fmt.Errorf("payload transaksi tidak valid: %w", err)
	}

	if existingID, found, err := syncmap.Resolve(r.db, deviceID, localID, "transaction"); err == nil && found {
		return existingID, nil
	}

	lockName := fmt.Sprintf("txcode:%s:%s", time_helper.ToSQLDate(time_helper.GetTimeNow()), tx.DeviceSource)
	var serverID int
	var err error
	const maxCodeRetries = 5
	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		serverID = 0
		err = r.db.Connection(func(conn *gorm.DB) error {
			var locked int
			if err := conn.Raw(`SELECT GET_LOCK(?, 5)`, lockName).Scan(&locked).Error; err != nil {
				return err
			}
			if locked != 1 {
				return fmt.Errorf("gagal mendapatkan lock generate kode sync transaksi (timeout)")
			}
			defer conn.Exec(`SELECT RELEASE_LOCK(?)`, lockName)
			return conn.Transaction(func(db *gorm.DB) error {
				return r.applySyncTransactionOnce(db, tx, deviceID, localID, cashDrawerRepo, &serverID)
			})
		})
		if err == nil || !isDuplicateTransactionCodeError(err) {
			break
		}
	}

	return serverID, err
}

func (r *transactionRepo) applySyncTransactionOnce(db *gorm.DB, tx dto_sync.SyncTransactionPayload, deviceID string, localID string, cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface, serverID *int) error {
	{
		now := time_helper.GetTimeNow()
		prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
		prefix, ok := prefixMap[tx.DeviceSource]
		if !ok {
			prefix = "DSK"
		}
		var count int
		if err := db.Raw(generateTransactionCodeQuery, time_helper.ToSQLDate(now), tx.DeviceSource).Scan(&count).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("%s-%s-%03d", prefix, now.Format("20060102"), count+1)

		// 3. Insert header transaksi
		if err := db.Exec(createTransactionQuery,
			code, tx.UserID, tx.ShiftID, now,
			tx.Subtotal, tx.Discount, tx.Tax, tx.TotalAmount,
			tx.PaymentMethod, tx.PaymentAmount, tx.ChangeAmount, float64(0),
			tx.CustomerID, tx.IsCredit, "completed", tx.DeviceSource,
		).Error; err != nil {
			return err
		}

		var transactionID int
		if err := db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error; err != nil {
			return err
		}

		if err := syncmap.Record(db, deviceID, localID, "transaction", transactionID); err != nil {
			return err
		}

		// 4. Kurangi stok (lewat ApplyStockDelta, celah #19) + insert item
		notes := fmt.Sprintf("Sync offline tx %s", localID)
		for _, item := range tx.Items {
			packageID := 0
			if item.UnitID != nil && *item.UnitID > 0 {
				packageID = *item.UnitID
			} else if resolved, ok := product_repo.ResolveDefaultPackageID(db, item.ProductID); ok {
				packageID = resolved
			}
			if packageID == 0 {
				return fmt.Errorf("produk %d tidak punya paket satuan yang bisa dipakai buat sync", item.ProductID)
			}

			// HPP pakai harga beli riil per paket kalau ada (sama seperti alur online
			// di createOnce()), fallback ke harga anchor × conversion_qty kalau kosong.
			var purchasePrice float64
			var pkgRows []*model_product.ProductPackage
			if err := db.Raw(getPackagesByProductQuery, item.ProductID).Scan(&pkgRows).Error; err == nil {
				for _, p := range pkgRows {
					if p.ID == packageID {
						purchasePrice = p.PurchasePrice
						break
					}
				}
			}
			if purchasePrice <= 0 {
				var basePurchasePrice float64
				if err := db.Raw(getProductPurchasePriceQuery, item.ProductID).Scan(&basePurchasePrice).Error; err != nil {
					return err
				}
				purchasePrice = basePurchasePrice * item.ConversionQty
			}

			if err := db.Exec(createTransactionItemQuery,
				transactionID, item.ProductID, item.ProductName,
				item.Quantity, item.Unit, item.Price, purchasePrice, item.Subtotal,
				item.DiscountItem, item.ConversionQty, packageID,
			).Error; err != nil {
				return err
			}

			userID := tx.UserID
			if _, err := product_repo.ApplyStockDelta(db, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     packageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockOut,
				MutationType:  "out",
				ReferenceType: "transaction",
				ReferenceID:   transactionID,
				Notes:         notes,
				UserID:        &userID,
			}); err != nil {
				if stderrors.Is(err, model_product.ErrInsufficientStock) {
					return fmt.Errorf("stok produk %d tidak mencukupi", item.ProductID)
				}
				return err
			}
		}

		// 5. Jika kredit â†’ buat piutang
		if tx.IsCredit && tx.CustomerID != nil {
			if err := db.Exec(createReceivableQuery, transactionID, *tx.CustomerID, tx.TotalAmount, tx.TotalAmount).Error; err != nil {
				return err
			}
		}

		if cashDrawerRepo != nil && tx.PaymentMethod == "cash" {
			cashDrawerTx := cashDrawerRepo.WithTx(db)
			drawer, err := cashDrawerTx.GetOpenCashDrawer(tx.UserID)
			if err != nil {
				return err
			}
			if drawer != nil {
				if err := cashDrawerTx.UpdateSales(drawer.ID, tx.TotalAmount, tx.TotalAmount, time_helper.GetTimeNow()); err != nil {
					return err
				}
			}
		}

		*serverID = transactionID
		return nil
	}
}

// ReturnStockForRejectSync membalikkan stok saat konflik sync transaksi
// offline ditolak. Celah #19: dulu punya bug SIMETRIS dengan
// ApplySyncTransaction (pakai item.Quantity mentah, bukan dikali faktor
// konversi) DAN memakai mutation_type 'REJECT_SYNC' yang bukan bagian dari
// enum stock_mutations.mutation_type (enum('in','out','adjustment','void',
// 'return','void_purchase','expired')) -- diganti 'void' (paling pas secara
// makna: membatalkan efek stok dari sebuah penjualan).
func (r *transactionRepo) ReturnStockForRejectSync(transactionID, resolvedBy int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		items, err := r.GetItems(transactionID)
		if err != nil {
			return err
		}

		notes := fmt.Sprintf("Reject sync konflik transaksi offline ID %d", transactionID)
		for _, item := range items {
			packageID := 0
			if item.UnitID != nil && *item.UnitID > 0 {
				packageID = *item.UnitID
			} else if resolved, ok := product_repo.ResolveDefaultPackageID(tx, item.ProductID); ok {
				packageID = resolved
			}
			if packageID == 0 {
				return fmt.Errorf("produk %d tidak punya paket satuan yang bisa dipakai buat reject sync", item.ProductID)
			}

			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     packageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockIn,
				MutationType:  "void",
				ReferenceType: "transaction",
				ReferenceID:   transactionID,
				Notes:         notes,
				UserID:        &resolvedBy,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

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
	updates["updated_at"] = time_helper.GetTimeNow()
	return r.db.Table("transactions").Where("id = ?", id).Updates(updates).Error
}
