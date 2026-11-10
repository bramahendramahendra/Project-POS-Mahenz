package model

type ProductPrice struct {
	ID        int     `gorm:"column:id"`
	ProductID int     `gorm:"column:product_id"`
	TierName  string  `gorm:"column:tier_name"`
	MinQty    float64 `gorm:"column:min_qty"`
	Price     float64 `gorm:"column:price"`
}
