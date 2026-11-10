package dto

import "time"

type (
	// REQUEST
	TopupRequest struct {
		CustomerID int     `json:"-"`
		UserID     int     `json:"-"`
		Amount     float64 `json:"amount" validate:"required,gt=0"`
		Notes      string  `json:"notes" validate:"max=500"`
	}

	RefundRequest struct {
		CustomerID int     `json:"-"`
		UserID     int     `json:"-"`
		Amount     float64 `json:"amount" validate:"required,gt=0"`
		Notes      string  `json:"notes" validate:"max=500"`
	}

	HistoryRequest struct {
		CustomerID int `json:"-"`
		Page       int `json:"page" validate:"min=1"`
		Limit      int `json:"limit" validate:"min=1"`
	}

	SaveToBalanceRequest struct {
		TransactionID int     `json:"transaction_id" validate:"required,gt=0"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		UserID        int     `json:"-"`
	}

	// RESPONSE
	MutationResponse struct {
		ID            int       `json:"id"`
		CustomerID    int       `json:"customer_id"`
		Amount        float64   `json:"amount"`
		BalanceAfter  float64   `json:"balance_after"`
		Type          string    `json:"type"`
		ReferenceType *string   `json:"reference_type"`
		ReferenceID   *int      `json:"reference_id"`
		Notes         *string   `json:"notes"`
		UserName      string    `json:"user_name"`
		CreatedAt     time.Time `json:"created_at"`
	}

	BalanceResponse struct {
		CustomerID   int     `json:"customer_id"`
		Balance      float64 `json:"balance"`
		BalanceAfter float64 `json:"balance_after"`
	}
)
