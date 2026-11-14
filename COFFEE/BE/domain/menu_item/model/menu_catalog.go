package model

// Baris mentah katalog kasir (tahap 4.7b).
type (
	CatalogCategory struct {
		ID   int    `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	CatalogItem struct {
		ID          int    `gorm:"column:id"`
		Name        string `gorm:"column:name"`
		CategoryID  int    `gorm:"column:category_id"`
		StockMode   string `gorm:"column:stock_mode"`
		ProductID   *int   `gorm:"column:product_id"`
		HasVariants bool   `gorm:"column:has_variants"`
		IsSoldOut   bool   `gorm:"column:is_sold_out"`
	}

	CatalogVariant struct {
		ID          int     `gorm:"column:id"`
		MenuItemID  int     `gorm:"column:menu_item_id"`
		SizeID      *int    `gorm:"column:size_id"`
		SizeName    string  `gorm:"column:size_name"`
		Temperature *string `gorm:"column:temperature"`
		Price       float64 `gorm:"column:price"`
	}

	CatalogGroup struct {
		MenuItemID    int    `gorm:"column:menu_item_id"`
		ID            int    `gorm:"column:id"`
		Name          string `gorm:"column:name"`
		SelectionType string `gorm:"column:selection_type"`
		IsRequired    bool   `gorm:"column:is_required"`
	}

	CatalogOption struct {
		ID        int     `gorm:"column:id"`
		GroupID   int     `gorm:"column:group_id"`
		Name      string  `gorm:"column:name"`
		Price     float64 `gorm:"column:price"`
		IsDefault bool    `gorm:"column:is_default"`
	}
)
