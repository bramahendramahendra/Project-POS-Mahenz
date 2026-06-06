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
