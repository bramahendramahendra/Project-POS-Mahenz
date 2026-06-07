package model

import "time"

type (
	Unit struct {
		ID           int       `gorm:"column:id"`
		Name         string    `gorm:"column:name"`
		Abbreviation string    `gorm:"column:abbreviation"`
		IsActive     bool      `gorm:"column:is_active"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at"`
	}

	UnitOption struct {
		ID           int    `gorm:"column:id"`
		Name         string `gorm:"column:name"`
		Abbreviation string `gorm:"column:abbreviation"`
	}
)
