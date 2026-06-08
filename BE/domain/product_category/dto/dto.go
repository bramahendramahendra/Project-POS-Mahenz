package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Search string `json:"search" validate:"max=100"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=500"`
		Code        string `json:"code"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID          int    `json:"-"`
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=500"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	CategoryResponse struct {
		ID                 int       `json:"id"`
		Name               string    `json:"name"`
		Code               string    `json:"code"`
		Description        string    `json:"description"`
		IsActive           bool      `json:"is_active"`
		ProductCount       int       `json:"product_count"`
		ActiveProductCount int       `json:"active_product_count"`
		CreatedAt          time.Time `json:"created_at"`
	}

	GetOptionResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
