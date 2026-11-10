package dto

import "time"

type (
	// ============ REQUEST ============

	// ListRequest — daftar baris penjualan yang modalnya perlu ditinjau manual.
	ListRequest struct {
		Page   int    `json:"page" validate:"required,min=1"`
		Limit  int    `json:"limit" validate:"required,min=1"`
		Search string `json:"search" validate:"max=100"`
		Status string `json:"status" validate:"omitempty,oneof=pending resolved skipped all"` // default: pending
	}

	// ResolveRequest — admin menetapkan modal (HPP) yang benar untuk 1 baris.
	ResolveRequest struct {
		ReviewID      int     `json:"review_id" validate:"required,min=1"`
		CorrectCost   float64 `json:"correct_cost" validate:"min=0"` // modal per satuan jual yang benar
		Note          string  `json:"note" validate:"max=500"`
		UserID        int     `json:"-"` // diisi server dari context
	}

	// SkipRequest — tandai baris sudah ditinjau tanpa mengubah modal
	// (mis. setelah dicek ternyata modal tersimpan sudah benar).
	SkipRequest struct {
		ReviewID int    `json:"review_id" validate:"required,min=1"`
		Note     string `json:"note" validate:"max=500"`
		UserID   int    `json:"-"`
	}

	// ============ RESPONSE ============

	// ListItem — 1 baris di tabel daftar rekonsiliasi modal.
	ListItem struct {
		ReviewID         int        `json:"review_id"`
		TransactionItemID int       `json:"transaction_item_id"`
		TransactionID    int        `json:"transaction_id"`
		TransactionCode  string     `json:"transaction_code"`
		TransactionDate  time.Time  `json:"transaction_date"`
		ProductID        int        `json:"product_id"`
		ProductName      string     `json:"product_name"`
		Unit             string     `json:"unit"`
		Quantity         float64    `json:"quantity"`
		SellPrice        float64    `json:"sell_price"`         // harga jual per satuan
		OldPurchasePrice float64    `json:"old_purchase_price"` // modal tersimpan saat ini
		SuggestedPrice   float64    `json:"suggested_price"`    // hasil hitung otomatis (meragukan)
		Reason           string     `json:"reason"`
		Status           string     `json:"status"`
		ResolvedPrice    *float64   `json:"resolved_price,omitempty"`
		Note             string     `json:"note"`
		ResolvedBy       string     `json:"resolved_by,omitempty"`
		ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	}

	// SummaryResponse — kartu ringkasan di atas tabel.
	SummaryResponse struct {
		Pending  int `json:"pending"`
		Resolved int `json:"resolved"`
		Skipped  int `json:"skipped"`
		Total    int `json:"total"`
	}
)
