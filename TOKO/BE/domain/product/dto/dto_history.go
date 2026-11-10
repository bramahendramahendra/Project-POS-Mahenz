package dto

import "time"

// ===================== REQUEST =====================

type ProductHistoryRequest struct {
	ID     int `uri:"id" validate:"required,min=1"`
	Page   int `json:"page" validate:"required,min=1"`
	Limit  int `json:"limit" validate:"required,min=1,max=50"`
	Status string `json:"status"` // only for sale-history: "completed" | "void" | ""
}

// ===================== RESPONSE =====================

type PurchaseHistoryItem struct {
	PurchaseDate   time.Time `json:"purchase_date"`
	PurchaseCode   string    `json:"purchase_code"`
	InvoiceNumber  string    `json:"invoice_number"`
	SupplierName   string    `json:"supplier_name"`
	Unit           string    `json:"unit"`
	Quantity       float64   `json:"quantity"`
	PurchasePrice  float64   `json:"purchase_price"`
	Subtotal       float64   `json:"subtotal"`
}

type PurchaseHistorySummary struct {
	TotalNotes    int     `json:"total_notes"`
	TotalQty      float64 `json:"total_qty"`
	TotalValue    float64 `json:"total_value"`
	AveragePrice  float64 `json:"average_price"`
}

type PurchaseHistoryResponse struct {
	Summary PurchaseHistorySummary `json:"summary"`
	Items   []PurchaseHistoryItem  `json:"items"`
	Total   int64                  `json:"total"`
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
}

type SaleHistoryItem struct {
	TransactionDate time.Time `json:"transaction_date"`
	TransactionCode string    `json:"transaction_code"`
	CustomerName    string    `json:"customer_name"`
	Unit            string    `json:"unit"`
	Quantity        float64   `json:"quantity"`
	Price           float64   `json:"price"`
	Subtotal        float64   `json:"subtotal"`
	DiscountItem    float64   `json:"discount_item"`
	Status          string    `json:"status"`
}

type SaleHistorySummary struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalQty          float64 `json:"total_qty"`
	TotalRevenue      float64 `json:"total_revenue"`
	AveragePerDay     float64 `json:"average_per_day"`
}

type SaleHistoryResponse struct {
	Summary SaleHistorySummary `json:"summary"`
	Items   []SaleHistoryItem  `json:"items"`
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	Limit   int                `json:"limit"`
}
