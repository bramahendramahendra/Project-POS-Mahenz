package dto

import "time"

type (
	GetSupplierByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	SupplierListRequest struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Search   string `json:"search" validate:"max=100"`
		IsActive *bool  `json:"is_active"`
	}

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

	SupplierOptionResponse struct {
		ID           int    `json:"id"`
		SupplierCode string `json:"supplier_code"`
		Name         string `json:"name"`
	}

	SupplierPurchaseItem struct {
		ID              int     `json:"id"`
		PurchaseCode    string  `json:"purchase_code"`
		PurchaseDate    string  `json:"purchase_date"`
		TotalAmount     float64 `json:"total_amount"`
		PaymentStatus   string  `json:"payment_status"`
		RemainingAmount float64 `json:"remaining_amount"`
	}

	SupplierReturnHistoryItem struct {
		ID          int     `json:"id"`
		ReturnCode  string  `json:"return_code"`
		ReturnDate  string  `json:"return_date"`
		TotalReturn float64 `json:"total_return"`
		Reason      string  `json:"reason"`
		Status      string  `json:"status"`
	}

	SupplierDetailResponse struct {
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
		PurchaseHistory []SupplierPurchaseItem      `json:"purchase_history"`
		ReturnHistory   []SupplierReturnHistoryItem `json:"return_history"`
	}

	CreateSupplierRequest struct {
		Name          string `json:"name" validate:"required,min=2,max=100"`
		Address       string `json:"address" validate:"omitempty,max=255"`
		Phone         string `json:"phone" validate:"omitempty,max=20"`
		Email         string `json:"email" validate:"omitempty,max=100,email"`
		ContactPerson string `json:"contact_person" validate:"omitempty,max=100"`
		Notes         string `json:"notes" validate:"omitempty,max=500"`
	}

	UpdateSupplierUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateSupplierRequest struct {
		ID            int    `json:"-"`
		Name          string `json:"name" validate:"required,min=2,max=100"`
		Address       string `json:"address" validate:"omitempty,max=255"`
		Phone         string `json:"phone" validate:"omitempty,max=20"`
		Email         string `json:"email" validate:"omitempty,max=100,email"`
		ContactPerson string `json:"contact_person" validate:"omitempty,max=100"`
		Notes         string `json:"notes" validate:"omitempty,max=500"`
	}

	DeleteSupplierRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ToggleStatusSupplierRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}
)
