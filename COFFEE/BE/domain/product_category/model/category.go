package model

import "time"

type (
	Category struct {
		ID                 int       `gorm:"column:id"`
		Name               string    `gorm:"column:name"`
		Code               string    `gorm:"column:code"`
		Description        string    `gorm:"column:description"`
		IsActive           bool      `gorm:"column:is_active"`
		CreatedAt          time.Time `gorm:"column:created_at"`
		UpdatedAt          time.Time `gorm:"column:updated_at"`
		ProductCount       int       `gorm:"column:product_count"`
		ActiveProductCount int       `gorm:"column:active_product_count"`
	}

	CategoryOption struct {
		ID   int    `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
)
