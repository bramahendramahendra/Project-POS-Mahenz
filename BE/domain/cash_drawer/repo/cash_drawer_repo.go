package repo

import (
	dto "pos_api/domain/cash_drawer/dto"
	model "pos_api/domain/cash_drawer/model"

	"gorm.io/gorm"
)

const (
	getCurrentCashDrawerQuery  = `SELECT cd.id, cd.user_id, u.full_name as user_name, cd.shift_id, s.name as shift_name, s.start_time as shift_start, s.end_time as shift_end, cd.open_time, cd.opening_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses, cd.expected_balance, cd.status, cd.open_notes FROM cash_drawer cd LEFT JOIN users u ON cd.user_id = u.id LEFT JOIN shifts s ON cd.shift_id = s.id WHERE cd.user_id = ? AND cd.status = 'open' LIMIT 1`
	getOpenCashDrawerQuery     = `SELECT id, user_id, shift_id, open_time, opening_balance, total_sales, total_cash_sales, total_expenses, expected_balance, status FROM cash_drawer WHERE user_id = ? AND status = 'open' LIMIT 1`
	getCashDrawerByIDQuery     = `SELECT id, user_id, shift_id, open_time, close_time, opening_balance, closing_balance, total_sales, total_cash_sales, total_expenses, expected_balance, difference, status, notes FROM cash_drawer WHERE id = ? LIMIT 1`
	openCashDrawerQuery        = `INSERT INTO cash_drawer (user_id, shift_id, open_time, opening_balance, open_notes, status) VALUES (?, ?, NOW(), ?, ?, 'open')`
	getLastCashDrawerInsertID  = `SELECT LAST_INSERT_ID()`
	closeCashDrawerQuery       = `UPDATE cash_drawer SET close_time = NOW(), closing_balance = ?, expected_balance = ?, difference = ?, status = 'closed', notes = ?, updated_at = NOW() WHERE id = ?`
	updateSalesQuery           = `UPDATE cash_drawer SET total_sales = total_sales + ?, total_cash_sales = total_cash_sales + ?, expected_balance = opening_balance + total_cash_sales - total_expenses, updated_at = NOW() WHERE id = ?`
	updateExpensesQuery        = `UPDATE cash_drawer SET total_expenses = total_expenses + ?, expected_balance = opening_balance + total_cash_sales - total_expenses, updated_at = NOW() WHERE id = ?`
	getCashDrawerHistoryBase   = `SELECT cd.id, u.full_name as user_name, s.name as shift_name, cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance, cd.expected_balance, CASE WHEN cd.status = 'closed' THEN cd.difference ELSE NULL END as difference, cd.total_sales, cd.total_cash_sales, cd.total_expenses, cd.status FROM cash_drawer cd LEFT JOIN users u ON cd.user_id = u.id LEFT JOIN shifts s ON cd.shift_id = s.id WHERE 1=1`
	countCashDrawerHistoryBase = `SELECT COUNT(*) FROM cash_drawer cd WHERE 1=1`

	getCashDrawerDetailQuery = `
		SELECT cd.id, cd.user_id, u.full_name as cashier_name, s.name as shift_name,
		       s.start_time as shift_start, s.end_time as shift_end,
		       cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance,
		       cd.expected_balance, cd.total_cash_sales, cd.total_expenses,
		       CASE WHEN cd.status = 'closed' THEN cd.difference ELSE NULL END as difference,
		       cd.status, cd.notes, cd.open_notes
		FROM cash_drawer cd
		LEFT JOIN users u ON cd.user_id = u.id
		LEFT JOIN shifts s ON cd.shift_id = s.id
		WHERE cd.id = ? LIMIT 1`

	getCashDrawerTransactionsQuery = `
		SELECT t.transaction_date, t.transaction_code, COALESCE(c.name, '') as customer_name, t.total_amount
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE t.user_id = ? AND t.payment_method = 'cash' AND t.status = 'completed'
		  AND t.transaction_date >= ? AND t.transaction_date < COALESCE(?, ?, NOW())
		ORDER BY t.transaction_date ASC`

	getCashDrawerExpensesQuery = `
		SELECT e.category, e.description, e.amount
		FROM expenses e
		WHERE e.user_id = ? AND e.created_at >= ? AND e.created_at < COALESCE(?, ?, NOW())
		ORDER BY e.created_at ASC`

	getNextSessionOpenTimeQuery = `SELECT MIN(open_time) FROM cash_drawer WHERE user_id = ? AND open_time > ? AND id != ?`

	getOpenYesterdayQuery = `SELECT id, user_id, open_time, opening_balance FROM cash_drawer WHERE status = 'open' AND DATE(open_time) < CURDATE()`

	calculateExpectedBalanceQuery = `
		SELECT ? +
		COALESCE((
			SELECT SUM(total_amount) FROM transactions
			WHERE user_id = ? AND payment_method = 'cash' AND status = 'completed'
			  AND transaction_date >= ?
			  AND transaction_date < CONCAT(DATE(?), ' 23:59:59')
		), 0) -
		COALESCE((
			SELECT SUM(amount) FROM expenses
			WHERE user_id = ? AND created_at >= ?
			  AND created_at < CONCAT(DATE(?), ' 23:59:59')
		), 0)`

	autoCloseCashDrawerQuery = `
		UPDATE cash_drawer SET
			close_time       = CONCAT(DATE(open_time), ' 23:59:59'),
			closing_balance  = ?,
			expected_balance = ?,
			difference       = 0,
			status           = 'closed',
			is_auto_closed   = TRUE,
			notes            = 'Ditutup otomatis oleh sistem karena pergantian hari',
			updated_at       = NOW()
		WHERE id = ?`

	getMyCashQuery = `
		SELECT cd.id, cd.user_id, s.name as shift_name,
		       s.start_time as shift_start, s.end_time as shift_end,
		       cd.open_time, cd.opening_balance, cd.total_cash_sales, cd.total_expenses,
		       cd.expected_balance, cd.status, cd.open_notes
		FROM cash_drawer cd
		LEFT JOIN shifts s ON cd.shift_id = s.id
		WHERE cd.user_id = ? AND DATE(cd.open_time) = CURDATE() AND cd.status = 'open'
		LIMIT 1`
)

