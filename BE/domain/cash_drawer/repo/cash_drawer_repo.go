package repo

import (
	request_helper "pos_api/helper/request"
	dto "pos_api/domain/cash_drawer/dto"
	model "pos_api/domain/cash_drawer/model"

	"gorm.io/gorm"
)

// liveExpectedBalanceExpr menghitung expected_balance LANGSUNG dari data sumber
// (transaksi tunai + pengeluaran sejak kas dibuka), bukan dari kolom cd.expected_balance
// yang di-cache dan harus di-update manual di setiap tempat yang menyentuh total kas.
// Kolom cache itu defaultnya 0 saat kas baru dibuka dan baru terisi saat UpdateSales/
// UpdateExpenses pertama kali dipanggil — kalau kas dibuka lalu ditutup tanpa ada
// transaksi/pengeluaran sama sekali, kolom itu tetap 0 (bukan opening_balance), membuat
// "difference" salah besar. Dengan dihitung langsung dari tabel sumber tiap kali dibaca,
// nilainya TIDAK PERNAH bisa basi/lupa di-update — pola ini sudah dipakai lebih dulu oleh
// AutoCloseYesterday() (lihat calculateExpectedBalanceQuery), di sini digeneralisasi jadi
// ekspresi inline yang bisa dipakai di semua query baca. Hanya berlaku untuk kas yang MASIH
// TERBUKA -- kas yang sudah ditutup tetap pakai nilai beku dari kolom (snapshot historis
// saat penutupan, sengaja tidak boleh berubah lagi walau ada transaksi terkait di-void
// belakangan).
const liveExpectedBalanceExpr = `(cd.opening_balance + COALESCE((SELECT SUM(total_amount) FROM transactions WHERE user_id = cd.user_id AND payment_method = 'cash' AND status = 'completed' AND transaction_date >= cd.open_time), 0) - COALESCE((SELECT SUM(amount) FROM expenses WHERE user_id = cd.user_id AND created_at >= cd.open_time), 0))`

const (
	getCurrentCashDrawerQuery = `SELECT cd.id, cd.user_id, u.full_name as user_name, cd.shift_id, s.name as shift_name, s.start_time as shift_start, s.end_time as shift_end, cd.open_time, cd.opening_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses, ` + liveExpectedBalanceExpr + ` as expected_balance, cd.status, cd.open_notes FROM cash_drawer cd LEFT JOIN users u ON cd.user_id = u.id LEFT JOIN shifts s ON cd.shift_id = s.id WHERE cd.user_id = ? AND cd.status = 'open' LIMIT 1`
	getOpenCashDrawerQuery     = `SELECT cd.id, cd.user_id, cd.shift_id, cd.open_time, cd.opening_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses, ` + liveExpectedBalanceExpr + ` as expected_balance, cd.status FROM cash_drawer cd WHERE cd.user_id = ? AND cd.status = 'open' LIMIT 1`
	getCashDrawerByIDQuery     = `SELECT cd.id, cd.user_id, cd.shift_id, cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance, cd.total_sales, cd.total_cash_sales, cd.total_expenses, CASE WHEN cd.status = 'closed' THEN cd.expected_balance ELSE ` + liveExpectedBalanceExpr + ` END as expected_balance, cd.difference, cd.status, cd.notes FROM cash_drawer cd WHERE cd.id = ? LIMIT 1`
	openCashDrawerQuery        = `INSERT INTO cash_drawer (user_id, shift_id, open_time, opening_balance, open_notes, status) VALUES (?, ?, NOW(), ?, ?, 'open')`
	getLastCashDrawerInsertID  = `SELECT LAST_INSERT_ID()`
	closeCashDrawerQuery       = `UPDATE cash_drawer SET close_time = NOW(), closing_balance = ?, expected_balance = ?, difference = ?, status = 'closed', notes = ?, updated_at = NOW() WHERE id = ?`
	updateSalesQuery           = `UPDATE cash_drawer SET total_sales = total_sales + ?, total_cash_sales = total_cash_sales + ?, expected_balance = opening_balance + total_cash_sales - total_expenses, updated_at = NOW() WHERE id = ?`
	updateExpensesQuery        = `UPDATE cash_drawer SET total_expenses = total_expenses + ?, expected_balance = opening_balance + total_cash_sales - total_expenses, updated_at = NOW() WHERE id = ?`
	countCashDrawerHistoryBase = `SELECT COUNT(*) FROM cash_drawer cd WHERE 1=1`
	getKasirOptionsQuery       = `SELECT DISTINCT u.id, u.full_name, u.username FROM cash_drawer cd JOIN users u ON cd.user_id = u.id ORDER BY u.full_name`
)

