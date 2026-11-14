package model

// MenuSize: ukuran jualan (mis. Regular, Large), dikelola admin (keputusan 85).
type MenuSize struct {
	ID         int    `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	SortOrder  int    `gorm:"column:sort_order"`
	IsActive   bool   `gorm:"column:is_active"`
	UsageCount int    `gorm:"column:usage_count"`
}
