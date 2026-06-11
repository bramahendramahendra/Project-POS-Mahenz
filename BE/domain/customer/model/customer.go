package model

import "time"

type Customer struct {
	ID           int       `gorm:"column:id"`
	CustomerCode string    `gorm:"column:customer_code"`
	Name         string    `gorm:"column:name"`
	Phone        string    `gorm:"column:phone"`
	Address      string    `gorm:"column:address"`
	CreditLimit  float64   `gorm:"column:credit_limit"`
	Notes        string    `gorm:"column:notes"`
	IsActive     bool      `gorm:"column:is_active"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}
