package model


type TransactionItem struct {
	ID            int     `gorm:"column:id"`
	TransactionID int     `gorm:"column:transaction_id"`
	ProductID     int     `gorm:"column:product_id"`
	ProductName   string  `gorm:"column:product_name"`
	Quantity      float64 `gorm:"column:quantity"`
	Unit          string  `gorm:"column:unit"`
	Price         float64 `gorm:"column:price"`
	PurchasePrice float64 `gorm:"column:purchase_price"`
	Subtotal      float64 `gorm:"column:subtotal"`
	DiscountItem  float64 `gorm:"column:discount_item"`
	ConversionQty float64 `gorm:"column:conversion_qty"`
	UnitID        *int    `gorm:"column:unit_id"`
}

