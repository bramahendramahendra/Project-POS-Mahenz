package model

type ProductPackage struct {
	ID            int     `gorm:"column:id"`
	ProductID     int     `gorm:"column:product_id"`
	UnitID        int     `gorm:"column:unit_id"`
	UnitName      string  `gorm:"column:unit_name"`
	Abbreviation  string  `gorm:"column:abbreviation"`
	PackageName   string  `gorm:"column:package_name"`
	ConversionQty float64 `gorm:"column:conversion_qty"`
	PurchasePrice float64 `gorm:"column:purchase_price"`
	SellingPrice  float64 `gorm:"column:selling_price"`
	IsDefault     bool    `gorm:"column:is_default"`
}
