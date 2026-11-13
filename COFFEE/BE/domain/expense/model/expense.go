package model

import "time"

type Expense struct {
	ID            int       `gorm:"column:id"`
	ExpenseDate   string    `gorm:"column:expense_date"`
	Category      string    `gorm:"column:category"`
	Description   string    `gorm:"column:description"`
	Amount        float64   `gorm:"column:amount"`
	PaymentMethod string    `gorm:"column:payment_method"`
	UserID        int       `gorm:"column:user_id"`
	UserName      string    `gorm:"column:user_name"`
	Notes         string    `gorm:"column:notes"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}