func (r *cashDrawerRepo) GetCurrent(userID int) (*dto.CurrentCashDrawerResponse, error) {
	var res dto.CurrentCashDrawerResponse
	err := r.db.Raw(getCurrentCashDrawerQuery, userID).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	if res.ID == 0 {
		return nil, nil
	}
	return &res, nil
}

func (r *cashDrawerRepo) GetOpenCashDrawer(userID int) (*model.CashDrawer, error) {
	var dataDB model.CashDrawer
	err := r.db.Raw(getOpenCashDrawerQuery, userID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *cashDrawerRepo) GetByID(id int) (*model.CashDrawer, error) {
	var dataDB model.CashDrawer
	err := r.db.Raw(getCashDrawerByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *cashDrawerRepo) GetHistory(req *dto.GetHistoryRequest) ([]*dto.CashDrawerHistoryResponse, int64, error) {
	var args []any
	conditions := ""

	if req.StartDate != "" {
		conditions += " AND cd.open_time >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		conditions += " AND cd.open_time < DATE_ADD(?, INTERVAL 1 DAY)"
		args = append(args, req.EndDate)
	}
	if req.UserID != nil {
		conditions += " AND cd.user_id = ?"
		args = append(args, *req.UserID)
	}
	if req.ShiftID != nil {
		conditions += " AND cd.shift_id = ?"
		args = append(args, *req.ShiftID)
	}
	if req.Status != "" {
		conditions += " AND cd.status = ?"
		args = append(args, req.Status)
	}

	var total int64
	if err := r.db.Raw(countCashDrawerHistoryBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := getCashDrawerHistoryBase + conditions + " ORDER BY cd.open_time DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*dto.CashDrawerHistoryResponse
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *cashDrawerRepo) GetDetailByID(id int) (*model.CashDrawerDetail, []model.CashDrawerTransactionItem, []model.CashDrawerExpenseItem, error) {
	var dataDB model.CashDrawerDetail
	err := r.db.Raw(getCashDrawerDetailQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, nil, nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil, nil, nil
	}

	var nextOpenTime *string
	r.db.Raw(getNextSessionOpenTimeQuery, dataDB.UserID, dataDB.OpenTime, dataDB.ID).Scan(&nextOpenTime)

	var transactions []model.CashDrawerTransactionItem
	if err := r.db.Raw(getCashDrawerTransactionsQuery, dataDB.UserID, dataDB.OpenTime, dataDB.CloseTime, nextOpenTime).Scan(&transactions).Error; err != nil {
		return nil, nil, nil, err
	}
	if transactions == nil {
		transactions = []model.CashDrawerTransactionItem{}
	}

	var expenses []model.CashDrawerExpenseItem
	if err := r.db.Raw(getCashDrawerExpensesQuery, dataDB.UserID, dataDB.OpenTime, dataDB.CloseTime, nextOpenTime).Scan(&expenses).Error; err != nil {
		return nil, nil, nil, err
	}
	if expenses == nil {
		expenses = []model.CashDrawerExpenseItem{}
	}

	return &dataDB, transactions, expenses, nil
}

func (r *cashDrawerRepo) Open(userID int, shiftID *int, openingBalance float64, notes string) (int64, error) {
	var id int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(openCashDrawerQuery, userID, shiftID, openingBalance, notes).Error; err != nil {
			return err
		}
		return tx.Raw(getLastCashDrawerInsertID).Scan(&id).Error
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *cashDrawerRepo) Close(id int, closingBalance, expectedBalance, difference float64, notes string) error {
	return r.db.Exec(closeCashDrawerQuery, closingBalance, expectedBalance, difference, notes, id).Error
}

func (r *cashDrawerRepo) UpdateSales(id int, totalSales, totalCashSales float64) error {
	return r.db.Exec(updateSalesQuery, totalSales, totalCashSales, id).Error
}

func (r *cashDrawerRepo) UpdateExpenses(id int, totalExpenses float64) error {
	return r.db.Exec(updateExpensesQuery, totalExpenses, id).Error
}

func (r *cashDrawerRepo) AutoCloseYesterday() (int, error) {
	type openDrawer struct {
		ID             int
		UserID         int
		OpenTime       string
		OpeningBalance float64
	}

	var drawers []openDrawer
	if err := r.db.Raw(getOpenYesterdayQuery).Scan(&drawers).Error; err != nil {
		return 0, err
	}
	if len(drawers) == 0 {
		return 0, nil
	}

	count := 0
	for _, cd := range drawers {
		var expected float64
		err := r.db.Raw(
			calculateExpectedBalanceQuery,
			cd.OpeningBalance,
			cd.UserID, cd.OpenTime, cd.OpenTime,
			cd.UserID, cd.OpenTime, cd.OpenTime,
		).Scan(&expected).Error
		if err != nil {
			return count, err
		}

		if err := r.db.Exec(autoCloseCashDrawerQuery, expected, expected, cd.ID).Error; err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (r *cashDrawerRepo) GetMyCash(userID int) (*model.CashDrawerDetail, []model.CashDrawerTransactionItem, []model.CashDrawerExpenseItem, error) {
	var dataDB model.CashDrawerDetail
	err := r.db.Raw(getMyCashQuery, userID).Scan(&dataDB).Error
	if err != nil {
		return nil, nil, nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil, nil, nil
	}

	var nextOpenTime *string
	r.db.Raw(getNextSessionOpenTimeQuery, dataDB.UserID, dataDB.OpenTime, dataDB.ID).Scan(&nextOpenTime)

	var transactions []model.CashDrawerTransactionItem
	if err := r.db.Raw(getCashDrawerTransactionsQuery, dataDB.UserID, dataDB.OpenTime, dataDB.CloseTime, nextOpenTime).Scan(&transactions).Error; err != nil {
		return nil, nil, nil, err
	}
	if transactions == nil {
		transactions = []model.CashDrawerTransactionItem{}
	}

	var expenses []model.CashDrawerExpenseItem
	if err := r.db.Raw(getCashDrawerExpensesQuery, dataDB.UserID, dataDB.OpenTime, dataDB.CloseTime, nextOpenTime).Scan(&expenses).Error; err != nil {
		return nil, nil, nil, err
	}
	if expenses == nil {
		expenses = []model.CashDrawerExpenseItem{}
	}

	return &dataDB, transactions, expenses, nil
}
