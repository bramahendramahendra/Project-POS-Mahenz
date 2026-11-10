package dto

type (
	// REQUEST
	GetWarningsRequest struct {
		Search string `json:"search"`
	}

	GetByProductUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ConfirmUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	ConfirmRequest struct {
		ID     int    `json:"-"`
		UserID int    `json:"-"`
		Notes  string `json:"notes"`
	}

	WriteOffUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	WriteOffRequest struct {
		ID     int    `json:"-"`
		UserID int    `json:"-"`
		Notes  string `json:"notes"`
	}

	// RESPONSE

	WarningResponse struct {
		ID          int     `json:"id"`
		ProductID   int     `json:"product_id"`
		ProductName string  `json:"product_name"`
		UnitName    string  `json:"unit_name"`
		Qty         float64 `json:"qty"`
		ExpiredDate string  `json:"expired_date"`
		Severity    string  `json:"severity"` // "near" | "expired"
		DaysLeft    int     `json:"days_left"`
	}

	ProductSeverityResponse struct {
		ProductID    int    `json:"product_id"`
		Severity     string `json:"severity"` // "near" | "expired"
		WarningCount int    `json:"warning_count"`
	}

	BatchHistoryResponse struct {
		ID          int     `json:"id"`
		UnitName    string  `json:"unit_name"`
		Qty         float64 `json:"qty"`
		ExpiredDate string  `json:"expired_date"`
		Status      string  `json:"status"` // "active" | "cleared" | "written_off"
	}
)
