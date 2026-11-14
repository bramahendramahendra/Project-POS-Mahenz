package model

// MenuCategory: kategori menu jualan (mis. Kopi, Non-Kopi, Makanan).
type MenuCategory struct {
	ID        int    `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	SortOrder int    `gorm:"column:sort_order"`
	IsActive  bool   `gorm:"column:is_active"`
	ItemCount int    `gorm:"column:item_count"`
}
