package model

import "time"

type RouteRegistry struct {
	ID        int        `gorm:"column:id"`
	Path      string     `gorm:"column:path"`
	Label     string     `gorm:"column:label"`
	IsActive  bool       `gorm:"column:is_active"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (RouteRegistry) TableName() string {
	return "route_registry"
}
