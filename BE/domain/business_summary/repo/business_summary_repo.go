package repo

import (
	"fmt"
	"time"

	dto "pos_api/domain/business_summary/dto"
)

const (
	todayStatsQuery = `
		SELECT COUNT(*) as total_transactions, COALESCE(SUM(total_amount),0) as total_sales,
		       COALESCE(SUM(discount),0) as total_discount
		FROM transactions WHERE DATE(transaction_date) = ? AND status = 'completed'`

	todayExpensesQuery = `SELECT COALESCE(SUM(amount),0) as total FROM expenses WHERE expense_date = ?`

	rangeStatsQuery = `
		SELECT COUNT(*) as total_transactions, COALESCE(SUM(total_amount),0) as total_sales,
		       COALESCE(SUM(discount),0) as total_discount
		FROM transactions WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'`

	rangeExpensesQuery = `SELECT COALESCE(SUM(amount),0) as total FROM expenses WHERE expense_date BETWEEN ? AND ?`

	monthStatsQuery = `
		SELECT COUNT(*) as total_transactions, COALESCE(SUM(total_amount),0) as total_sales
		FROM transactions WHERE MONTH(transaction_date) = ?
		AND YEAR(transaction_date) = ? AND status = 'completed'`

	monthExpensesQuery = `
		SELECT COALESCE(SUM(amount),0) as total FROM expenses
		WHERE MONTH(expense_date) = ? AND YEAR(expense_date) = ?`

	salesTrendQuery = `
		SELECT DATE_FORMAT(transaction_date, '%Y-%m-%d') as label,
		       COALESCE(SUM(total_amount),0) as total_sales,
		       COUNT(*) as total_transactions
		FROM transactions
		WHERE transaction_date >= DATE_SUB(?, INTERVAL ? DAY) AND status = 'completed'
		GROUP BY DATE(transaction_date) ORDER BY label`

	topProductsQuery = `
		SELECT ti.product_id, p.name as product_name,
		       SUM(ti.quantity) as total_qty,
		       SUM(ti.subtotal) as total_value
		FROM transaction_items ti
		LEFT JOIN products p ON ti.product_id = p.id
		LEFT JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY ti.product_id, p.name
		ORDER BY total_qty DESC LIMIT ?`

	topProductsByValueQuery = `
		SELECT ti.product_id, p.name as product_name,
		       SUM(ti.quantity) as total_qty,
		       SUM(ti.subtotal) as total_value
		FROM transaction_items ti
		LEFT JOIN products p ON ti.product_id = p.id
		LEFT JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY ti.product_id, p.name
		ORDER BY total_value DESC LIMIT ?`

	topCategoriesQuery = `
		SELECT c.id as category_id, c.name as category_name,
		       SUM(ti.subtotal) as total_sales
		FROM transaction_items ti
		LEFT JOIN products p ON ti.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY c.id, c.name ORDER BY total_sales DESC LIMIT ?`

	paymentMethodsQuery = `
		SELECT payment_method, COUNT(*) as count, SUM(total_amount) as total
		FROM transactions WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'
		GROUP BY payment_method`

	lowStockCountQuery        = `SELECT COUNT(*) FROM products WHERE stock <= min_stock AND is_active = 1`
	openReceivablesCountQuery = `SELECT COUNT(*) FROM receivables WHERE status != 'paid'`

	highestTransactionQuery = `
		SELECT COALESCE(MAX(total_amount),0) as total_amount, COALESCE(transaction_code,'') as transaction_code
		FROM transactions
		WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'
		ORDER BY total_amount DESC LIMIT 1`

	peakHourQuery = `
		SELECT COALESCE(sub.hour, 0) as hour, COALESCE(sub.count, 0) as count
		FROM (
			SELECT HOUR(transaction_date) as hour, COUNT(*) as count
			FROM transactions
			WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'
			GROUP BY HOUR(transaction_date)
			ORDER BY count DESC LIMIT 1
		) sub RIGHT JOIN (SELECT 1) dummy ON 1=1`

	avgTransactionQuery = `
		SELECT COALESCE(AVG(total_amount),0) as avg_amount, COUNT(*) as total_count
		FROM transactions
		WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'`
)


func (r *businessSummaryRepo) GetTodayStats(date string) (*dto.TodayStats, error) {
	var result dto.TodayStats
	row := r.db.Raw(todayStatsQuery, date).Row()
	if err := row.Scan(&result.TotalTransactions, &result.TotalSales, &result.TotalDiscount); err != nil {
		return nil, fmt.Errorf("GetTodayStats: %w", err)
	}
	return &result, nil
}

func (r *businessSummaryRepo) GetTodayExpenses(date string) (float64, error) {
	var total float64
	row := r.db.Raw(todayExpensesQuery, date).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetTodayExpenses: %w", err)
	}
	return total, nil
}

func (r *businessSummaryRepo) GetStatsByRange(startDate, endDate string) (*dto.TodayStats, error) {
	var result dto.TodayStats
	row := r.db.Raw(rangeStatsQuery, startDate, endDate).Row()
	if err := row.Scan(&result.TotalTransactions, &result.TotalSales, &result.TotalDiscount); err != nil {
		return nil, fmt.Errorf("GetStatsByRange: %w", err)
	}
	return &result, nil
}

