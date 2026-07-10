package dto

type (
	// REQUEST
	CreateRequest struct {
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

	PackageRequest struct {
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		PackageName   string  `json:"package_name"`
		ConversionQty float64 `json:"conversion_qty" validate:"required,min=0.001"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"min=0"`
		IsDefault     bool    `json:"is_default"`
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
