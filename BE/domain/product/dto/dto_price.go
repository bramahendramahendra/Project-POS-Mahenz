package dto

type (
	// REQUEST
	PriceByProductRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
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
