package dto

// Katalog kasir (tahap 4.7b): menu aktif dari kategori aktif beserta kombinasi yang dijual dan
// kelompok tambahan aktif, dalam satu panggilan agar halaman kasir cepat di tablet.
type (
	CatalogResponse struct {
		Categories []CatalogCategory `json:"categories"`
		Items      []CatalogItem     `json:"items"`
	}

	CatalogCategory struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	CatalogItem struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		CategoryID  int    `json:"category_id"`
		StockMode   string `json:"stock_mode"`
		HasVariants bool   `json:"has_variants"`
		IsSoldOut   bool   `json:"is_sold_out"`
		// Stock: sisa stok produk barang jadi (satuan dasar) yang bisa dijual; nil untuk racikan.
		Stock          *float64               `json:"stock"`
		Variants       []CatalogVariant       `json:"variants"`
		ModifierGroups []CatalogModifierGroup `json:"modifier_groups"`
	}

	CatalogVariant struct {
		ID          int     `json:"id"`
		SizeID      *int    `json:"size_id"`
		SizeName    string  `json:"size_name"`
		Temperature *string `json:"temperature"`
		Price       float64 `json:"price"`
	}

	CatalogModifierGroup struct {
		ID            int                     `json:"id"`
		Name          string                  `json:"name"`
		SelectionType string                  `json:"selection_type"`
		IsRequired    bool                    `json:"is_required"`
		Options       []CatalogModifierOption `json:"options"`
	}

	CatalogModifierOption struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		IsDefault bool    `json:"is_default"`
	}
)