var getCashDrawerHistoryBase = `SELECT cd.id, u.full_name as user_name, s.name as shift_name, cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance, CASE WHEN cd.status = 'closed' THEN cd.expected_balance ELSE ` + liveExpectedBalanceExpr + ` END as expected_balance, CASE WHEN cd.status = 'closed' THEN cd.difference ELSE NULL END as difference, cd.total_sales, cd.total_cash_sales, cd.total_expenses, cd.status FROM cash_drawer cd LEFT JOIN users u ON cd.user_id = u.id LEFT JOIN shifts s ON cd.shift_id = s.id WHERE 1=1`

var getCashDrawerDetailQuery = `
	SELECT cd.id, cd.user_id, u.full_name as cashier_name, s.name as shift_name,
	       s.start_time as shift_start, s.end_time as shift_end,
	       cd.open_time, cd.close_time, cd.opening_balance, cd.closing_balance,
	       CASE WHEN cd.status = 'closed' THEN cd.expected_balance ELSE ` + liveExpectedBalanceExpr + ` END as expected_balance,
	       cd.total_cash_sales, cd.total_expenses,
	       CASE WHEN cd.status = 'closed' THEN cd.difference ELSE NULL END as difference,
	       cd.status, cd.notes, cd.open_notes
	FROM cash_drawer cd
	LEFT JOIN users u ON cd.user_id = u.id
	LEFT JOIN shifts s ON cd.shift_id = s.id
	WHERE cd.id = ? LIMIT 1`

