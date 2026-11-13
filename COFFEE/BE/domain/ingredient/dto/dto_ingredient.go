package dto

// Jenis entri riwayat stok bahan (lihat migrasi 003).
const (
	MutationUsage  = "usage"  // terpakai
	MutationWaste  = "waste"  // terbuang
	MutationCancel = "cancel" // pembatalan entri terpakai/terbuang
)

type (
	// REQUEST
	GetAllRequest struct {
		Page      int    `json:"page" validate:"required,min=1"`
		Limit     int    `json:"limit" validate:"required,min=1"`
		Search    string `json:"search" validate:"max=100"`
		IsActive  *bool  `json:"is_active"`
		LowStock  bool   `json:"low_stock"`
		SortBy    string `json:"sort_by"`
		SortOrder string `json:"sort_order"`
	}

	IDUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		Name     string   `json:"name" validate:"required,max=100,notblank"`
		UnitID   int      `json:"unit_id" validate:"required,min=1"`
		MinStock *float64 `json:"min_stock" validate:"omitempty,gte=0,lte=999999999"`
	}

	UpdateRequest struct {
		ID       int      `json:"-"`
		Name     string   `json:"name" validate:"required,max=100,notblank"`
		UnitID   int      `json:"unit_id" validate:"required,min=1"`
		MinStock *float64 `json:"min_stock" validate:"omitempty,gte=0,lte=999999999"`
	}

	// RecordUsageRequest: catat terpakai (satu bahan per input, keputusan 14.10).
	RecordUsageRequest struct {
		IngredientID int     `json:"ingredient_id" validate:"required,min=1"`
		Qty          float64 `json:"qty" validate:"required,gt=0,lte=999999999"`
		Notes        string  `json:"notes" validate:"max=255"`
	}

	// RecordWasteRequest: catat terbuang, alasan wajib (keputusan 83).
	RecordWasteRequest struct {
		IngredientID int     `json:"ingredient_id" validate:"required,min=1"`
		Qty          float64 `json:"qty" validate:"required,gt=0,lte=999999999"`
		Reason       string  `json:"reason" validate:"required,max=255,notblank"`
	}

	GetMutationsRequest struct {
		Page         int    `json:"page" validate:"required,min=1"`
		Limit        int    `json:"limit" validate:"required,min=1"`
		IngredientID *int   `json:"ingredient_id"`
		Type         string `json:"type" validate:"omitempty,oneof=usage waste cancel"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
	}

	// CancelMutationRequest: batalkan entri terpakai/terbuang (keputusan 81), alasan wajib.
	CancelMutationRequest struct {
		ID     int    `json:"-"`
		Reason string `json:"reason" validate:"required,max=255,notblank"`
	}

	// RESPONSE
	IngredientResponse struct {
		ID          int      `json:"id"`
		Name        string   `json:"name"`
		UnitID      int      `json:"unit_id"`
		UnitName    string   `json:"unit_name"`
		UnitAbbr    string   `json:"unit_abbr"`
		Stock       float64  `json:"stock"`
		TotalUsed   float64  `json:"total_used"`
		TotalWasted float64  `json:"total_wasted"`
		MinStock    *float64 `json:"min_stock"`
		IsLowStock  bool     `json:"is_low_stock"`
		IsActive    bool     `json:"is_active"`
		HasHistory  bool     `json:"has_history"`
	}

	IngredientOptionResponse struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		UnitAbbr string  `json:"unit_abbr"`
		Stock    float64 `json:"stock"`
	}

	MutationResponse struct {
		ID              int     `json:"id"`
		IngredientID    int     `json:"ingredient_id"`
		IngredientName  string  `json:"ingredient_name"`
		UnitAbbr        string  `json:"unit_abbr"`
		Type            string  `json:"type"`
		QtyChange       float64 `json:"qty_change"`
		StockBefore     float64 `json:"stock_before"`
		StockAfter      float64 `json:"stock_after"`
		Notes           string  `json:"notes"`
		ReferenceID     *int    `json:"reference_id"`
		IsCancelled     bool    `json:"is_cancelled"`
		CancelledAt     string  `json:"cancelled_at,omitempty"`
		CancelledByName string  `json:"cancelled_by_name,omitempty"`
		CancelReason    string  `json:"cancel_reason,omitempty"`
		CreatedByName   string  `json:"created_by_name"`
		CreatedAt       string  `json:"created_at"`
	}

	// RecordResponse: hasil catat terpakai/terbuang. NegativeStock = stok menjadi minus
	// (diizinkan dengan peringatan, keputusan 14.9).
	RecordResponse struct {
		Mutation      MutationResponse `json:"mutation"`
		StockAfter    float64          `json:"stock_after"`
		NegativeStock bool             `json:"negative_stock"`
	}
)
