package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page      int    `json:"page" validate:"required,min=1"`
		Limit     int    `json:"limit" validate:"required,min=1"`
		Search    string `json:"search" validate:"max=100"`
		IsActive  *bool  `json:"is_active"`
		SortBy    string `json:"sort_by"`
		SortOrder string `json:"sort_order"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		Name         string `json:"name" validate:"required,min=2,max=100"`
		Abbreviation string `json:"abbreviation" validate:"required,min=2,max=20"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID           int    `json:"-"`
		Name         string `json:"name" validate:"required,min=2,max=100"`
		Abbreviation string `json:"abbreviation" validate:"required,min=2,max=20"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	UnitResponse struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		Abbreviation string    `json:"abbreviation"`
		IsActive     bool      `json:"is_active"`
		CreatedAt    time.Time `json:"created_at"`
	}

	GetOptionResponse struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
	}
)
