package repo

import (
	stderrors "errors"
	"fmt"
	"strings"

	model_product "coffee_api/domain/product/model"
	product_repo "coffee_api/domain/product/repo"
	queue_repo "coffee_api/domain/queue_number/repo"
	"coffee_api/domain/transaction/dto"
	"coffee_api/domain/transaction/model"
	app_errors "coffee_api/errors"
	request_helper "coffee_api/helper/request"
	time_helper "coffee_api/helper/time"

	"github.com/go-sql-driver/mysql"
)

// Mutasi stok produk untuk barang jadi yang dibuang saat void/batal (keputusan 100, migrasi 007).
const MutationWaste = "waste"

const (
	getPackagesByProductQuery    = `SELECT pp.id, pp.ref_package_id, pp.qty, pp.ref_qty, pp.purchase_price, COALESCE(u.name, '') AS unit_name FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.product_id = ?`
	generateTransactionCodeQuery = `SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(transaction_code, '-', -1) AS UNSIGNED)), 0) FROM transactions WHERE DATE(transaction_date) = ? AND device_source = ?`
	createTransactionQuery       = `INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, subtotal, discount, tax, total_amount, payment_method, payment_amount, change_amount, balance_used, customer_id, is_credit, status, device_source, order_type, table_id, table_name, queue_date, queue_number, order_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createTransactionItemQuery   = `INSERT INTO transaction_items (transaction_id, menu_item_id, variant_id, product_id, product_name, variant_name, notes, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	createItemModifierQuery      = `INSERT INTO transaction_item_modifiers (transaction_item_id, option_id, group_name, option_name, price) VALUES (?, ?, ?, ?, ?)`
	voidTransactionQuery         = `UPDATE transactions SET status = 'void', updated_at = ? WHERE id = ?`
	createReceivableQuery        = `INSERT INTO receivables (transaction_id, customer_id, total_amount, remaining_amount, status) VALUES (?, ?, ?, ?, 'unpaid')`
	updateReceivableVoidQuery    = `UPDATE receivables SET status = 'void', updated_at = ? WHERE transaction_id = ?`
	getTransactionItemsQuery     = `SELECT id, transaction_id, menu_item_id, variant_id, COALESCE(product_id, 0) AS product_id, product_name, COALESCE(variant_name, '') AS variant_name, COALESCE(notes, '') AS notes, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id FROM transaction_items WHERE transaction_id = ? ORDER BY id`
	getItemModifiersQuery        = `SELECT m.transaction_item_id, m.option_id, m.group_name, m.option_name, m.price FROM transaction_item_modifiers m JOIN transaction_items i ON i.id = m.transaction_item_id WHERE i.transaction_id = ? ORDER BY m.id`
	getTransactionForVoidQuery   = `SELECT status FROM transactions WHERE id = ? LIMIT 1 FOR UPDATE`

	transactionColumns = `
		SELECT t.id, t.transaction_code, t.user_id, COALESCE(u.full_name, '') AS kasir_name,
		       t.shift_id, t.transaction_date,
		       t.subtotal, t.discount, t.tax, t.total_amount, t.payment_method,
		       t.payment_amount, t.change_amount, t.balance_used, t.customer_id, COALESCE(c.name, '') AS customer_name,
		       t.is_credit, t.status, t.device_source, t.cash_drawer_id,
		       COALESCE(t.order_type, '') AS order_type, t.table_id, COALESCE(t.table_name, '') AS table_name,
		       COALESCE(t.queue_number, 0) AS queue_number, COALESCE(t.order_name, '') AS order_name
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN customers c ON c.id = t.customer_id`
	getTransactionByIDQuery = transactionColumns + ` WHERE t.id = ? LIMIT 1`
	getAllTransactionsBase  = transactionColumns + ` WHERE 1=1`
	countTransactionsBase   = `SELECT COUNT(*) FROM transactions t LEFT JOIN customers c ON c.id = t.customer_id WHERE 1=1`

	// Data menu untuk pesanan (tahap 4.7a).
	getSaleVariantQuery = `SELECT v.id AS variant_id, mi.id AS menu_item_id, mi.name AS menu_name, mi.stock_mode, mi.product_id,
		mi.is_active AS menu_active, mi.is_sold_out AS sold_out, v.is_active AS variant_active, v.price,
		COALESCE(s.name, '') AS size_name, v.temperature
		FROM menu_item_variants v
		JOIN menu_items mi ON mi.id = v.menu_item_id
		LEFT JOIN menu_sizes s ON s.id = v.size_id
		WHERE v.id = ? LIMIT 1`
	getMenuModifierOptionsQuery = `SELECT o.id AS option_id, o.name AS option_name, o.price, o.is_active AS option_active,
		g.id AS group_id, g.name AS group_name, g.selection_type, g.is_required, g.is_active AS group_active
		FROM menu_item_modifier_groups mg
		JOIN modifier_groups g ON g.id = mg.group_id
		JOIN modifier_options o ON o.group_id = g.id
		WHERE mg.menu_item_id = ?
		ORDER BY mg.sort_order, g.name, o.sort_order, o.id`
	getTableForOrderQuery = `SELECT name, is_active FROM dining_tables WHERE id = ? LIMIT 1`
)

func scanTransaction(scan func(dest ...any) error, t *dto.TransactionResponse) error {
	if err := scan(
		&t.ID, &t.TransactionCode, &t.UserID, &t.KasirName, &t.ShiftID, &t.TransactionDate,
		&t.Subtotal, &t.Discount, &t.Tax, &t.TotalAmount, &t.PaymentMethod,
		&t.PaymentAmount, &t.ChangeAmount, &t.BalanceUsed, &t.CustomerID, &t.CustomerName,
		&t.IsCredit, &t.Status, &t.DeviceSource, &t.CashDrawerID,
		&t.OrderType, &t.TableID, &t.TableName, &t.QueueNumber, &t.OrderName,
	); err != nil {
		return err
	}
	if t.QueueNumber > 0 {
		t.QueueLabel = queue_repo.FormatQueueNumber(t.QueueNumber)
	}
	return nil
}

func (r *transactionRepo) GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error) {
	var args []interface{}
	conditions := ""

	if req.Status != "" {
		conditions += " AND t.status = ?"
		args = append(args, req.Status)
	}
	if req.PaymentMethod != "" {
		conditions += " AND t.payment_method = ?"
		args = append(args, req.PaymentMethod)
	}
	if req.OrderType != "" {
		conditions += " AND t.order_type = ?"
		args = append(args, req.OrderType)
	}
	if req.DateFrom != "" {
		conditions += " AND DATE(t.transaction_date) >= ?"
		args = append(args, req.DateFrom)
	}
	if req.DateTo != "" {
		conditions += " AND DATE(t.transaction_date) <= ?"
		args = append(args, req.DateTo)
	}
	if req.UserID != nil {
		conditions += " AND t.user_id = ?"
		args = append(args, *req.UserID)
	}
	if req.Search != "" {
		// Cari kode transaksi, pelanggan, nama pemesan, meja, atau nomor antrean (mis. "007").
		conditions += " AND (t.transaction_code LIKE ? OR c.name LIKE ? OR t.order_name LIKE ? OR t.table_name LIKE ? OR LPAD(t.queue_number, 3, '0') = ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like, like, like, strings.TrimSpace(req.Search))
	}

	var total int64
	if err := r.db.Raw(countTransactionsBase+conditions, args...).Scan(&total).Error; err != nil {
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
		"queue_number":     "t.queue_number",
	}
	const defaultOrder = " ORDER BY t.transaction_date DESC, t.id DESC"
	query := getAllTransactionsBase + conditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, defaultOrder) +
		fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := []*dto.TransactionResponse{}
	for rows.Next() {
		var t dto.TransactionResponse
		if err := scanTransaction(rows.Scan, &t); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, &t)
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
	if err := scanTransaction(rows.Scan, &t); err != nil {
		return nil, err
	}
	rows.Close()

	items, err := r.GetItems(id)
	if err != nil {
		return nil, err
	}
	mods, err := r.getItemModifiers(id)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		resp := dto.TransactionItemResponse{
			ID:            item.ID,
			MenuItemID:    item.MenuItemID,
			VariantID:     item.VariantID,
			ProductName:   item.ProductName,
			VariantName:   item.VariantName,
			Notes:         item.Notes,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			Price:         item.Price,
			Subtotal:      item.Subtotal,
			DiscountItem:  item.DiscountItem,
			ConversionQty: item.ConversionQty,
			UnitID:        item.UnitID,
			Modifiers:     []dto.ItemModifierResponse{},
		}
		if item.ProductID > 0 {
			pid := item.ProductID
			resp.ProductID = &pid
		}
		for _, m := range mods {
			if m.TransactionItemID == item.ID {
				resp.Modifiers = append(resp.Modifiers, dto.ItemModifierResponse{
					OptionID: m.OptionID, GroupName: m.GroupName, OptionName: m.OptionName, Price: m.Price,
				})
			}
		}
		t.Items = append(t.Items, resp)
	}
	return &t, nil
}

