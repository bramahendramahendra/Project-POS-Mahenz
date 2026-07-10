package dto

type (
	// REQUEST
	GetAllRequest struct {
		Page       int    `json:"page" validate:"required,min=1"`
		Limit      int    `json:"limit" validate:"required,min=1"`
		Search     string `json:"search" validate:"max=100"`
		CategoryID *int   `json:"category_id"`
		IsActive   *bool  `json:"is_active"`
		LowStock   bool   `json:"low_stock"`
		SortBy     string `json:"sort_by"`
		SortOrder  string `json:"sort_order"`
	}

	SearchRequest struct {
		Q     string `json:"q" validate:"required"`
		Limit int    `json:"limit"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	GetByBarcodeRequest struct {
		Barcode string `uri:"barcode" validate:"required"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID            int     `json:"-"`
		Barcode       string  `json:"barcode" validate:"required,max=100"`
		SKU           string  `json:"sku" validate:"required,max=50"`
		Name          string  `json:"name" validate:"required,max=200"`
		CategoryID    *int    `json:"category_id" validate:"required"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"required,min=0"`
		Stock         float64 `json:"stock" validate:"min=0"`
		MinStock      float64 `json:"min_stock" validate:"min=0"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	ProductResponse struct {
		ID               int     `json:"id"`
		Barcode          string  `json:"barcode"`
		SKU              string  `json:"sku"`
		Name             string  `json:"name"`
		CategoryID       *int    `json:"category_id"`
		CategoryName     string  `json:"category_name"`
		PurchasePrice    float64 `json:"purchase_price"`
		SellingPrice     float64 `json:"selling_price"`
		Stock            float64 `json:"stock"`
		ReservedQty      float64 `json:"reserved_qty"`
		MinStock         float64 `json:"min_stock"`
		UnitID           int     `json:"unit_id"`
		UnitName         string  `json:"unit_name"`
		UnitAbbreviation string  `json:"unit_abbreviation"`
		IsActive         bool    `json:"is_active"`
		ExtraPackages    int     `json:"extra_packages"`
		PriceTiersCount  int     `json:"price_tiers_count"`
	}

	GetOptionResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	SearchResponse struct {
		ID           int     `json:"id"`
		Barcode      string  `json:"barcode"`
		Name         string  `json:"name"`
		SellingPrice float64 `json:"selling_price"`
		Stock        float64 `json:"stock"`
		MinStock     float64 `json:"min_stock"`
		UnitID       int     `json:"unit_id"`
		UnitName     string  `json:"unit_name"`
	}
)
