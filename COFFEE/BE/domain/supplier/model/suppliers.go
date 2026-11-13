package model

import "time"

type Supplier struct {
	ID            int        `gorm:"column:id"`
	SupplierCode  string     `gorm:"column:supplier_code"`
	Name          string     `gorm:"column:name"`
	Address       string     `gorm:"column:address"`
	Phone         string     `gorm:"column:phone"`
	Email         string     `gorm:"column:email"`
	ContactPerson string     `gorm:"column:contact_person"`
	Notes         string     `gorm:"column:notes"`
	IsActive      bool       `gorm:"column:is_active"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

type SupplierOption struct {
	ID           int    `gorm:"column:id"`
	SupplierCode string `gorm:"column:supplier_code"`
	Name         string `gorm:"column:name"`
}

type SupplierPurchase struct {
	ID              int     `gorm:"column:id"`
	PurchaseCode    string  `gorm:"column:purchase_code"`
	PurchaseDate    string  `gorm:"column:purchase_date"`
	TotalAmount     float64 `gorm:"column:total_amount"`
	PaymentStatus   string  `gorm:"column:payment_status"`
	RemainingAmount float64 `gorm:"column:remaining_amount"`
}

type SupplierReturn struct {
	ID          int     `gorm:"column:id"`
	ReturnCode  string  `gorm:"column:return_code"`
	ReturnDate  string  `gorm:"column:return_date"`
	TotalReturn float64 `gorm:"column:total_return"`
	Reason      string  `gorm:"column:reason"`
	Status      string  `gorm:"column:status"`
}
