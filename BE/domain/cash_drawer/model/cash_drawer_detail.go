package model

import "time"

type CashDrawerTransactionItem struct {
	TransactionDate time.Time `gorm:"column:transaction_date"`
	TransactionCode string    `gorm:"column:transaction_code"`
	CustomerName    string    `gorm:"column:customer_name"`
	TotalAmount     float64   `gorm:"column:total_amount"`
}

type CashDrawerExpenseItem struct {
	Category    string  `gorm:"column:category"`
	Description string  `gorm:"column:description"`
	Amount      float64 `gorm:"column:amount"`
}

type CashDrawerDetail struct {
	ID              int                         `gorm:"column:id"`
	UserID          int                         `gorm:"column:user_id"`
	CashierName     string                      `gorm:"column:cashier_name"`
	ShiftName       *string                     `gorm:"column:shift_name"`
	ShiftStart      *string                     `gorm:"column:shift_start"`
	ShiftEnd        *string                     `gorm:"column:shift_end"`
	OpenTime        time.Time                   `gorm:"column:open_time"`
	CloseTime       *time.Time                  `gorm:"column:close_time"`
	OpeningBalance  float64                     `gorm:"column:opening_balance"`
	ClosingBalance  *float64                    `gorm:"column:closing_balance"`
	ExpectedBalance float64                     `gorm:"column:expected_balance"`
	TotalCashSales  float64                     `gorm:"column:total_cash_sales"`
	TotalExpenses   float64                     `gorm:"column:total_expenses"`
	Difference      *float64                    `gorm:"column:difference"`
	Status          string                      `gorm:"column:status"`
	Notes           *string                     `gorm:"column:notes"`
	OpenNotes       *string                     `gorm:"column:open_notes"`
	Transactions    []CashDrawerTransactionItem
	Expenses        []CashDrawerExpenseItem
}
