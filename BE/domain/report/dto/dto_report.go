package dto

// FilterParams adalah parameter filter tanggal untuk semua laporan
type FilterParams struct {
	DateFrom string
	DateTo   string
}

// ─── Sales Report ──────────────────────────────────────────────

type SalesListRequest struct {
	DateFrom      string `json:"date_from"`
	DateTo        string `json:"date_to"`
	PaymentMethod string `json:"payment_method"`
	UserID        *int   `json:"user_id"`
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
}

type SalesItem struct {
	ID              int     `json:"id"`
	TransactionCode string  `json:"transaction_code"`
	TransactionDate string  `json:"transaction_date"`
	CashierName     string  `json:"cashier_name"`
	CustomerName    string  `json:"customer_name"`
	TotalAmount     float64 `json:"total_amount"`
	Discount        float64 `json:"discount"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
}

type SalesSummary struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalRevenue      float64 `json:"total_revenue"`
	AvgPerTransaction float64 `json:"avg_per_transaction"`
	TotalDiscount     float64 `json:"total_discount"`
	TotalTax          float64 `json:"total_tax"`
}

type SalesChartItem struct {
	Label             string  `json:"label"`
	TotalSales        float64 `json:"total_sales"`
	TotalTransactions int     `json:"total_transactions"`
}

type SalesReportResponse struct {
	Summary SalesSummary `json:"summary"`
	Items   []SalesItem  `json:"items"`
}

// ─── Profit/Loss Report ────────────────────────────────────────

type ProfitLossRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type ProfitLossItem struct {
	ProductID     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	QtySold       float64 `json:"qty_sold"`
	PurchasePrice float64 `json:"purchase_price"`
	TotalCOGS     float64 `json:"total_cogs"`
	TotalRevenue  float64 `json:"total_revenue"`
	GrossProfit   float64 `json:"gross_profit"`
}

type ExpenseSummaryItem struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type ProfitLossResponse struct {
	TotalRevenue  float64              `json:"total_revenue"`
	TotalCOGS     float64              `json:"total_cogs"`
	GrossProfit   float64              `json:"gross_profit"`
	TotalExpenses float64              `json:"total_expenses"`
	NetProfit     float64              `json:"net_profit"`
	Items         []ProfitLossItem     `json:"items"`
	Expenses      []ExpenseSummaryItem `json:"expenses"`
}

// ─── Stock Report ──────────────────────────────────────────────

type StockListRequest struct {
	Search     string `json:"search"`
	CategoryID *int   `json:"category_id"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type StockSummaryRequest struct {
	Search     string `json:"search"`
	CategoryID *int   `json:"category_id"`
}

type StockItem struct {
	ID           int     `json:"id"`
	ProductCode  string  `json:"product_code"`
	ProductName  string  `json:"product_name"`
	CategoryName string  `json:"category_name"`
	CurrentStock float64 `json:"current_stock"`
	MinStock     float64 `json:"min_stock"`
	Unit         string  `json:"unit"`
	CostPrice    float64 `json:"cost_price"`
	StockValue   float64 `json:"stock_value"`
	IsLowStock   bool    `json:"is_low_stock"`
}

type StockSummary struct {
	TotalProducts   int     `json:"total_products"`
	LowStockCount   int     `json:"low_stock_count"`
	TotalStockValue float64 `json:"total_stock_value"`
}

type StockReportResponse struct {
	TotalProducts   int         `json:"total_products"`
	LowStockCount   int         `json:"low_stock_count"`
	TotalStockValue float64     `json:"total_stock_value"`
	Items           []StockItem `json:"items"`
}

// ─── Cashier Report ────────────────────────────────────────────

type CashierReportRequest struct {
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type CashierItem struct {
	UserID            int     `json:"user_id"`
	CashierName       string  `json:"cashier_name"`
	TotalTransactions int     `json:"total_transactions"`
	TotalSales        float64 `json:"total_sales"`
	TotalCash         float64 `json:"total_cash"`
	TotalNonCash      float64 `json:"total_non_cash"`
	AvgPerTransaction float64 `json:"avg_per_transaction"`
	VoidCount         int     `json:"void_count"`
}
