package repo

import (
	"fmt"
	"time"

	cash_drawer_model "pos_api/domain/cash_drawer/model"
	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const liveExpectedBalanceExpr = `(cd.opening_balance + COALESCE((SELECT SUM(total_amount) FROM transactions WHERE user_id = cd.user_id AND payment_method = 'cash' AND status = 'completed' AND transaction_date >= cd.open_time), 0) - COALESCE((SELECT SUM(amount) FROM expenses WHERE user_id = cd.user_id AND created_at >= cd.open_time), 0))`

const (
	openBackdateCashDrawerQuery = `INSERT INTO cash_drawer (user_id, shift_id, open_time, opening_balance, open_notes, status, is_backdate, created_by) VALUES (?, ?, ?, ?, ?, 'open', 1, ?)`

	getCurrentBackdateQuery = `
		SELECT cd.id, cd.user_id, u.full_name as user_name, cd.shift_id, s.name as shift_name,
		       cd.open_time, cd.opening_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses,
		       ` + liveExpectedBalanceExpr + ` as expected_balance, cd.status, cd.open_notes,
		       cd.created_by, COALESCE(cb.full_name, '') as created_by_name
		FROM cash_drawer cd
		LEFT JOIN users u ON cd.user_id = u.id
		LEFT JOIN shifts s ON cd.shift_id = s.id
		LEFT JOIN users cb ON cd.created_by = cb.id
		WHERE cd.created_by = ? AND cd.is_backdate = 1 AND cd.status = 'open'
		LIMIT 1`

	closeBackdateCashDrawerQuery = `UPDATE cash_drawer SET close_time = ?, closing_balance = ?, expected_balance = ?, difference = ?, status = 'closed', notes = ?, updated_at = ? WHERE id = ? AND is_backdate = 1`

	getBackdateCashDrawerByIDQuery = `SELECT cd.id, cd.user_id, cd.shift_id, cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses, CASE WHEN cd.status = 'closed' THEN cd.expected_balance ELSE ` + liveExpectedBalanceExpr + ` END as expected_balance, cd.difference, cd.status, cd.notes, cd.is_backdate, cd.created_by FROM cash_drawer cd WHERE cd.id = ? AND cd.is_backdate = 1 LIMIT 1`

	generateBackdateTransactionCodeQuery = `SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(transaction_code, '-', -1) AS UNSIGNED)), 0) FROM transactions WHERE DATE(transaction_date) = ? AND device_source = ?`

	createBackdateTransactionQuery = `INSERT INTO transactions (transaction_code, user_id, shift_id, transaction_date, subtotal, discount, tax, total_amount, payment_method, payment_amount, change_amount, balance_used, customer_id, is_credit, status, device_source, cash_drawer_id, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	createBackdateTransactionItemQuery = `INSERT INTO transaction_items (transaction_id, product_id, product_name, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updateBackdateSalesQuery = `UPDATE cash_drawer SET total_sales = total_sales + ?, total_cash_sales = total_cash_sales + ?, updated_at = ? WHERE id = ?`

	createReceivableQuery = `INSERT INTO receivables (transaction_id, customer_id, total_amount, remaining_amount, status) VALUES (?, ?, ?, ?, 'unpaid')`
)

type BackdateRepoInterface interface {
	OpenCashDrawer(userID int, shiftID *int, openTime time.Time, openingBalance float64, notes string, createdBy int) (int64, error)
	GetCurrentCashDrawer(createdBy int) (*cash_drawer_model.CashDrawer, error)
	GetCashDrawerByID(id int) (*cash_drawer_model.CashDrawer, error)
	CloseCashDrawer(id int, closingBalance, expectedBalance, difference float64, notes string, closeTime time.Time) error
	CreateTransaction(userID int, shiftID *int, transactionDate time.Time, cashDrawerID int, createdBy int, req *CreateTransactionParams) (*CreateTransactionResult, error)
	GetDB() *gorm.DB
}

type CreateTransactionParams struct {
	Subtotal      float64
	Discount      float64
	Tax           float64
	TotalAmount   float64
	PaymentMethod string
	PaymentAmount float64
	ChangeAmount  float64
	BalanceUsed   float64
	CustomerID    *int
	IsCredit      bool
	DeviceSource  string
	Items         []CreateTransactionItemParams
}

type CreateTransactionItemParams struct {
	ProductID    int
	ProductName  string
	Quantity     float64
	Unit         string
	Price        float64
	Subtotal     float64
	DiscountItem float64
	UnitID       *int
}

type CreateTransactionResult struct {
	ID              int
	TransactionCode string
	TransactionDate time.Time
	TotalAmount     float64
}

type backdateRepo struct {
	db *gorm.DB
}

func NewBackdateRepo(db *gorm.DB) *backdateRepo {
	return &backdateRepo{db: db}
}

func (r *backdateRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *backdateRepo) OpenCashDrawer(userID int, shiftID *int, openTime time.Time, openingBalance float64, notes string, createdBy int) (int64, error) {
	if err := r.db.Exec(openBackdateCashDrawerQuery, userID, shiftID, openTime, openingBalance, notes, createdBy).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *backdateRepo) GetCurrentCashDrawer(createdBy int) (*cash_drawer_model.CashDrawer, error) {
	var data struct {
		ID              int       `gorm:"column:id"`
		UserID          int       `gorm:"column:user_id"`
		UserName        string    `gorm:"column:user_name"`
		ShiftID         *int      `gorm:"column:shift_id"`
		ShiftName       *string   `gorm:"column:shift_name"`
		OpenTime        time.Time `gorm:"column:open_time"`
		OpeningBalance  float64   `gorm:"column:opening_balance"`
		TotalSales      float64   `gorm:"column:total_sales"`
		TotalCashSales  float64   `gorm:"column:total_cash_sales"`
		TotalExpenses   float64   `gorm:"column:total_expenses"`
		ExpectedBalance float64   `gorm:"column:expected_balance"`
		Status          string    `gorm:"column:status"`
		OpenNotes       *string   `gorm:"column:open_notes"`
		CreatedBy       *int      `gorm:"column:created_by"`
		CreatedByName   string    `gorm:"column:created_by_name"`
	}
	if err := r.db.Raw(getCurrentBackdateQuery, createdBy).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	cd := &cash_drawer_model.CashDrawer{
		ID:              data.ID,
		UserID:          data.UserID,
		ShiftID:         data.ShiftID,
		OpenTime:        data.OpenTime,
		OpeningBalance:  data.OpeningBalance,
		TotalSales:      data.TotalSales,
		TotalCashSales:  data.TotalCashSales,
		TotalExpenses:   data.TotalExpenses,
		ExpectedBalance: data.ExpectedBalance,
		Status:          data.Status,
		IsBackdate:      true,
		CreatedBy:       data.CreatedBy,
	}
	return cd, nil
}

func (r *backdateRepo) GetCashDrawerByID(id int) (*cash_drawer_model.CashDrawer, error) {
	var cd cash_drawer_model.CashDrawer
	if err := r.db.Raw(getBackdateCashDrawerByIDQuery, id).Scan(&cd).Error; err != nil {
		return nil, err
	}
	if cd.ID == 0 {
		return nil, nil
	}
	return &cd, nil
}

func (r *backdateRepo) CloseCashDrawer(id int, closingBalance, expectedBalance, difference float64, notes string, closeTime time.Time) error {
	return r.db.Exec(closeBackdateCashDrawerQuery, closeTime, closingBalance, expectedBalance, difference, notes, time_helper.GetTimeNow(), id).Error
}

func (r *backdateRepo) CreateTransaction(userID int, shiftID *int, transactionDate time.Time, cashDrawerID int, createdBy int, req *CreateTransactionParams) (*CreateTransactionResult, error) {
	var result CreateTransactionResult

	prefixMap := map[string]string{"desktop": "DSK", "web": "WEB", "android": "AND"}
	prefix, ok := prefixMap[req.DeviceSource]
	if !ok {
		prefix = "POS"
	}

	targetDateStr := time_helper.ToSQLDate(transactionDate)
	lockName := fmt.Sprintf("txcode:%s:%s", targetDateStr, req.DeviceSource)

	var code string
	var transactionID int
	if err := r.db.Connection(func(conn *gorm.DB) error {
		var locked int
		if err := conn.Raw(`SELECT GET_LOCK(?, 5)`, lockName).Scan(&locked).Error; err != nil {
			return err
		}
		if locked != 1 {
			return fmt.Errorf("gagal mendapatkan lock generate kode transaksi backdate (timeout)")
		}
		defer conn.Exec(`SELECT RELEASE_LOCK(?)`, lockName)

		var count int
		if err := conn.Raw(generateBackdateTransactionCodeQuery, targetDateStr, req.DeviceSource).Scan(&count).Error; err != nil {
			return err
		}
		code = fmt.Sprintf("%s-%s-%03d", prefix, transactionDate.Format("20060102"), count+1)

		if err := conn.Exec(createBackdateTransactionQuery,
			code, userID, shiftID, transactionDate,
			req.Subtotal, req.Discount, req.Tax, req.TotalAmount,
			req.PaymentMethod, req.PaymentAmount, req.ChangeAmount, req.BalanceUsed,
			req.CustomerID, req.IsCredit, "completed", req.DeviceSource,
			cashDrawerID, createdBy,
		).Error; err != nil {
			return err
		}

		return conn.Raw(`SELECT LAST_INSERT_ID()`).Scan(&transactionID).Error
	}); err != nil {
		return nil, err
	}

	// Insert items + apply stock
	for _, item := range req.Items {
		packageID := 0
		if item.UnitID != nil && *item.UnitID > 0 {
			packageID = *item.UnitID
		} else if resolved, ok := product_repo.ResolveDefaultPackageID(r.db, item.ProductID); ok {
			packageID = resolved
		}

		// Get purchase price
		var purchasePrice float64
		r.db.Raw(`SELECT purchase_price FROM products WHERE id = ? LIMIT 1`, item.ProductID).Scan(&purchasePrice)

		// Hitung conversion qty
		conversionQty := item.Quantity
		if packageID > 0 {
			var pkgQty float64
			r.db.Raw(`SELECT qty FROM product_packages WHERE id = ?`, packageID).Scan(&pkgQty)
			if pkgQty > 0 {
				conversionQty = item.Quantity * pkgQty
			}
		}

		if err := r.db.Exec(createBackdateTransactionItemQuery,
			transactionID, item.ProductID, item.ProductName, item.Quantity, item.Unit,
			item.Price, purchasePrice, item.Subtotal, item.DiscountItem, conversionQty, item.UnitID,
		).Error; err != nil {
			return nil, err
		}

		// Apply stock delta (reduce)
		if packageID == 0 {
			continue
		}
		notes := fmt.Sprintf("Penjualan backdate %s", code)
		if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
			ProductID:     item.ProductID,
			PackageID:     packageID,
			Quantity:      item.Quantity,
			Direction:     model_product.StockOut,
			MutationType:  "sale",
			ReferenceType: "transaction",
			ReferenceID:   transactionID,
			Notes:         notes,
			UserID:        &userID,
		}); err != nil {
			return nil, product_repo.WrapStockError(err, item.ProductID)
		}
	}

	// Create receivable if kredit
	if req.IsCredit && req.CustomerID != nil {
		if err := r.db.Exec(createReceivableQuery, transactionID, *req.CustomerID, req.TotalAmount, req.TotalAmount).Error; err != nil {
			return nil, err
		}
	}

	result.ID = transactionID
	result.TransactionCode = code
	result.TransactionDate = transactionDate
	result.TotalAmount = req.TotalAmount
	return &result, nil
}
