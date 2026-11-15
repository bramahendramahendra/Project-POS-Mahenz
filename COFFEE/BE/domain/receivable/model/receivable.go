package model

import "time"

type Receivable struct {
	ID              int        `gorm:"column:id"`
	TransactionID   int        `gorm:"column:transaction_id"`
	CustomerID      int        `gorm:"column:customer_id"`
	TotalAmount     float64    `gorm:"column:total_amount"`
	PaidAmount      float64    `gorm:"column:paid_amount"`
	RemainingAmount float64    `gorm:"column:remaining_amount"`
	Status          string     `gorm:"column:status"`
	DueDate         *time.Time `gorm:"column:due_date"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
}
