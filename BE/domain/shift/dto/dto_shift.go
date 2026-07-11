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
		Name      string `json:"name" validate:"required,min=2,max=100,notblank"`
		StartTime string `json:"start_time" validate:"required,timeformat"`
		EndTime   string `json:"end_time" validate:"required,timeformat"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID        int    `json:"-"`
		Name      string `json:"name" validate:"required,min=2,max=100,notblank"`
		StartTime string `json:"start_time" validate:"required,timeformat"`
		EndTime   string `json:"end_time" validate:"required,timeformat"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	ShiftResponse struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		StartTime string    `json:"start_time"`
		EndTime   string    `json:"end_time"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
	}

	ShiftActiveResponse struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	ShiftSummaryResponse struct {
		ShiftID           int     `json:"shift_id"`
		ShiftName         string  `json:"shift_name"`
		TotalTransactions int     `json:"total_transactions"`
		TotalSales        float64 `json:"total_sales"`
		TotalCash         float64 `json:"total_cash"`
		TotalNonCash      float64 `json:"total_non_cash"`
	}
)
