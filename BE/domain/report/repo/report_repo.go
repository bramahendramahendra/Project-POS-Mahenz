package repo

import (
	dto "pos_api/domain/report/dto"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"
)

const (
	salesReportQuery = `
		SELECT t.id, t.transaction_code, DATE_FORMAT(t.transaction_date, '%Y-%m-%d %H:%i:%s') as transaction_date,
		       u.full_name as cashier_name, COALESCE(cu.name,'') as customer_name,
		       t.total_amount, t.discount, t.payment_method, t.status
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN customers cu ON t.customer_id = cu.id
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?
		ORDER BY t.transaction_date DESC`

	salesListBase = `
		SELECT t.id, t.transaction_code, DATE_FORMAT(t.transaction_date, '%Y-%m-%d %H:%i:%s') as transaction_date,
		       u.full_name as cashier_name, COALESCE(cu.name,'') as customer_name,
		       t.total_amount, t.discount, t.payment_method, t.status
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN customers cu ON t.customer_id = cu.id
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?`

	salesListCountBase = `
		SELECT COUNT(*) FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN customers cu ON t.customer_id = cu.id
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?`

	salesSummaryQuery = `
		SELECT COUNT(*) as total_transactions,
		       COALESCE(SUM(total_amount),0) as total_revenue,
		       COALESCE(AVG(total_amount),0) as avg_per_transaction,
		       COALESCE(SUM(discount),0) as total_discount, COALESCE(SUM(tax),0) as total_tax
		FROM transactions WHERE status = 'completed' AND transaction_date BETWEEN ? AND ?`

	salesSummaryBase = `
		SELECT COUNT(*) as total_transactions,
		       COALESCE(SUM(t.total_amount),0) as total_revenue,
		       COALESCE(AVG(t.total_amount),0) as avg_per_transaction,
		       COALESCE(SUM(t.discount),0) as total_discount, COALESCE(SUM(t.tax),0) as total_tax
		FROM transactions t
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?`

	salesChartQuery = `
		SELECT DATE(transaction_date) as label, COALESCE(SUM(total_amount),0) as total_sales, COUNT(*) as total_transactions
		FROM transactions WHERE status = 'completed' AND transaction_date BETWEEN ? AND ?
		GROUP BY DATE(transaction_date) ORDER BY label`

	// total_revenue diprorata dari diskon level-transaksi (header discount) supaya net dari
	// SEMUA diskon (item + header), tapi tetap exclude pajak (pajak dititipkan ke negara,
	// bukan pendapatan bisnis). Basis: ti.subtotal * (t.total_amount - t.tax) / t.subtotal,
	// yang secara aljabar sama dengan ti.subtotal * (t.subtotal - t.discount) / t.subtotal
	// karena total_amount = subtotal - discount + tax. NULLIF menghindari divide-by-zero.
	profitLossQuery = `
		SELECT ti.product_id, p.name as product_name,
		       SUM(ti.quantity) as qty_sold,
		       p.purchase_price,
		       COALESCE(SUM(ti.quantity * p.purchase_price),0) as total_cogs,
		       COALESCE(SUM(ti.subtotal * (t.total_amount - t.tax) / NULLIF(t.subtotal,0)),0) as total_revenue
		FROM transaction_items ti
		LEFT JOIN products p ON ti.product_id = p.id
		LEFT JOIN transactions t ON ti.transaction_id = t.id
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?
		GROUP BY ti.product_id, p.name, p.purchase_price`

	expenseSummaryQuery = `
		SELECT category, COALESCE(SUM(amount),0) as total
		FROM expenses WHERE expense_date BETWEEN ? AND ?
		GROUP BY category`

	stockReportQuery = `
		SELECT p.id, COALESCE(p.sku,'') as product_code, p.name as product_name,
		       COALESCE(c.name,'') as category_name, p.stock as current_stock, p.min_stock,
		       COALESCE(u.name,'') as unit, p.purchase_price as cost_price,
		       (p.stock * p.purchase_price) as stock_value,
		       CASE WHEN p.stock <= p.min_stock THEN 1 ELSE 0 END as is_low_stock
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1`

	stockListBase = `
		SELECT p.id, COALESCE(p.sku,'') as product_code, p.name as product_name,
		       COALESCE(c.name,'') as category_name, p.stock as current_stock, p.min_stock,
		       COALESCE(u.name,'') as unit, p.purchase_price as cost_price,
		       (p.stock * p.purchase_price) as stock_value,
		       CASE WHEN p.stock <= p.min_stock THEN 1 ELSE 0 END as is_low_stock
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1`

	stockListCountBase = `
		SELECT COUNT(*) FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_active = 1`

	stockSummaryBase = `
		SELECT COUNT(*) as total_products,
		       SUM(CASE WHEN p.stock <= p.min_stock THEN 1 ELSE 0 END) as low_stock_count,
		       COALESCE(SUM(p.stock * p.purchase_price),0) as total_stock_value
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_active = 1`

	cashierReportBase = `
		SELECT t.user_id, u.full_name as cashier_name,
		       COUNT(CASE WHEN t.status = 'completed' THEN 1 END) as total_transactions,
		       COALESCE(SUM(CASE WHEN t.status = 'completed' THEN t.total_amount ELSE 0 END),0) as total_sales,
		       COALESCE(SUM(CASE WHEN t.status = 'completed' AND t.payment_method='cash' THEN t.total_amount ELSE 0 END),0) as total_cash,
		       COALESCE(SUM(CASE WHEN t.status = 'completed' AND t.payment_method!='cash' THEN t.total_amount ELSE 0 END),0) as total_non_cash,
		       COALESCE(AVG(CASE WHEN t.status = 'completed' THEN t.total_amount END),0) as avg_per_transaction,
		       COUNT(CASE WHEN t.status = 'void' THEN 1 END) as void_count
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.transaction_date BETWEEN ? AND ?`

	cashierReportGroupBy = ` GROUP BY t.user_id, u.full_name`

	cashierCountBase = `
		SELECT COUNT(DISTINCT t.user_id)
		FROM transactions t
		WHERE t.transaction_date BETWEEN ? AND ?`
)

