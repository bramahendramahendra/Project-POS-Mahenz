package model

import "time"

type CustomerBalanceMutation struct {
	ID            int       `gorm:"column:id"`
	CustomerID    int       `gorm:"column:customer_id"`
	Amount        float64   `gorm:"column:amount"`
	BalanceAfter  float64   `gorm:"column:balance_after"`
	Type          string    `gorm:"column:type"`
	ReferenceType *string   `gorm:"column:reference_type"`
	ReferenceID   *int      `gorm:"column:reference_id"`
	Notes         *string   `gorm:"column:notes"`
	UserID        int       `gorm:"column:user_id"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}
