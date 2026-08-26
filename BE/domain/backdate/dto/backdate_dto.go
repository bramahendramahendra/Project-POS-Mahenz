package dto

import "time"

// ===================== REQUEST =====================

type OpenCashDrawerRequest struct {
	UserID         int     `json:"user_id" validate:"required,min=1"`
	Date           string  `json:"date" validate:"required"`
	ShiftID        *int    `json:"shift_id"`
	OpeningBalance float64 `json:"opening_balance" validate:"min=0"`
	Notes          string  `json:"notes"`
}

type CloseCashDrawerRequest struct {
	ID             int     `uri:"id" validate:"required,min=1"`
	ClosingBalance float64 `json:"closing_balance" validate:"min=0"`
	Notes          string  `json:"notes"`
}

type CreateTransactionItemRequest struct {
	ProductID    int     `json:"product_id" validate:"required,min=1"`
	ProductName  string  `json:"product_name" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,min=0.001"`
	Unit         string  `json:"unit" validate:"required"`
	Price        float64 `json:"price" validate:"required,min=0"`
	Subtotal     float64 `json:"subtotal" validate:"required,min=0"`
	DiscountItem float64 `json:"discount_item" validate:"gte=0"`
	UnitID       *int    `json:"unit_id"`
}

type CreateTransactionRequest struct {
	TransactionTime string                         `json:"transaction_time" validate:"required"`
	ShiftID         *int                           `json:"shift_id"`
	Subtotal        float64                        `json:"subtotal" validate:"required,min=0"`
	Discount        float64                        `json:"discount" validate:"gte=0"`
	Tax             float64                        `json:"tax" validate:"gte=0"`
	TotalAmount     float64                        `json:"total_amount" validate:"required,min=0"`
	PaymentMethod   string                         `json:"payment_method" validate:"required,oneof=cash transfer qris card kredit balance"`
	PaymentAmount   float64                        `json:"payment_amount" validate:"min=0"`
	ChangeAmount    float64                        `json:"change_amount"`
	BalanceUsed     float64                        `json:"balance_used" validate:"min=0"`
	CustomerID      *int                           `json:"customer_id"`
	IsCredit        bool                           `json:"is_credit"`
	DeviceSource    string                         `json:"device_source" validate:"required,oneof=desktop web android"`
	Items           []CreateTransactionItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ===================== RESPONSE =====================

type OpenCashDrawerResponse struct {
	ID int `json:"id"`
}

type CurrentCashDrawerResponse struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	UserName        string     `json:"user_name"`
	ShiftID         *int       `json:"shift_id"`
	ShiftName       *string    `json:"shift_name"`
	OpenTime        time.Time  `json:"open_time"`
	OpeningBalance  float64    `json:"opening_balance"`
	TotalSales      float64    `json:"total_sales"`
	TotalCashSales  float64    `json:"total_cash_sales"`
	TotalExpenses   float64    `json:"total_expenses"`
	ExpectedBalance float64    `json:"expected_balance"`
	Status          string     `json:"status"`
	OpenNotes       *string    `json:"open_notes"`
	CreatedBy       *int       `json:"created_by"`
	CreatedByName   string     `json:"created_by_name"`
}

type CloseCashDrawerResponse struct {
	ExpectedBalance float64 `json:"expected_balance"`
	ClosingBalance  float64 `json:"closing_balance"`
	Difference      float64 `json:"difference"`
}

type CreateTransactionResponse struct {
	ID              int       `json:"id"`
	TransactionCode string    `json:"transaction_code"`
	TransactionDate time.Time `json:"transaction_date"`
	TotalAmount     float64   `json:"total_amount"`
}
