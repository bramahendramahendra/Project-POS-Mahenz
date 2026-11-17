package model

type SupplierReturnRow struct {
	ID                int                  `gorm:"column:id"`
	ReturnCode        string               `gorm:"column:return_code"`
	PurchaseID        int                  `gorm:"column:purchase_id"`
	SupplierID        *int                 `gorm:"column:supplier_id"`
	SupplierName      string               `gorm:"column:supplier_name"`
	ReturnDate        string               `gorm:"column:return_date"`
	TotalReturnAmount float64              `gorm:"column:total_return_amount"`
	Reason            string               `gorm:"column:reason"`
	Status            string               `gorm:"column:status"`
	UserName          string               `gorm:"column:user_name"`
	Notes             string               `gorm:"column:notes"`
	Items             []SupplierReturnItem `gorm:"-"`
}

type SupplierReturnPurchaseRef struct {
	PurchaseID        int     `gorm:"column:purchase_id"`
	TotalReturnAmount float64 `gorm:"column:total_return_amount"`
}

type SupplierReturnItem struct {
	ID             int     `gorm:"column:id"`
	ReturnID       int     `gorm:"column:return_id"`
	PurchaseItemID int     `gorm:"column:purchase_item_id"`
	ProductID      int     `gorm:"column:product_id"`
	ProductName    string  `gorm:"column:product_name"`
	PackageID      *int    `gorm:"column:package_id"`
	Quantity       float64 `gorm:"column:quantity"`
	Unit           string  `gorm:"column:unit"`
	PurchasePrice  float64 `gorm:"column:purchase_price"`
	Subtotal       float64 `gorm:"column:subtotal"`
}
