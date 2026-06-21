package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
		Search    string `json:"search" validate:"max=100"`
		IsActive  *bool  `json:"is_active"`
		SortBy    string `json:"sort_by"`
		SortOrder string `json:"sort_order"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		Name        string  `json:"name" validate:"required,min=2,max=100"`
		Phone       string  `json:"phone" validate:"max=20"`
		Address     string  `json:"address" validate:"max=255"`
		CreditLimit float64 `json:"credit_limit" validate:"min=0"`
		Notes       string  `json:"notes" validate:"max=500"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID          int     `json:"-"`
		Name        string  `json:"name" validate:"required,min=2,max=100"`
		Phone       string  `json:"phone" validate:"max=20"`
		Address     string  `json:"address" validate:"max=255"`
		CreditLimit float64 `json:"credit_limit" validate:"min=0"`
		Notes       string  `json:"notes" validate:"max=500"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	CustomerResponse struct {
		ID           int       `json:"id"`
		CustomerCode string    `json:"customer_code"`
		Name         string    `json:"name"`
		Phone        string    `json:"phone"`
		Address      string    `json:"address"`
		CreditLimit  float64   `json:"credit_limit"`
		IsActive     bool      `json:"is_active"`
		CreatedAt    time.Time `json:"created_at"`
	}

	CustomerDetailResponse struct {
		ID           int       `json:"id"`
		CustomerCode string    `json:"customer_code"`
		Name         string    `json:"name"`
		Phone        string    `json:"phone"`
		Address      string    `json:"address"`
		CreditLimit  float64   `json:"credit_limit"`
		Notes        string    `json:"notes"`
		IsActive     bool      `json:"is_active"`
		CreatedAt    time.Time `json:"created_at"`
	}

	CustomerActiveItem struct {
		ID           int     `json:"id"`
		Name         string  `json:"name"`
		CustomerCode string  `json:"customer_code"`
		CreditLimit  float64 `json:"credit_limit"`
	}
)
