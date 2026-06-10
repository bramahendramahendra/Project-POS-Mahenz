package model

import "time"

type CashDrawer struct {
	ID              int        `gorm:"column:id"`
	UserID          int        `gorm:"column:user_id"`
	ShiftID         *int       `gorm:"column:shift_id"`
	OpenTime        time.Time  `gorm:"column:open_time"`
	CloseTime       *time.Time `gorm:"column:close_time"`
	OpeningBalance  float64    `gorm:"column:opening_balance"`
	ClosingBalance  *float64   `gorm:"column:closing_balance"`
	TotalSales      float64    `gorm:"column:total_sales"`
	TotalCashSales  float64    `gorm:"column:total_cash_sales"`
	TotalExpenses   float64    `gorm:"column:total_expenses"`
	ExpectedBalance float64    `gorm:"column:expected_balance"`
	Difference      *float64   `gorm:"column:difference"`
	Status          string     `gorm:"column:status"`
	Notes           *string    `gorm:"column:notes"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}
