package model

import "time"

type Menu struct {
	ID         int        `gorm:"column:id"`
	ParentID   *int       `gorm:"column:parent_id"`
	KeyName    string     `gorm:"column:key_name"`
	Label      string     `gorm:"column:label"`
	Icon       *string    `gorm:"column:icon"`
	Path       *string    `gorm:"column:path"`
	OrderIndex int        `gorm:"column:order_index"`
	IsActive   bool       `gorm:"column:is_active"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
}
