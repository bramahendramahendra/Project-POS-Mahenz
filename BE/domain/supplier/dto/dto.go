package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Search   string `json:"search" validate:"max=100"`
		IsActive *bool  `json:"is_active"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		Name          string `json:"name" validate:"required,min=2,max=100"`
		Address       string `json:"address" validate:"omitempty,max=255"`
		Phone         string `json:"phone" validate:"omitempty,max=20"`
		Email         string `json:"email" validate:"omitempty,max=100,email"`
		ContactPerson string `json:"contact_person" validate:"omitempty,max=100"`
		Notes         string `json:"notes" validate:"omitempty,max=500"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID            int    `json:"-"`
		Name          string `json:"name" validate:"required,min=2,max=100"`
		Address       string `json:"address" validate:"omitempty,max=255"`
		Phone         string `json:"phone" validate:"omitempty,max=20"`
		Email         string `json:"email" validate:"omitempty,max=100,email"`
		ContactPerson string `json:"contact_person" validate:"omitempty,max=100"`
		Notes         string `json:"notes" validate:"omitempty,max=500"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	SupplierResponse struct {
		ID            int       `json:"id"`
		SupplierCode  string    `json:"supplier_code"`
		Name          string    `json:"name"`
		Address       string    `json:"address"`
		Phone         string    `json:"phone"`
		Email         string    `json:"email"`
		ContactPerson string    `json:"contact_person"`
		Notes         string    `json:"notes"`
		IsActive      bool      `json:"is_active"`
		CreatedAt     time.Time `json:"created_at"`
	}

	GetOptionResponse struct {
		ID           int    `json:"id"`
		SupplierCode string `json:"supplier_code"`
		Name         string `json:"name"`
	}

	GetDetailPurchaseResponse struct {
		ID              int     `json:"id"`
		PurchaseCode    string  `json:"purchase_code"`
		PurchaseDate    string  `json:"purchase_date"`
		TotalAmount     float64 `json:"total_amount"`
		PaymentStatus   string  `json:"payment_status"`
		RemainingAmount float64 `json:"remaining_amount"`
	}

	GetDetailReturnResponse struct {
		ID          int     `json:"id"`
		ReturnCode  string  `json:"return_code"`
		ReturnDate  string  `json:"return_date"`
		TotalReturn float64 `json:"total_return"`
		Reason      string  `json:"reason"`
		Status      string  `json:"status"`
	}

	GetDetailResponse struct {
		ID              int                         `json:"id"`
		SupplierCode    string                      `json:"supplier_code"`
		Name            string                      `json:"name"`
		Address         string                      `json:"address"`
		Phone           string                      `json:"phone"`
		Email           string                      `json:"email"`
		ContactPerson   string                      `json:"contact_person"`
		Notes           string                      `json:"notes"`
		IsActive        bool                        `json:"is_active"`
		TotalPurchases  int                         `json:"total_purchases"`
		TotalAmount     float64                     `json:"total_amount"`
		TotalDebt       float64                     `json:"total_debt"`
		TotalReturn     float64                     `json:"total_return"`
		PurchaseHistory []GetDetailPurchaseResponse `json:"purchase_history"`
		ReturnHistory   []GetDetailReturnResponse   `json:"return_history"`
	}
)
