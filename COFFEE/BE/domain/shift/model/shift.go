package model

import "time"

type Shift struct {
	ID        int       `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	StartTime string    `gorm:"column:start_time"`
	EndTime   string    `gorm:"column:end_time"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
