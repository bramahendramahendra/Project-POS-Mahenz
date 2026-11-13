package model

import "time"

// App: aplikasi yang bisa dibuka dari Portal.
type App struct {
	ID          int    `gorm:"column:id"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Path        string `gorm:"column:path"`
	Icon        string `gorm:"column:icon"`
	SortOrder   int    `gorm:"column:sort_order"`
	IsActive    bool   `gorm:"column:is_active"`
}

// AppAdmin: data aplikasi untuk keperluan admin (tanpa membuka hash kunci).
type AppAdmin struct {
	ID          int       `gorm:"column:id"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Path        string    `gorm:"column:path"`
	Icon        string    `gorm:"column:icon"`
	SortOrder   int       `gorm:"column:sort_order"`
	IsActive    bool      `gorm:"column:is_active"`
	HasKey      bool      `gorm:"column:has_key"`
	UserCount   int       `gorm:"column:user_count"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

// UserAccess: status akses seorang user terhadap satu aplikasi.
type UserAccess struct {
	AppID    int    `gorm:"column:app_id"`
	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	IsActive bool   `gorm:"column:is_active"`
	Granted  bool   `gorm:"column:granted"`
}