func (r *reportRepo) GetSalesItems(params dto.FilterParams) ([]dto.SalesItem, error) {
	rows, err := r.db.Raw(salesReportQuery, params.DateFrom, params.DateTo).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.SalesItem
	for rows.Next() {
		var item dto.SalesItem
		if err := rows.Scan(&item.ID, &item.TransactionCode, &item.TransactionDate,
			&item.CashierName, &item.CustomerName, &item.TotalAmount, &item.Discount, &item.PaymentMethod, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.SalesItem{}
	}
	return items, nil
}

func (r *reportRepo) GetSalesItemsPaginated(req *dto.SalesListRequest) ([]dto.SalesItem, int64, error) {
	dateFrom, dateTo := resolveSalesDates(req.DateFrom, req.DateTo)
	args := []any{dateFrom, dateTo}
	conditions := ""

	if req.PaymentMethod != "" {
		conditions += " AND t.payment_method = ?"
		args = append(args, req.PaymentMethod)
	}
	if req.UserID != nil {
		conditions += " AND t.user_id = ?"
		args = append(args, *req.UserID)
	}

	var total int64
	r.db.Raw(salesListCountBase+conditions, args...).Scan(&total)

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortFields := map[string]string{
		"transaction_date": "t.transaction_date",
		"transaction_code": "t.transaction_code",
		"cashier_name":     "u.full_name",
		"customer_name":    "cu.name",
		"total_amount":     "t.total_amount",
		"payment_method":   "t.payment_method",
		"status":           "t.status",
	}
	const salesListDefaultOrder = " ORDER BY t.transaction_date DESC"
	listArgs := append(args, limit, offset)
	query := salesListBase + conditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, salesListDefaultOrder) +
		" LIMIT ? OFFSET ?"

	rows, err := r.db.Raw(query, listArgs...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.SalesItem
	for rows.Next() {
		var item dto.SalesItem
		if err := rows.Scan(&item.ID, &item.TransactionCode, &item.TransactionDate,
			&item.CashierName, &item.CustomerName, &item.TotalAmount, &item.Discount, &item.PaymentMethod, &item.Status); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.SalesItem{}
	}
	return items, total, nil
}

func (r *reportRepo) GetSalesSummaryWithFilters(req *dto.SalesListRequest) (*dto.SalesSummary, error) {
	dateFrom, dateTo := resolveSalesDates(req.DateFrom, req.DateTo)
	args := []any{dateFrom, dateTo}
	conditions := ""

	if req.PaymentMethod != "" {
		conditions += " AND t.payment_method = ?"
		args = append(args, req.PaymentMethod)
	}
	if req.UserID != nil {
		conditions += " AND t.user_id = ?"
		args = append(args, *req.UserID)
	}

	var summary dto.SalesSummary
	if err := r.db.Raw(salesSummaryBase+conditions, args...).Scan(&summary).Error; err != nil {
		return nil, err
	}
	return &summary, nil
}

func resolveSalesDates(dateFrom, dateTo string) (string, string) {
	return time_helper.NormalizeDateRange(dateFrom, dateTo)
}

func (r *reportRepo) GetSalesSummary(params dto.FilterParams) (*dto.SalesSummary, error) {
	var summary dto.SalesSummary
	if err := r.db.Raw(salesSummaryQuery, params.DateFrom, params.DateTo).Scan(&summary).Error; err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *reportRepo) GetSalesChart(params dto.FilterParams) ([]dto.SalesChartItem, error) {
	rows, err := r.db.Raw(salesChartQuery, params.DateFrom, params.DateTo).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.SalesChartItem
	for rows.Next() {
		var item dto.SalesChartItem
		if err := rows.Scan(&item.Label, &item.TotalSales, &item.TotalTransactions); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.SalesChartItem{}
	}
	return items, nil
}

func (r *reportRepo) GetProfitLossItems(params dto.FilterParams) ([]dto.ProfitLossItem, error) {
	rows, err := r.db.Raw(profitLossQuery, params.DateFrom, params.DateTo).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.ProfitLossItem
	for rows.Next() {
		var item dto.ProfitLossItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.QtySold,
			&item.PurchasePrice, &item.TotalCOGS, &item.TotalRevenue); err != nil {
			return nil, err
		}
		item.GrossProfit = item.TotalRevenue - item.TotalCOGS
		items = append(items, item)
	}
	if items == nil {
		items = []dto.ProfitLossItem{}
	}
	return items, nil
}

