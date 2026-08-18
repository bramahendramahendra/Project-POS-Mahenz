package dto

type (
	// REQUEST
	CreateRequest struct {
		UserID        int                   `json:"-"`
		Barcode       string                `json:"barcode" validate:"required,max=100"`
		SKU           string                `json:"sku" validate:"required,max=50"`
		Name          string                `json:"name" validate:"required,max=200"`
		CategoryID    *int                  `json:"category_id" validate:"required"`
		PurchasePrice float64               `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64               `json:"selling_price" validate:"required,min=0"`
		Stock         float64               `json:"stock" validate:"min=0"`
		MinStock      float64               `json:"min_stock" validate:"min=0"`
		UnitID        int                   `json:"unit_id" validate:"required,min=1"`
		Packages      []PackageDraftRequest `json:"packages" validate:"dive"`
	}

	PackageDraftRequest struct {
		TempID        int     `json:"temp_id" validate:"required,min=1"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		PackageName   string  `json:"package_name" validate:"max=100"`
		RefTempID     int     `json:"ref_temp_id"`
		Qty           float64 `json:"qty" validate:"required,min=0.001"`
		RefQty        float64 `json:"ref_qty" validate:"required,min=0.001"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"min=0"`
	}

	// RESPONSE
	GetLowStockResponse struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Stock    float64 `json:"stock"`
		MinStock float64 `json:"min_stock"`
		UnitName string  `json:"unit_name"`
	}

	GetUnitInfoResponse struct {
		Name         string
		Abbreviation string
	}
)
