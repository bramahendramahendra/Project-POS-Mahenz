package dto

import "time"

type (
	// ============ REQUEST ============

	// ListRequest — daftar produk rekonsiliasi stok (lama vs baru).
	ListRequest struct {
		Page       int    `json:"page" validate:"required,min=1"`
		Limit      int    `json:"limit" validate:"required,min=1"`
		Search     string `json:"search" validate:"max=100"`
		CategoryID *int   `json:"category_id"`
		OnlyReview bool   `json:"only_review"` // true = hanya produk needs_stock_review
	}

	// DetailRequest — detail 1 produk (breakdown + kartu stok + level satuan).
	DetailRequest struct {
		ProductID int `uri:"product_id" validate:"required,min=1"`
	}

	// AdjustLevel — 1 baris koreksi stok per level satuan.
	AdjustLevel struct {
		PackageID int     `json:"package_id" validate:"required,min=1"`
		NewStock  float64 `json:"new_stock" validate:"min=0"`
	}

	// AdjustRequest — koreksi manual stok per level satuan.
	AdjustRequest struct {
		ProductID int           `json:"product_id" validate:"required,min=1"`
		Levels    []AdjustLevel `json:"levels" validate:"required,min=1,dive"`
		Note      string        `json:"note" validate:"max=500"`
		UserID    int           `json:"-"` // diisi server dari context, bukan dari client
	}

	// MarkReviewedRequest — tandai selesai tanpa mengubah stok.
	MarkReviewedRequest struct {
		ProductID int    `json:"product_id" validate:"required,min=1"`
		Note      string `json:"note" validate:"max=500"`
		UserID    int    `json:"-"`
	}

	// ============ RESPONSE ============

	// ListItem — 1 baris di tabel daftar rekonsiliasi.
	ListItem struct {
		ProductID         int     `json:"product_id"`
		ProductCode       string  `json:"product_code"`
		ProductName       string  `json:"product_name"`
		BaseUnit          string  `json:"base_unit"`
		OldStock          float64 `json:"old_stock"`
		OldStockAvailable bool    `json:"old_stock_available"`
		NewStock          float64 `json:"new_stock"`
		Diff              float64 `json:"diff"`
		NeedsStockReview  bool    `json:"needs_stock_review"`
		StockReviewNote   string  `json:"stock_review_note"`
	}

	// SummaryResponse — kartu ringkasan di atas tabel.
	SummaryResponse struct {
		TotalProducts   int  `json:"total_products"`
		NeedsReview     int  `json:"needs_review"`
		Matched         int  `json:"matched"`
		BackupAvailable bool `json:"backup_available"`
	}

	// Breakdown — rincian "stok baru dari mana" (dalam satuan dasar/anchor).
	Breakdown struct {
		PurchaseIn       float64 `json:"purchase_in"`       // in
		PurchaseVoid     float64 `json:"purchase_void"`     // void_purchase
		SaleOut          float64 `json:"sale_out"`          // out
		SaleVoid         float64 `json:"sale_void"`         // void
		SupplierReturn   float64 `json:"supplier_return"`   // return
		Expired          float64 `json:"expired"`           // expired
		Adjustment       float64 `json:"adjustment"`        // adjustment (bertanda)
		ComputedNewStock float64 `json:"computed_new_stock"` // hasil rumus
	}

	// PackageLevel — 1 level satuan produk + stok saat ini + faktor ke dasar.
	PackageLevel struct {
		PackageID    int     `json:"package_id"`
		UnitName     string  `json:"unit_name"`
		PackageName  string  `json:"package_name"`
		IsDefault    bool    `json:"is_default"`
		FactorToBase float64 `json:"factor_to_base"`
		CurrentStock float64 `json:"current_stock"`
	}

	// KartuStokRow — 1 baris riwayat mutasi (kartu stok).
	KartuStokRow struct {
		ID            int       `json:"id"`
		MutationType  string    `json:"mutation_type"`
		Quantity      float64   `json:"quantity"`
		StockBefore   float64   `json:"stock_before"`
		StockAfter    float64   `json:"stock_after"`
		ReferenceType string    `json:"reference_type"`
		ReferenceID   int       `json:"reference_id"`
		Notes         string    `json:"notes"`
		UserName      string    `json:"user_name"`
		CreatedAt     time.Time `json:"created_at"`
	}

	// DetailResponse — payload lengkap layar detail.
	DetailResponse struct {
		ProductID         int            `json:"product_id"`
		ProductCode       string         `json:"product_code"`
		ProductName       string         `json:"product_name"`
		BaseUnit          string         `json:"base_unit"`
		OldStock          float64        `json:"old_stock"`
		OldStockAvailable bool           `json:"old_stock_available"`
		NewStock          float64        `json:"new_stock"`
		Diff              float64        `json:"diff"`
		NeedsStockReview  bool           `json:"needs_stock_review"`
		StockReviewNote   string         `json:"stock_review_note"`
		Breakdown         Breakdown      `json:"breakdown"`
		Packages          []PackageLevel `json:"packages"`
		KartuStok         []KartuStokRow `json:"kartu_stok"`
	}
)
