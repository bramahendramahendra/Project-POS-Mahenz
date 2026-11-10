package repo

import (
	"fmt"

	dto "pos_api/domain/dashboard/dto"
)

const (
	recentTransactionsQuery = `
		SELECT id, transaction_code, transaction_date, total_amount, payment_method, status
		FROM transactions
		WHERE user_id = ? AND status != 'pending'
		ORDER BY transaction_date DESC LIMIT ?`

	todaySummaryQuery = `
		SELECT COUNT(*) as total_transactions, COALESCE(SUM(total_amount),0) as total_sales
		FROM transactions
		WHERE user_id = ? AND DATE(transaction_date) = ? AND status = 'completed'`
)

func (r *dashboardRepo) GetRecentTransactions(userID int, limit int) ([]dto.RecentTransactionItem, error) {
	rows, err := r.db.Raw(recentTransactionsQuery, userID, limit).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetRecentTransactions: %w", err)
	}
	defer rows.Close()

	var items []dto.RecentTransactionItem
	for rows.Next() {
		var item dto.RecentTransactionItem
		if err := rows.Scan(&item.ID, &item.TransactionCode, &item.TransactionDate, &item.TotalAmount, &item.PaymentMethod, &item.Status); err != nil {
			return nil, fmt.Errorf("GetRecentTransactions scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *dashboardRepo) GetTodaySummary(userID int, date string) (*dto.TodaySummaryResponse, error) {
	var result dto.TodaySummaryResponse
	row := r.db.Raw(todaySummaryQuery, userID, date).Row()
	if err := row.Scan(&result.TotalTransactions, &result.TotalSales); err != nil {
		return nil, fmt.Errorf("GetTodaySummary: %w", err)
	}
	return &result, nil
}