func (r *transactionRepo) getItemModifiers(transactionID int) ([]model.TransactionItemModifier, error) {
	var mods []model.TransactionItemModifier
	if err := r.db.Raw(getItemModifiersQuery, transactionID).Scan(&mods).Error; err != nil {
		return nil, err
	}
	return mods, nil
}

func (r *transactionRepo) GetSaleVariant(variantID int) (*model.SaleVariant, error) {
	var v model.SaleVariant
	if err := r.db.Raw(getSaleVariantQuery, variantID).Scan(&v).Error; err != nil {
		return nil, err
	}
	if v.VariantID == 0 {
		return nil, nil
	}
	return &v, nil
}

func (r *transactionRepo) GetMenuModifierOptions(menuItemID int) ([]model.SaleModifierOption, error) {
	var opts []model.SaleModifierOption
	if err := r.db.Raw(getMenuModifierOptionsQuery, menuItemID).Scan(&opts).Error; err != nil {
		return nil, err
	}
	return opts, nil
}

// GetTableForOrder: nama dan status meja; found = false bila meja tidak ada.
func (r *transactionRepo) GetTableForOrder(id int) (name string, active bool, found bool, err error) {
	var row struct {
		Name     string
		IsActive bool
	}
	res := r.db.Raw(getTableForOrderQuery, id).Scan(&row)
	if res.Error != nil {
		return "", false, false, res.Error
	}
	if res.RowsAffected == 0 {
		return "", false, false, nil
	}
	return row.Name, row.IsActive, true, nil
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

// createOnce menyimpan pesanan dalam transaksi pemanggil: kode transaksi, nomor antrean, item,
// tambahan, potong stok barang jadi, dan utang. Semua gagal bersama bila satu langkah gagal.
func (r *transactionRepo) createOnce(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	var resp dto.CreateTransactionResponse

	now := time_helper.GetTimeNow()

	prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
	prefix, ok := prefixMap[req.DeviceSource]
	if !ok {
		prefix = "POS"
	}

	// Kode transaksi urut per hari per perangkat, diserialkan dengan named lock di koneksi transaksi
	// yang sama (pola Toko, lihat docs/BUG_TRANSAKSI_CANGKANG_HEADER_TANPA_ITEM.md di Toko).
	lockName := fmt.Sprintf("txcode:%s:%s", time_helper.ToSQLDate(now), req.DeviceSource)
	var locked int
	if err := r.db.Raw(`SELECT GET_LOCK(?, 5)`, lockName).Scan(&locked).Error; err != nil {
		return nil, err
	}
	if locked != 1 {
		return nil, fmt.Errorf("gagal mendapatkan lock generate kode transaksi (timeout)")
	}
	defer r.db.Exec(`SELECT RELEASE_LOCK(?)`, lockName)

	var count int
	if err := r.db.Raw(generateTransactionCodeQuery, time_helper.ToSQLDate(now), req.DeviceSource).Scan(&count).Error; err != nil {
		return nil, err
	}
	code := fmt.Sprintf("%s-%s-%03d", prefix, now.Format("20060102"), count+1)

	// Nomor antrean harian diambil saat pesanan dibuat (keputusan 95, 96); ikut batal bila gagal disimpan.
	queueNumber, err := queue_repo.NextQueueNumber(r.db, now)
	if err != nil {
		return nil, err
	}

	var tableName *string
	if req.TableID != nil {
		tableName = &req.TableName
	}
	var orderName *string
	if req.OrderName != "" {
		orderName = &req.OrderName
	}

	if err := r.db.Exec(createTransactionQuery,
		code, userID, req.ShiftID, now,
		req.Subtotal, req.Discount, req.Tax, req.TotalAmount,
		req.PaymentMethod, req.PaymentAmount, req.ChangeAmount,
		req.CustomerID, req.IsCredit, "completed", req.DeviceSource,
		req.OrderType, req.TableID, tableName, time_helper.ToSQLDate(now), queueNumber, orderName,
	).Error; err != nil {
		return nil, err
	}

	var transactionID int
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error; err != nil {
		return nil, err
	}

	for i := range req.Items {
		item := &req.Items[i]
		if err := r.insertItem(transactionID, code, item, userID); err != nil {
			return nil, err
		}
	}

	if req.IsCredit && req.CustomerID != nil && req.TotalAmount > 0 {
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
	resp.OrderType = req.OrderType
	resp.TableID = req.TableID
	resp.TableName = req.TableName
	resp.QueueNumber = queueNumber
	resp.QueueLabel = queue_repo.FormatQueueNumber(queueNumber)
	resp.OrderName = req.OrderName
	return &resp, nil
}

// insertItem menyimpan satu baris pesanan + tambahannya, lalu memotong stok bila menu barang jadi
// (keputusan 40: saat item masuk pesanan, keputusan 41: diblokir bila stok kurang).
func (r *transactionRepo) insertItem(transactionID int, code string, item *dto.CreateTransactionItemRequest, userID int) error {
	var productID, packageID *int
	unitName := ""
	purchasePrice := 0.0
	if item.ProductID != nil {
		resolved, ok := product_repo.ResolveDefaultPackageID(r.db, *item.ProductID)
		if !ok {
			return fmt.Errorf("produk untuk menu %s tidak punya satuan stok", item.MenuName)
		}
		productID, packageID = item.ProductID, &resolved
		var pkgRows []*model_product.ProductPackage
		if err := r.db.Raw(getPackagesByProductQuery, *item.ProductID).Scan(&pkgRows).Error; err != nil {
			return err
		}
		for _, p := range pkgRows {
			if p.ID == resolved {
				unitName = p.UnitName
				purchasePrice = p.PurchasePrice
			}
		}
	}

	var variantName, notes *string
	if item.VariantName != "" {
		variantName = &item.VariantName
	}
	if item.Notes != "" {
		notes = &item.Notes
	}
	if err := r.db.Exec(createTransactionItemQuery,
		transactionID, item.MenuItemID, item.VariantID, productID, item.MenuName, variantName, notes,
		item.Quantity, unitName, item.Price, purchasePrice, item.Subtotal, item.DiscountItem, 1, packageID,
	).Error; err != nil {
		return err
	}
	var itemID int
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&itemID).Error; err != nil {
		return err
	}
	for _, m := range item.Modifiers {
		if err := r.db.Exec(createItemModifierQuery, itemID, m.OptionID, m.GroupName, m.OptionName, m.Price).Error; err != nil {
			return err
		}
	}

	if productID == nil {
		return nil
	}
	if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
		ProductID:     *productID,
		PackageID:     *packageID,
		Quantity:      item.Quantity,
		Direction:     model_product.StockOut,
		MutationType:  "out",
		ReferenceType: "transaction",
		ReferenceID:   transactionID,
		Notes:         fmt.Sprintf("Transaksi %s", code),
		UserID:        &userID,
	}); err != nil {
		// Keputusan 41: penjualan diblokir bila stok barang jadi kurang; pesan memakai nama menu.
		if stderrors.Is(err, model_product.ErrInsufficientStock) {
			return &app_errors.BadRequestError{Message: fmt.Sprintf("Stok %s tidak cukup untuk %g item", item.MenuName, item.Quantity)}
		}
		return product_repo.WrapStockError(err, *productID)
	}
	return nil
}

