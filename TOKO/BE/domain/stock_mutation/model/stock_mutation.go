package model

import "time"

type StockMutation struct {
	ID            int       `gorm:"column:id"`
	ProductID     int       `gorm:"column:product_id"`
	ProductName   string    `gorm:"column:product_name"`
	MutationType  string    `gorm:"column:mutation_type"`
	Quantity      float64   `gorm:"column:quantity"`
	StockBefore   float64   `gorm:"column:stock_before"`
	StockAfter    float64   `gorm:"column:stock_after"`
	ReferenceType string    `gorm:"column:reference_type"`
	ReferenceID   int       `gorm:"column:reference_id"`
	Notes         string    `gorm:"column:notes"`
	UserID        int       `gorm:"column:user_id"`
	UserName      string    `gorm:"column:user_name"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}
