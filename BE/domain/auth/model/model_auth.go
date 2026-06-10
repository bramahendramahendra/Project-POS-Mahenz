package model

import "time"

type User struct {
	ID        int       `gorm:"column:id"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	FullName  string    `gorm:"column:full_name"`
	RoleID    int       `gorm:"column:role_id"`
	RoleName  string    `gorm:"column:role_name"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Session struct {
	ID           int       `gorm:"column:id"`
	UserID       int       `gorm:"column:user_id"`
	UserRole     string    `gorm:"column:user_role"`
	Token        string    `gorm:"column:token"`
	RefreshToken string    `gorm:"column:refresh_token"`
	DeviceInfo   string    `gorm:"column:device_info"`
	IPAddress    string    `gorm:"column:ip_address"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
}
