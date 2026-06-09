package dto

type (
	// REQUEST
	PriceByProductRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	PriceRequest struct {
		TierName string  `json:"tier_name" validate:"required"`
		MinQty   float64 `json:"min_qty" validate:"min=0"`
		Price    float64 `json:"price" validate:"required,min=0"`
	}

	SavePriceRequest struct {
		ProductID int            `json:"-"`
		Prices    []PriceRequest `json:"prices" validate:"dive"`
	}

	// RESPONSE
	PriceResponse struct {
		ID        int     `json:"id"`
		ProductID int     `json:"product_id"`
		TierName  string  `json:"tier_name"`
		MinQty    float64 `json:"min_qty"`
		Price     float64 `json:"price"`
	}
)
