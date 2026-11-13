package model

import "time"

// Code: satu baris kode SSO beserta data pendukung (hasil join apps dan users).
type Code struct {
	ID         int        `gorm:"column:id"`
	UserID     int        `gorm:"column:user_id"`
	AppID      int        `gorm:"column:app_id"`
	ExpiresAt  time.Time  `gorm:"column:expires_at"`
	UsedAt     *time.Time `gorm:"column:used_at"`
	AppCode    string     `gorm:"column:app_code"`
	AppKeyHash string     `gorm:"column:app_key_hash"`
	AppActive  bool       `gorm:"column:app_active"`
	Username   string     `gorm:"column:username"`
	FullName   string     `gorm:"column:full_name"`
	UserActive bool       `gorm:"column:user_active"`
}

// LaunchApp: aplikasi tujuan yang sudah dipastikan boleh dibuka oleh user.
type LaunchApp struct {
	ID   int    `gorm:"column:id"`
	Code string `gorm:"column:code"`
	Path string `gorm:"column:path"`
}
