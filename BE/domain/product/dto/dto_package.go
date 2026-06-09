package dto

type (
	// REQUEST
	PackageByProductRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	SavePackageRequest struct {
		ProductID int              `json:"-"`
		Packages  []PackageRequest `json:"packages" validate:"required,min=1,dive"`
	}

	DeletePackageRequest struct {
		ID        int `uri:"id" validate:"required,min=1"`
		PackageID int `uri:"package_id" validate:"required,min=1"`
	}

	// RESPONSE
	PackageResponse struct {
		ID            int     `json:"id"`
		ProductID     int     `json:"product_id"`
		UnitID        int     `json:"unit_id"`
		UnitName      string  `json:"unit_name"`
		Abbreviation  string  `json:"abbreviation"`
		PackageName   string  `json:"package_name"`
		ConversionQty float64 `json:"conversion_qty"`
		PurchasePrice float64 `json:"purchase_price"`
		SellingPrice  float64 `json:"selling_price"`
		IsDefault     bool    `json:"is_default"`
	}
)
