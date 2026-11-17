package repo

import (
	"fmt"
	"time"

	dto "pos_api/domain/business_summary/dto"
	product_repo "pos_api/domain/product/repo"
)

const (
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

	// HPP dihitung dari ti.purchase_price (snapshot harga beli saat transaksi terjadi),
	// pola identik dengan profitLossQuery di report_repo.go -- supaya "Laba Kotor" di sini
	// dan di Laporan Laba Rugi sama-sama Pendapatan - HPP, bukan Pendapatan - Pengeluaran.
	rangeCOGSQuery = `
		SELECT COALESCE(SUM(ti.quantity * ti.purchase_price),0) as total_cogs
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'`

	monthCOGSQuery = `
		SELECT COALESCE(SUM(ti.quantity * ti.purchase_price),0) as total_cogs
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE MONTH(t.transaction_date) = ? AND YEAR(t.transaction_date) = ? AND t.status = 'completed'`

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

	lowStockCandidatesQuery   = `SELECT id, min_stock FROM products WHERE is_active = 1`
	openReceivablesCountQuery = `SELECT COUNT(*) FROM receivables WHERE status != 'paid'`
)

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

func (r *businessSummaryRepo) GetCOGSByRange(startDate, endDate string) (float64, error) {
	var total float64
	row := r.db.Raw(rangeCOGSQuery, startDate, endDate).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetCOGSByRange: %w", err)
	}
	return total, nil
}

func (r *businessSummaryRepo) GetMonthCOGS(month int, year int) (float64, error) {
	var total float64
	row := r.db.Raw(monthCOGSQuery, month, year).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetMonthCOGS: %w", err)
	}
	return total, nil
}

// GetStockStatusCounts mengembalikan (jumlah HABIS, jumlah MENIPIS) — terpisah,
// supaya dashboard bisa menampilkan dua kartu berbeda seperti Laporan Stok.
// Memakai BuildStockSummaries/ComputeStockSummary yang sama dengan Laporan Stok.
func (r *businessSummaryRepo) GetStockStatusCounts() (outCount int64, lowCount int64, err error) {
	type candidateRow struct {
		ID       int     `gorm:"column:id"`
		MinStock float64 `gorm:"column:min_stock"`
	}
	var candidates []*candidateRow
	if err := r.db.Raw(lowStockCandidatesQuery).Scan(&candidates).Error; err != nil {
		return 0, 0, fmt.Errorf("GetStockStatusCounts: %w", err)
	}
	if len(candidates) == 0 {
		return 0, 0, nil
	}

	minStockByProduct := make(map[int]float64, len(candidates))
	for _, c := range candidates {
		minStockByProduct[c.ID] = c.MinStock
	}
	summaries, err := product_repo.BuildStockSummaries(r.db, minStockByProduct)
	if err != nil {
		return 0, 0, fmt.Errorf("GetStockStatusCounts: %w", err)
	}

	// BuildStockSummaries mengembalikan hasil utk SEMUA produk berpackage aktif
	// (scope global); difilter lewat candidates (produk is_active=1) agar
	// konsisten dengan Laporan Stok (report_repo.go).
	for _, c := range candidates {
		s, ok := summaries[c.ID]
		if !ok {
			continue
		}
		switch s.Status {
		case "out": // model_product.StockStatusOut
			outCount++
		case "low": // model_product.StockStatusLow
			lowCount++
		}
	}
	return outCount, lowCount, nil
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
