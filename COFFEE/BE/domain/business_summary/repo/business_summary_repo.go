package repo

import (
	"fmt"
	"time"

	dto "coffee_api/domain/business_summary/dto"
	product_repo "coffee_api/domain/product/repo"
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
	// Coffee (keputusan 114, 117): HPP = HPP barang jadi + biaya bahan (terpakai + terbuang x harga beli rata-rata).
	rangeCOGSQuery = `
		SELECT (SELECT COALESCE(SUM(ti.quantity * ti.purchase_price),0)
		        FROM transaction_items ti JOIN transactions t ON ti.transaction_id = t.id
		        WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed')
		     + (SELECT COALESCE(SUM(-m.qty_change * COALESCE(ap.avg_price, 0)), 0)
		        FROM ingredient_mutations m LEFT JOIN (SELECT pi.ingredient_id, SUM(pi.quantity * pi.purchase_price) / NULLIF(SUM(pi.quantity), 0) AS avg_price FROM purchase_items pi JOIN purchases p ON p.id = pi.purchase_id WHERE pi.item_type = 'ingredient' AND p.status <> 'void' GROUP BY pi.ingredient_id) ap ON ap.ingredient_id = m.ingredient_id
		        WHERE m.type IN ('usage', 'waste') AND m.cancelled_at IS NULL AND DATE(m.created_at) BETWEEN ? AND ?) as total_cogs`

	monthCOGSQuery = `
		SELECT (SELECT COALESCE(SUM(ti.quantity * ti.purchase_price),0)
		        FROM transaction_items ti JOIN transactions t ON ti.transaction_id = t.id
		        WHERE MONTH(t.transaction_date) = ? AND YEAR(t.transaction_date) = ? AND t.status = 'completed')
		     + (SELECT COALESCE(SUM(-m.qty_change * COALESCE(ap.avg_price, 0)), 0)
		        FROM ingredient_mutations m LEFT JOIN (SELECT pi.ingredient_id, SUM(pi.quantity * pi.purchase_price) / NULLIF(SUM(pi.quantity), 0) AS avg_price FROM purchase_items pi JOIN purchases p ON p.id = pi.purchase_id WHERE pi.item_type = 'ingredient' AND p.status <> 'void' GROUP BY pi.ingredient_id) ap ON ap.ingredient_id = m.ingredient_id
		        WHERE m.type IN ('usage', 'waste') AND m.cancelled_at IS NULL AND MONTH(m.created_at) = ? AND YEAR(m.created_at) = ?) as total_cogs`

	salesTrendQuery = `
		SELECT DATE_FORMAT(transaction_date, '%Y-%m-%d') as label,
		       COALESCE(SUM(total_amount),0) as total_sales,
		       COUNT(*) as total_transactions
		FROM transactions
		WHERE transaction_date >= DATE_SUB(?, INTERVAL ? DAY) AND status = 'completed'
		GROUP BY DATE(transaction_date) ORDER BY label`

	// Coffee (keputusan 117): menu terlaris (racikan dan barang jadi), dikelompokkan per menu.
	topProductsQuery = `
		SELECT COALESCE(ti.menu_item_id, 0) as product_id, MAX(ti.product_name) as product_name,
		       SUM(ti.quantity) as total_qty,
		       SUM(ti.subtotal) as total_value
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY COALESCE(ti.menu_item_id, 0)
		ORDER BY total_qty DESC LIMIT ?`

	topProductsByValueQuery = `
		SELECT COALESCE(ti.menu_item_id, 0) as product_id, MAX(ti.product_name) as product_name,
		       SUM(ti.quantity) as total_qty,
		       SUM(ti.subtotal) as total_value
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY COALESCE(ti.menu_item_id, 0)
		ORDER BY total_value DESC LIMIT ?`

	// Coffee (keputusan 117): kategori menu.
	topCategoriesQuery = `
		SELECT COALESCE(mc.id, 0) as category_id, COALESCE(mc.name, 'Tanpa kategori') as category_name,
		       SUM(ti.subtotal) as total_sales
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		LEFT JOIN menu_items mi ON mi.id = ti.menu_item_id
		LEFT JOIN menu_categories mc ON mc.id = mi.category_id
		WHERE DATE(t.transaction_date) BETWEEN ? AND ? AND t.status = 'completed'
		GROUP BY mc.id, mc.name ORDER BY total_sales DESC LIMIT ?`

	paymentMethodsQuery = `
		SELECT payment_method, COUNT(*) as count, SUM(total_amount) as total
		FROM transactions WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'
		GROUP BY payment_method`

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
	row := r.db.Raw(rangeCOGSQuery, startDate, endDate, startDate, endDate).Row()
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetCOGSByRange: %w", err)
	}
	return total, nil
}

func (r *businessSummaryRepo) GetMonthCOGS(month int, year int) (float64, error) {
	var total float64
	row := r.db.Raw(monthCOGSQuery, month, year, month, year).Row()
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
