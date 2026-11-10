package model

import "time"

type Role struct {
	ID          int        `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	DisplayName string     `gorm:"column:display_name"`
	Description *string    `gorm:"column:description"`
	IsSystem    bool       `gorm:"column:is_system"`
	IsActive    bool       `gorm:"column:is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}
