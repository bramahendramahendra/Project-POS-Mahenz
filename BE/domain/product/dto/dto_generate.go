package dto

type (
	// REQUEST
	GenerateSkuRequest struct {
		CategoryID int `json:"category_id" validate:"required,min=1"`
	}

	// RESPONSE
	GenerateBarcodeResponse struct {
		Barcode string `json:"barcode"`
	}

	GenerateSkuResponse struct {
		SKU string `json:"sku"`
	}
)
