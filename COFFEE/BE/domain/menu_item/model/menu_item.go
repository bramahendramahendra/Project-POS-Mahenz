package model

// MenuItem: menu jualan. StockMode: recipe (racikan) atau direct (langsung / barang jadi).
type MenuItem struct {
	ID           int     `gorm:"column:id"`
	CategoryID   int     `gorm:"column:category_id"`
	CategoryName string  `gorm:"column:category_name"`
	Name         string  `gorm:"column:name"`
	StockMode    string  `gorm:"column:stock_mode"`
	ProductID    *int    `gorm:"column:product_id"`
	ProductName  string  `gorm:"column:product_name"`
	HasVariants  bool    `gorm:"column:has_variants"`
	IsActive     bool    `gorm:"column:is_active"`
	IsSoldOut    bool    `gorm:"column:is_sold_out"`
	PriceMin     float64 `gorm:"column:price_min"`
	PriceMax     float64 `gorm:"column:price_max"`
	VariantCount int     `gorm:"column:variant_count"`
}

// Variant: satu kombinasi ukuran x suhu dan harganya. Menu tanpa varian punya satu baris
// dengan SizeID dan Temperature kosong.
type Variant struct {
	ID          int     `gorm:"column:id"`
	MenuItemID  int     `gorm:"column:menu_item_id"`
	SizeID      *int    `gorm:"column:size_id"`
	SizeName    string  `gorm:"column:size_name"`
	SizeOrder   int     `gorm:"column:size_order"`
	Temperature *string `gorm:"column:temperature"`
	Price       float64 `gorm:"column:price"`
	IsActive    bool    `gorm:"column:is_active"`
}

// RecipeLine: satu baris resep (bahan + takaran teks), keputusan 86.
type RecipeLine struct {
	ID             int    `gorm:"column:id"`
	VariantID      int    `gorm:"column:variant_id"`
	IngredientID   int    `gorm:"column:ingredient_id"`
	IngredientName string `gorm:"column:ingredient_name"`
	Amount         string `gorm:"column:amount"`
	SortOrder      int    `gorm:"column:sort_order"`
}

// ItemModifierGroup: kelompok tambahan yang dipasang ke menu.
type ItemModifierGroup struct {
	GroupID   int    `gorm:"column:group_id"`
	GroupName string `gorm:"column:group_name"`
	SortOrder int    `gorm:"column:sort_order"`
}
