package model

import "time"

type ExpiryBatch struct {
	ID             int        `gorm:"column:id"`
	ProductID      int        `gorm:"column:product_id"`
	ProductName    string     `gorm:"column:product_name"`
	PurchaseItemID int        `gorm:"column:purchase_item_id"`
	Qty            float64    `gorm:"column:qty"`
	ExpiredDate    time.Time  `gorm:"column:expired_date"`
	Status         string     `gorm:"column:status"`
	ResolvedBy     *int       `gorm:"column:resolved_by"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at"`
	Notes          string     `gorm:"column:notes"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
}
