package model

import "time"

type User struct {
	ID        int       `gorm:"column:id"`
	Username  string    `gorm:"column:username"`
	FullName  string    `gorm:"column:full_name"`
	RoleID    int       `gorm:"column:role_id"`
	RoleName  string    `gorm:"column:role_name"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
