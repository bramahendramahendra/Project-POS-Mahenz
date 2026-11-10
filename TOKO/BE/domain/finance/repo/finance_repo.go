package repo

import (
	request_helper "pos_api/helper/request"
	"fmt"

	dto "pos_api/domain/finance/dto"
)

const (
	getTotalIncomeQuery    = `SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE status = 'completed' AND DATE(transaction_date) BETWEEN ? AND ?`
	getTotalExpenseQuery   = `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE expense_date BETWEEN ? AND ?`
	getTotalReceivableQuery = `SELECT COALESCE(SUM(remaining_amount), 0) FROM receivables WHERE status != 'paid'`
	// HPP dihitung dari ti.purchase_price (snapshot harga beli saat transaksi terjadi), pola
	// identik dengan profitLossQuery di report_repo.go -- "Laba" di sini adalah laba BERSIH
	// (Pendapatan - HPP - Pengeluaran), jadi HPP wajib ikut dikurangkan, bukan cuma Pengeluaran.
	getTotalCOGSQuery = `
		SELECT COALESCE(SUM(ti.quantity * ti.purchase_price), 0)
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE t.status = 'completed' AND DATE(t.transaction_date) BETWEEN ? AND ?`

	countCashflowQuery = `
		SELECT COUNT(*) FROM (
			SELECT id FROM transactions WHERE status = 'completed' AND DATE(transaction_date) BETWEEN ? AND ?
			UNION ALL
			SELECT id FROM expenses WHERE expense_date BETWEEN ? AND ?
		) AS combined`

	countCashflowByTypeIncomeQuery = `SELECT COUNT(*) FROM transactions WHERE status = 'completed' AND DATE(transaction_date) BETWEEN ? AND ?`
	countCashflowByTypeExpenseQuery = `SELECT COUNT(*) FROM expenses WHERE expense_date BETWEEN ? AND ?`

	getCashflowQuery = `
		SELECT id, type, category, amount, description, date FROM (
			SELECT id, 'income' AS type, 'Penjualan' AS category, total_amount AS amount, transaction_code AS description, DATE(transaction_date) AS date
			FROM transactions WHERE status = 'completed' AND DATE(transaction_date) BETWEEN ? AND ?
			UNION ALL
			SELECT id, 'expense' AS type, category, amount, description, expense_date AS date
			FROM expenses WHERE expense_date BETWEEN ? AND ?
		) AS combined
		ORDER BY date DESC, id DESC
		LIMIT ? OFFSET ?`

	getCashflowByTypeIncomeQuery = `
		SELECT id, 'income' AS type, 'Penjualan' AS category, total_amount AS amount, transaction_code AS description, DATE(transaction_date) AS date
		FROM transactions WHERE status = 'completed' AND DATE(transaction_date) BETWEEN ? AND ?
		ORDER BY transaction_date DESC
		LIMIT ? OFFSET ?`

	getCashflowByTypeExpenseQuery = `
		SELECT id, 'expense' AS type, category, amount, description, expense_date AS date
		FROM expenses WHERE expense_date BETWEEN ? AND ?
		ORDER BY expense_date DESC
		LIMIT ? OFFSET ?`
)

func (r *financeRepo) GetSummary(req *dto.GetSummaryRequest) (*dto.SummaryResponse, error) {
	dateFrom := req.DateFrom
	dateTo := req.DateTo

	var totalIncome float64
	if err := r.db.Raw(getTotalIncomeQuery, dateFrom, dateTo).Scan(&totalIncome).Error; err != nil {
		return nil, err
	}

	var totalExpense float64
	if err := r.db.Raw(getTotalExpenseQuery, dateFrom, dateTo).Scan(&totalExpense).Error; err != nil {
		return nil, err
	}

	var totalCOGS float64
	if err := r.db.Raw(getTotalCOGSQuery, dateFrom, dateTo).Scan(&totalCOGS).Error; err != nil {
		return nil, err
	}

	var totalReceivable float64
	if err := r.db.Raw(getTotalReceivableQuery).Scan(&totalReceivable).Error; err != nil {
		return nil, err
	}

	periodLabel := fmt.Sprintf("%s s/d %s", dateFrom, dateTo)

	return &dto.SummaryResponse{
		TotalIncome:     totalIncome,
		TotalExpense:    totalExpense,
		TotalCOGS:       totalCOGS,
		NetProfit:       totalIncome - totalCOGS - totalExpense,
		TotalReceivable: totalReceivable,
		PeriodLabel:     periodLabel,
	}, nil
}

func (r *financeRepo) GetCashflow(req *dto.GetCashflowRequest) ([]dto.CashflowItemResponse, int64, error) {
	dateFrom := req.DateFrom
	dateTo := req.DateTo

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	var total int64
	var items []dto.CashflowItemResponse

	switch req.Type {
	case "income":
		if err := r.db.Raw(countCashflowByTypeIncomeQuery, dateFrom, dateTo).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := r.db.Raw(getCashflowByTypeIncomeQuery, dateFrom, dateTo, limit, offset).Scan(&items).Error; err != nil {
			return nil, 0, err
		}
	case "expense":
		if err := r.db.Raw(countCashflowByTypeExpenseQuery, dateFrom, dateTo).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := r.db.Raw(getCashflowByTypeExpenseQuery, dateFrom, dateTo, limit, offset).Scan(&items).Error; err != nil {
			return nil, 0, err
		}
	default:
		if err := r.db.Raw(countCashflowQuery, dateFrom, dateTo, dateFrom, dateTo).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := r.db.Raw(getCashflowQuery, dateFrom, dateTo, dateFrom, dateTo, limit, offset).Scan(&items).Error; err != nil {
			return nil, 0, err
		}
	}

	if items == nil {
		items = []dto.CashflowItemResponse{}
	}

	return items, total, nil
}
