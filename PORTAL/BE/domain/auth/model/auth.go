package model

import "time"

// User: profil lokal Portal, ditautkan ke akun BE User lewat AuthUserID.
// Password, role, dan status aktif dikelola BE User.
type User struct {
	ID         int    `gorm:"column:id"`
	AuthUserID int    `gorm:"column:auth_user_id"`
	Username   string `gorm:"column:username"`
	FullName   string `gorm:"column:full_name"`
}

type Session struct {
	ID                 int       `gorm:"column:id"`
	UserID             int       `gorm:"column:user_id"`      // id profil lokal
	AuthUserID         int       `gorm:"column:auth_user_id"` // hanya terisi saat dibaca (join users)
	UserRole           string    `gorm:"column:user_role"`
	MustChangePassword bool      `gorm:"column:must_change_password"`
	Token              string    `gorm:"column:token"`
	RefreshToken       string    `gorm:"column:refresh_token"`
	DeviceInfo         string    `gorm:"column:device_info"`
	IPAddress          string    `gorm:"column:ip_address"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	ExpiresAt          time.Time `gorm:"column:expires_at"`
}