func (r *reportRepo) GetExpenseSummary(params dto.FilterParams) ([]dto.ExpenseSummaryItem, error) {
	rows, err := r.db.Raw(expenseSummaryQuery, params.DateFrom, params.DateTo).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.ExpenseSummaryItem
	for rows.Next() {
		var item dto.ExpenseSummaryItem
		if err := rows.Scan(&item.Category, &item.Total); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.ExpenseSummaryItem{}
	}
	return items, nil
}

func (r *reportRepo) GetStockItems() ([]dto.StockItem, error) {
	rows, err := r.db.Raw(stockReportQuery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.StockItem
	for rows.Next() {
		var item dto.StockItem
		var isLowInt int
		if err := rows.Scan(&item.ID, &item.ProductCode, &item.ProductName, &item.CategoryName,
			&item.CurrentStock, &item.MinStock, &item.Unit, &item.CostPrice, &item.StockValue, &isLowInt); err != nil {
			return nil, err
		}
		item.IsLowStock = isLowInt == 1
		items = append(items, item)
	}
	if items == nil {
		items = []dto.StockItem{}
	}
	return items, nil
}

func (r *reportRepo) GetStockItemsPaginated(req *dto.StockListRequest) ([]dto.StockItem, int64, error) {
	args := []any{}
	listConditions := ""
	countConditions := ""

	if req.Search != "" {
		listConditions += " AND (p.name LIKE ? OR p.sku LIKE ?)"
		countConditions += " AND (p.name LIKE ? OR p.sku LIKE ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like)
	}
	if req.CategoryID != nil {
		listConditions += " AND p.category_id = ?"
		countConditions += " AND p.category_id = ?"
		args = append(args, *req.CategoryID)
	}

	var total int64
	r.db.Raw(stockListCountBase+countConditions, args...).Scan(&total)

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedStockSortFields := map[string]string{
		"product_code":  "p.sku",
		"product_name":  "p.name",
		"category_name": "c.name",
		"current_stock": "p.stock",
		"stock_value":   "(p.stock * p.purchase_price)",
	}
	const stockListDefaultOrder = " ORDER BY p.name ASC"
	listArgs := append(args, limit, offset)
	query := stockListBase + listConditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedStockSortFields, stockListDefaultOrder) +
		" LIMIT ? OFFSET ?"

	rows, err := r.db.Raw(query, listArgs...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.StockItem
	for rows.Next() {
		var item dto.StockItem
		var isLowInt int
		if err := rows.Scan(&item.ID, &item.ProductCode, &item.ProductName, &item.CategoryName,
			&item.CurrentStock, &item.MinStock, &item.Unit, &item.CostPrice, &item.StockValue, &isLowInt); err != nil {
			return nil, 0, err
		}
		item.IsLowStock = isLowInt == 1
		items = append(items, item)
	}
	if items == nil {
		items = []dto.StockItem{}
	}
	return items, total, nil
}

func (r *reportRepo) GetStockSummaryWithFilters(req *dto.StockSummaryRequest) (*dto.StockSummary, error) {
	args := []any{}
	conditions := ""

	if req.Search != "" {
		conditions += " AND (p.name LIKE ? OR p.sku LIKE ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like)
	}
	if req.CategoryID != nil {
		conditions += " AND p.category_id = ?"
		args = append(args, *req.CategoryID)
	}

	var summary dto.StockSummary
	if err := r.db.Raw(stockSummaryBase+conditions, args...).Scan(&summary).Error; err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *reportRepo) GetCashierItems(params dto.FilterParams) ([]dto.CashierItem, error) {
	query := cashierReportBase + cashierReportGroupBy + " ORDER BY total_sales DESC"
	rows, err := r.db.Raw(query, params.DateFrom, params.DateTo).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.CashierItem
	for rows.Next() {
		var item dto.CashierItem
		if err := rows.Scan(&item.UserID, &item.CashierName, &item.TotalTransactions,
			&item.TotalSales, &item.TotalCash, &item.TotalNonCash, &item.AvgPerTransaction, &item.VoidCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.CashierItem{}
	}
	return items, nil
}

func (r *reportRepo) GetCashierItemsPaginated(req *dto.CashierReportRequest) ([]dto.CashierItem, int64, error) {
	dateFrom, dateTo := time_helper.NormalizeDateRange(req.DateFrom, req.DateTo)

	var total int64
	if err := r.db.Raw(cashierCountBase, dateFrom, dateTo).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortColumns := map[string]string{
		"cashier_name":        "cashier_name",
		"total_transactions":  "total_transactions",
		"total_sales":         "total_sales",
		"avg_per_transaction": "avg_per_transaction",
		"void_count":          "void_count",
	}
	sortCol := "total_sales"
	if col, ok := allowedSortColumns[req.SortBy]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if req.SortOrder == "asc" {
		sortDir = "ASC"
	}

	query := cashierReportBase + cashierReportGroupBy +
		" ORDER BY " + sortCol + " " + sortDir + " LIMIT ? OFFSET ?"
	rows, err := r.db.Raw(query, dateFrom, dateTo, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.CashierItem
	for rows.Next() {
		var item dto.CashierItem
		if err := rows.Scan(&item.UserID, &item.CashierName, &item.TotalTransactions,
			&item.TotalSales, &item.TotalCash, &item.TotalNonCash, &item.AvgPerTransaction, &item.VoidCount); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.CashierItem{}
	}
	return items, total, nil
}
