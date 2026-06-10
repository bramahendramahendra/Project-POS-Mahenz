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

type ReceivablePayment struct {
	ID            int       `gorm:"column:id"`
	ReceivableID  int       `gorm:"column:receivable_id"`
	PaymentDate   time.Time `gorm:"column:payment_date"`
	Amount        float64   `gorm:"column:amount"`
	PaymentMethod string    `gorm:"column:payment_method"`
	Notes         string    `gorm:"column:notes"`
	UserID        int       `gorm:"column:user_id"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}
