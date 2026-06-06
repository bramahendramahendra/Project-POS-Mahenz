package dto

import "time"

type (
	GetCategoryByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CategoryListRequest struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Search string `json:"search" validate:"max=100"`
	}

	CategoryOptionResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
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
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=500"`
		Code        string `json:"code"`
	}

	UpdateCategoryUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateCategoryRequest struct {
		ID          int    `json:"-"`
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description" validate:"max=500"`
	}

	DeleteCategoryRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusCategoryRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}
)
