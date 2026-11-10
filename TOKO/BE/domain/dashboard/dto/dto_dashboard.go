package dto

type RecentTransactionItem struct {
	ID              int64   `json:"id"`
	TransactionCode string  `json:"transaction_code"`
	TransactionDate string  `json:"transaction_date"`
	TotalAmount     float64 `json:"total_amount"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
}

type TodaySummaryResponse struct {
	TotalTransactions int64   `json:"total_transactions"`
	TotalSales        float64 `json:"total_sales"`
}
