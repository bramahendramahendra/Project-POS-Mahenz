package dto

import "time"

type GetAllRequest struct {
	Page    int    `json:"page" validate:"required,min=1"`
	Limit   int    `json:"limit" validate:"required,min=1"`
	Search  string `json:"search"`
	Status  string `json:"status"`
	DueTerm string `json:"due_term"` // "", "overdue" (lewat jatuh tempo), "not_due" (belum jatuh tempo)
}

type GetByIDRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type PayUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ReceivableResponse struct {
	ID              int        `json:"id"`
	TransactionID   int        `json:"transaction_id"`
	TransactionCode string     `json:"transaction_code"`
	CustomerName    string     `json:"customer_name"`
	TotalAmount     float64    `json:"total_amount"`
	PaidAmount      float64    `json:"paid_amount"`
	RemainingAmount float64    `json:"remaining_amount"`
	Status          string     `json:"status"`
	DueDate         *time.Time `json:"due_date"`
	CreatedAt       time.Time  `json:"created_at"`
}

// TransactionItemBrief: rincian ringkas barang dari transaksi asal piutang —
// cukup untuk tampilan detail & cetak nota.
type TransactionItemBrief struct {
	ProductName  string  `json:"product_name"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
	DiscountItem float64 `json:"discount_item"`
}

type ReceivableDetailResponse struct {
	ID              int                    `json:"id"`
	TransactionID   int                    `json:"transaction_id"`
	TransactionCode string                 `json:"transaction_code"`
	CustomerID      int                    `json:"customer_id"`
	CustomerName    string                 `json:"customer_name"`
	TotalAmount     float64                `json:"total_amount"`
	PaidAmount      float64                `json:"paid_amount"`
	RemainingAmount float64                `json:"remaining_amount"`
	Status          string                 `json:"status"`
	DueDate         *time.Time             `json:"due_date"`
	Notes           string                 `json:"notes"`
	CreatedAt       time.Time              `json:"created_at"`
	Items           []TransactionItemBrief `json:"items" gorm:"-"`
	Payments        []PaymentResponse      `json:"payments" gorm:"-"`
}

type ReceivableSummaryItem struct {
	CustomerID      int     `json:"customer_id"`
	CustomerName    string  `json:"customer_name"`
	TotalReceivable float64 `json:"total_receivable"`
	TotalPaid       float64 `json:"total_paid"`
	TotalRemaining  float64 `json:"total_remaining"`
	Count           int     `json:"count"`
}

// ReceivableStats: angka agregat untuk kartu ringkasan di atas daftar.
type ReceivableStats struct {
	TotalRemaining   float64 `json:"total_remaining"`   // total sisa piutang belum lunas
	CustomerCount    int     `json:"customer_count"`    // jumlah pelanggan yang masih berutang
	OverdueCount     int     `json:"overdue_count"`     // jumlah piutang lewat jatuh tempo
	OverdueRemaining float64 `json:"overdue_remaining"` // total sisa piutang yang lewat jatuh tempo
}

type PaymentResponse struct {
	ID            int       `json:"id"`
	PaymentDate   time.Time `json:"payment_date"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	UserName      string    `json:"user_name"`
	Notes         string    `json:"notes"`
}

type PayRequest struct {
	ID          int     `json:"-"`
	UserID      int     `json:"-"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	PaymentDate string  `json:"payment_date" validate:"required,datetime=2006-01-02"`
	Notes       string  `json:"notes"`
}

type PayResponse struct {
	ReceivableID    int     `json:"receivable_id"`
	PaidAmount      float64 `json:"paid_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Status          string  `json:"status"`
}
