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

	// WarningResponse: satu baris batch yang statusnya masih "active" DAN expired_date-nya
	// sudah mendekati/lewat hari ini — dihitung murni dari tanggal, bukan dari estimasi
	// sisa stok (lihat catatan desain di migrasi 003_expiry_batch.sql).
	WarningResponse struct {
		ID          int     `json:"id"`
		ProductID   int     `json:"product_id"`
		ProductName string  `json:"product_name"`
		Qty         float64 `json:"qty"`
		ExpiredDate string  `json:"expired_date"`
		Severity    string  `json:"severity"` // "near" | "expired"
		DaysLeft    int     `json:"days_left"`
	}

	// ProductSeverityResponse: ringkasan per produk (severity terburuk + jumlah batch
	// warning) dipakai FE untuk badge di list Produk tanpa perlu load semua detail batch.
	ProductSeverityResponse struct {
		ProductID    int    `json:"product_id"`
		Severity     string `json:"severity"` // "near" | "expired"
		WarningCount int    `json:"warning_count"`
	}

	// BatchHistoryResponse: SEMUA batch expired produk ini apapun statusnya — dipakai di
	// Detail Produk supaya staf bisa lihat histori lengkap (bukan cuma yang lagi warning).
	BatchHistoryResponse struct {
		ID          int     `json:"id"`
		Qty         float64 `json:"qty"`
		ExpiredDate string  `json:"expired_date"`
		Status      string  `json:"status"` // "active" | "cleared" | "written_off"
	}
)
