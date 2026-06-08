package model_product

import "time"

type Product struct {
	ID               int       `gorm:"column:id"`
	Barcode          string    `gorm:"column:barcode"`
	SKU              string    `gorm:"column:sku"`
	Name             string    `gorm:"column:name"`
	CategoryID       *int      `gorm:"column:category_id"`
	CategoryName     string    `gorm:"column:category_name"`
	PurchasePrice    float64   `gorm:"column:purchase_price"`
	SellingPrice     float64   `gorm:"column:selling_price"`
	Stock            float64   `gorm:"column:stock"`
	MinStock         float64   `gorm:"column:min_stock"`
	UnitID           int       `gorm:"column:unit_id"`
	UnitName         string    `gorm:"column:unit_name"`
	UnitAbbreviation string    `gorm:"column:unit_abbreviation"`
	IsActive         bool      `gorm:"column:is_active"`
	ExtraPackages    int       `gorm:"column:extra_packages"`
	PriceTiersCount  int       `gorm:"column:price_tiers_count"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

type ProductOption struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type ProductSearchResult struct {
	ID           int     `gorm:"column:id"`
	Barcode      string  `gorm:"column:barcode"`
	Name         string  `gorm:"column:name"`
	SellingPrice float64 `gorm:"column:selling_price"`
	Stock        float64 `gorm:"column:stock"`
	UnitID       int     `gorm:"column:unit_id"`
	UnitName     string  `gorm:"column:unit_name"`
}

type LowStockProduct struct {
	ID       int     `gorm:"column:id"`
	Name     string  `gorm:"column:name"`
	Stock    float64 `gorm:"column:stock"`
	MinStock float64 `gorm:"column:min_stock"`
	UnitName string  `gorm:"column:unit_name"`
}