func (r *businessSummaryRepo) GetExpensesByRange(startDate, endDate string) (float64, error) {
	var total float64
	row := r.db.Raw(rangeExpensesQuery, startDate, endDate).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetExpensesByRange: %w", err)
	}
	return total, nil
}

func (r *businessSummaryRepo) GetMonthStats(month int, year int) (*dto.MonthStats, error) {
	var result dto.MonthStats
	row := r.db.Raw(monthStatsQuery, month, year).Row()
	if err := row.Scan(&result.TotalTransactions, &result.TotalSales); err != nil {
		return nil, fmt.Errorf("GetMonthStats: %w", err)
	}
	return &result, nil
}

func (r *businessSummaryRepo) GetMonthExpenses(month int, year int) (float64, error) {
	var total float64
	row := r.db.Raw(monthExpensesQuery, month, year).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetMonthExpenses: %w", err)
	}
	return total, nil
}

func (r *businessSummaryRepo) GetLowStockCount() (int64, error) {
	var count int64
	row := r.db.Raw(lowStockCountQuery).Row()
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("GetLowStockCount: %w", err)
	}
	return count, nil
}

func (r *businessSummaryRepo) GetOpenReceivablesCount() (int64, error) {
	var count int64
	row := r.db.Raw(openReceivablesCountQuery).Row()
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("GetOpenReceivablesCount: %w", err)
	}
	return count, nil
}

func (r *businessSummaryRepo) GetSalesTrend(days int, now time.Time) ([]dto.SalesTrendItem, error) {
	rows, err := r.db.Raw(salesTrendQuery, now, days).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetSalesTrend: %w", err)
	}
	defer rows.Close()

	var items []dto.SalesTrendItem
	for rows.Next() {
		var item dto.SalesTrendItem
		if err := rows.Scan(&item.Label, &item.TotalSales, &item.TotalTransactions); err != nil {
			return nil, fmt.Errorf("GetSalesTrend scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *businessSummaryRepo) GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error) {
	q := topProductsQuery
	if filter.SortBy == "value" {
		q = topProductsByValueQuery
	}
	rows, err := r.db.Raw(q, filter.StartDate, filter.EndDate, filter.Limit).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetTopProducts: %w", err)
	}
	defer rows.Close()

	var items []dto.TopProductItem
	for rows.Next() {
		var item dto.TopProductItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.TotalQty, &item.TotalValue); err != nil {
			return nil, fmt.Errorf("GetTopProducts scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *businessSummaryRepo) GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error) {
	rows, err := r.db.Raw(topCategoriesQuery, filter.StartDate, filter.EndDate, filter.Limit).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetTopCategories: %w", err)
	}
	defer rows.Close()

	var items []dto.TopCategoryItem
	for rows.Next() {
		var item dto.TopCategoryItem
		if err := rows.Scan(&item.CategoryID, &item.CategoryName, &item.TotalSales); err != nil {
			return nil, fmt.Errorf("GetTopCategories scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *businessSummaryRepo) GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error) {
	rows, err := r.db.Raw(paymentMethodsQuery, filter.StartDate, filter.EndDate).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetPaymentMethods: %w", err)
	}
	defer rows.Close()

	var items []dto.PaymentMethodItem
	for rows.Next() {
		var item dto.PaymentMethodItem
		if err := rows.Scan(&item.PaymentMethod, &item.Count, &item.Total); err != nil {
			return nil, fmt.Errorf("GetPaymentMethods scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *businessSummaryRepo) GetHighestTransaction(filter dto.DateRangeFilter) (*dto.HighestTransactionItem, error) {
	var result dto.HighestTransactionItem
	row := r.db.Raw(highestTransactionQuery, filter.StartDate, filter.EndDate).Row()
	if err := row.Scan(&result.TotalAmount, &result.TransactionCode); err != nil {
		return nil, fmt.Errorf("GetHighestTransaction: %w", err)
	}
	if result.TransactionCode == "" {
		return nil, nil
	}
	return &result, nil
}

func (r *businessSummaryRepo) GetPeakHour(filter dto.DateRangeFilter) (*dto.PeakHourItem, error) {
	var result dto.PeakHourItem
	row := r.db.Raw(peakHourQuery, filter.StartDate, filter.EndDate).Row()
	if err := row.Scan(&result.Hour, &result.Count); err != nil {
		return nil, fmt.Errorf("GetPeakHour: %w", err)
	}
	if result.Count == 0 {
		return nil, nil
	}
	return &result, nil
}

func (r *businessSummaryRepo) GetAvgTransaction(filter dto.DateRangeFilter) (*dto.AvgTransactionItem, error) {
	var result dto.AvgTransactionItem
	row := r.db.Raw(avgTransactionQuery, filter.StartDate, filter.EndDate).Row()
	if err := row.Scan(&result.AvgAmount, &result.TotalCount); err != nil {
		return nil, fmt.Errorf("GetAvgTransaction: %w", err)
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result, nil
}
