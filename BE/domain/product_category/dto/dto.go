package dto

import "time"

type (
	GetCategoryByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

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

	CreateCategoryRequest struct {
		Name        string `json:"name" validate:"required,min=2"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}

	UpdateCategoryUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateCategoryRequest struct {
		ID          int    `json:"-"`
		Name        string `json:"name" validate:"required,min=2"`
		Description string `json:"description"`
	}

	DeleteCategoryRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusCategoryRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}
)