const (

	getCashDrawerTransactionsQuery = `
		SELECT t.transaction_date, t.transaction_code, COALESCE(c.name, '') as customer_name, t.total_amount
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE t.user_id = ? AND t.payment_method = 'cash' AND t.status = 'completed'
		  AND t.transaction_date >= ? AND t.transaction_date < COALESCE(?, ?, NOW())
		ORDER BY t.transaction_date ASC`

	getNonCashTransactionsQuery = `
		SELECT t.transaction_date, t.transaction_code, COALESCE(c.name, '') as customer_name, pm.label as payment_method_label, t.total_amount
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		JOIN payment_methods pm ON t.payment_method = pm.code
		WHERE t.user_id = ? AND t.payment_method != 'cash' AND t.status = 'completed'
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

	getNonCashSalesQuery = `
		SELECT t.payment_method, pm.label, COALESCE(SUM(t.total_amount), 0) as total
		FROM transactions t
		JOIN payment_methods pm ON t.payment_method = pm.code
		WHERE t.user_id = ?
		  AND t.status = 'completed'
		  AND t.payment_method != 'cash'
		  AND t.transaction_date >= ?
		  AND t.transaction_date < COALESCE(?, NOW())
		GROUP BY t.payment_method, pm.label
		ORDER BY total DESC`

	getCashDrawerSummaryAggregateQuery = `SELECT COALESCE(SUM(cd.opening_balance), 0) as total_opening, COALESCE(SUM(CASE WHEN cd.status = 'closed' THEN cd.closing_balance ELSE 0 END), 0) as total_closing, COALESCE(SUM(cd.total_expenses), 0) as total_expenses, COALESCE(SUM(CASE WHEN cd.status = 'closed' THEN cd.difference ELSE 0 END), 0) as total_difference FROM cash_drawer cd WHERE 1=1`

)

var getMyCashQuery = `
	SELECT cd.id, cd.user_id, s.name as shift_name,
	       s.start_time as shift_start, s.end_time as shift_end,
	       cd.open_time, cd.opening_balance, cd.total_cash_sales, cd.total_expenses,
	       ` + liveExpectedBalanceExpr + ` as expected_balance, cd.status, cd.open_notes
	FROM cash_drawer cd
	LEFT JOIN shifts s ON cd.shift_id = s.id
	WHERE cd.user_id = ? AND DATE(cd.open_time) = CURDATE() AND cd.status = 'open'
	LIMIT 1`

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

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortColumns := map[string]string{
		"open_time":        "cd.open_time",
		"user_name":        "u.full_name",
		"opening_balance":  "cd.opening_balance",
		"total_cash_sales": "cd.total_cash_sales",
		"total_expenses":   "cd.total_expenses",
	}
	sortCol := "cd.open_time"
	if col, ok := allowedSortColumns[req.SortBy]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if req.SortOrder == "asc" {
		sortDir = "ASC"
	}

	query := getCashDrawerHistoryBase + conditions + " ORDER BY " + sortCol + " " + sortDir + " LIMIT ? OFFSET ?"
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

func (r *cashDrawerRepo) GetNonCashTransactions(userID int, openTime string, closeTime *string, nextOpenTime *string) ([]model.CashDrawerNonCashTransactionItem, error) {
	var result []model.CashDrawerNonCashTransactionItem
	if err := r.db.Raw(getNonCashTransactionsQuery, userID, openTime, closeTime, nextOpenTime).Scan(&result).Error; err != nil {
		return nil, err
	}
	if result == nil {
		result = []model.CashDrawerNonCashTransactionItem{}
	}
	return result, nil
}

func (r *cashDrawerRepo) GetNonCashSales(userID int, openTime string, closeTime *string) ([]dto.NonCashSaleItem, error) {
	var result []dto.NonCashSaleItem
	if err := r.db.Raw(getNonCashSalesQuery, userID, openTime, closeTime).Scan(&result).Error; err != nil {
		return nil, err
	}
	if result == nil {
		result = []dto.NonCashSaleItem{}
	}
	return result, nil
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

func (r *cashDrawerRepo) GetSummary(req *dto.GetHistoryRequest) (*dto.CashDrawerSummaryResponse, error) {
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

	type aggregateResult struct {
		TotalOpening    float64 `gorm:"column:total_opening"`
		TotalClosing    float64 `gorm:"column:total_closing"`
		TotalExpenses   float64 `gorm:"column:total_expenses"`
		TotalDifference float64 `gorm:"column:total_difference"`
	}
	var agg aggregateResult
	if err := r.db.Raw(getCashDrawerSummaryAggregateQuery+conditions, args...).Scan(&agg).Error; err != nil {
		return nil, err
	}

	recordQuery := getCashDrawerHistoryBase + conditions + " ORDER BY cd.open_time DESC LIMIT 1000"
	var records []*dto.CashDrawerHistoryResponse
	if err := r.db.Raw(recordQuery, args...).Scan(&records).Error; err != nil {
		return nil, err
	}
	if records == nil {
		records = []*dto.CashDrawerHistoryResponse{}
	}

	return &dto.CashDrawerSummaryResponse{
		TotalOpening:  agg.TotalOpening,
		TotalClosing:  agg.TotalClosing,
		TotalExpenses: agg.TotalExpenses,
		Net:           agg.TotalDifference,
		Records:       records,
	}, nil
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

func (r *cashDrawerRepo) GetKasirOptions() ([]dto.KasirOptionResponse, error) {
	var options []dto.KasirOptionResponse
	if err := r.db.Raw(getKasirOptionsQuery).Scan(&options).Error; err != nil {
		return nil, err
	}
	if options == nil {
		options = []dto.KasirOptionResponse{}
	}
	return options, nil
}