// Void menandai transaksi void di dalam transaksi pemanggil. Barang jadi: semua kembali ke stok,
// lalu item yang dipilih "dibuang" langsung dikurangi lagi sebagai mutasi "waste" beserta alasannya
// (keputusan 100). Menu racikan tidak menyentuh stok. Mengembalikan false bila sudah void.
func (r *transactionRepo) Void(id, userID int, wasteItemIDs map[int]bool, reason string) error {
	now := time_helper.GetTimeNow()

	var status string
	if err := r.db.Raw(getTransactionForVoidQuery, id).Scan(&status).Error; err != nil {
		return err
	}
	if status == "void" {
		return ErrAlreadyVoid
	}
	if err := r.db.Exec(voidTransactionQuery, now, id).Error; err != nil {
		return err
	}

	items, err := r.GetItems(id)
	if err != nil {
		return err
	}

	notes := fmt.Sprintf("Void transaksi ID %d", id)
	for _, item := range items {
		if item.ProductID == 0 {
			continue
		}
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
		if !wasteItemIDs[item.ID] {
			continue
		}
		if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
			ProductID:     item.ProductID,
			PackageID:     packageID,
			Quantity:      item.Quantity,
			Direction:     model_product.StockOut,
			MutationType:  MutationWaste,
			ReferenceType: "transaction",
			ReferenceID:   id,
			Notes:         fmt.Sprintf("Dibuang saat void transaksi ID %d: %s", id, reason),
			UserID:        &userID,
		}); err != nil {
			return product_repo.WrapStockError(err, item.ProductID)
		}
	}

	return r.db.Exec(updateReceivableVoidQuery, now, id).Error
}

// ErrAlreadyVoid: transaksi sudah void saat baris dikunci (void bersamaan).
var ErrAlreadyVoid = stderrors.New("transaksi sudah di-void")

func (r *transactionRepo) GetItems(transactionID int) ([]model.TransactionItem, error) {
	var items []model.TransactionItem
	if err := r.db.Raw(getTransactionItemsQuery, transactionID).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
