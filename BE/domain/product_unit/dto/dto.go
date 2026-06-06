package dto

import "time"

type (
	GetUnitByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UnitListRequest struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Search string `json:"search"`
	}

	UnitResponse struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		Abbreviation string    `json:"abbreviation"`
		IsActive     bool      `json:"is_active"`
		CreatedAt    time.Time `json:"created_at"`
	}

	UnitActiveResponse struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
	}

	CreateUnitRequest struct {
		Name         string `json:"name" validate:"required,min=1"`
		Abbreviation string `json:"abbreviation" validate:"required,min=1"`
	}

	UpdateUnitUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateUnitRequest struct {
		ID           int    `json:"-"`
		Name         string `json:"name" validate:"required,min=1"`
		Abbreviation string `json:"abbreviation" validate:"required,min=1"`
	}

	DeleteUnitRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusUnitRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}
)
