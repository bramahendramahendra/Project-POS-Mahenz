package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page          int    `json:"page" validate:"required,min=1"`
		Limit         int    `json:"limit" validate:"required,min=1"`
		ProductID     *int   `json:"product_id"`
		MutationType  string `json:"mutation_type" validate:"max=50"`
		ReferenceType string `json:"reference_type" validate:"max=50"`
		DateFrom      string `json:"date_from"`
		DateTo        string `json:"date_to"`
	}

	GetByProductRequest struct {
		ProductID int `uri:"product_id" validate:"required,min=1"`
	}

	// RESPONSE
	StockMutationResponse struct {
		ID            int       `json:"id"`
		ProductID     int       `json:"product_id"`
		ProductName   string    `json:"product_name"`
		MutationType  string    `json:"mutation_type"`
		Quantity      float64   `json:"quantity"`
		StockBefore   float64   `json:"stock_before"`
		StockAfter    float64   `json:"stock_after"`
		ReferenceType string    `json:"reference_type"`
		ReferenceID   int       `json:"reference_id"`
		Notes         string    `json:"notes"`
		UserName      string    `json:"user_name"`
		CreatedAt     time.Time `json:"created_at"`
	}

	StockMutationByProductResponse struct {
		ID            int       `json:"id"`
		MutationType  string    `json:"mutation_type"`
		Quantity      float64   `json:"quantity"`
		StockBefore   float64   `json:"stock_before"`
		StockAfter    float64   `json:"stock_after"`
		ReferenceType string    `json:"reference_type"`
		ReferenceID   int       `json:"reference_id"`
		Notes         string    `json:"notes"`
		UserName      string    `json:"user_name"`
		CreatedAt     time.Time `json:"created_at"`
	}
)
