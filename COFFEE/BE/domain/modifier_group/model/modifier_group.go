package model

// ModifierGroup: kelompok tambahan (mis. Level gula, Topping), keputusan 54 dan 87.
type ModifierGroup struct {
	ID            int    `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	SelectionType string `gorm:"column:selection_type"`
	IsRequired    bool   `gorm:"column:is_required"`
	IsActive      bool   `gorm:"column:is_active"`
	OptionCount   int    `gorm:"column:option_count"`
	MenuCount     int    `gorm:"column:menu_count"`
}

// ModifierOption: pilihan dalam kelompok. Harga 0 = gratis.
type ModifierOption struct {
	ID        int     `gorm:"column:id"`
	GroupID   int     `gorm:"column:group_id"`
	Name      string  `gorm:"column:name"`
	Price     float64 `gorm:"column:price"`
	IsDefault bool    `gorm:"column:is_default"`
	IsActive  bool    `gorm:"column:is_active"`
	SortOrder int     `gorm:"column:sort_order"`
}
