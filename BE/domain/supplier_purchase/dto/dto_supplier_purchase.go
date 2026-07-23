package dto

type (
	// REQUEST
	GetAllRequest struct {
		Page          int    `json:"page" validate:"required,min=1"`
		Limit         int    `json:"limit" validate:"required,min=1"`
		Search        string `json:"search"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		SupplierID    *int   `json:"supplier_id"`
		PaymentStatus string `json:"payment_status"`
		SortBy        string `json:"sort_by"`
		SortOrder     string `json:"sort_order"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	VoidUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	GenerateCodeResponse struct {
		PurchaseCode string `json:"purchase_code"`
	}

	// ExpiryBatchDraft: satu baris rincian tanggal expired untuk sebagian qty item PO.
	// Opsional — dikirim hanya kalau staf mencentang "Produk ini ada tanggal expired".
	// Total Qty di semua baris untuk 1 item harus persis sama dengan PurchaseRequest.Quantity
	// (dicek di service, bukan lewat tag validate karena butuh baca field lain).
	ExpiryBatchDraft struct {
		Qty         float64 `json:"qty" validate:"required,gt=0"`
		ExpiredDate string  `json:"expired_date" validate:"required"`
	}

	// PackageID: id product_packages yang dipilih user (mis. "Slop (Lama)" vs "Slop
	// (Baru)"). Kalau diisi, conversion_qty dihitung ulang di server dari sini —
	// ConversionQty yang dikirim klien cuma dipakai sebagai fallback kalau PackageID
	// kosong (mis. baris lama sebelum redesain ini).
	PurchaseRequest struct {
		ProductID     int                `json:"product_id" validate:"required,gt=0"`
		PackageID     *int               `json:"package_id"`
		Quantity      float64            `json:"quantity" validate:"required,gt=0"`
		Unit          string             `json:"unit" validate:"required"`
		ConversionQty float64            `json:"conversion_qty"`
		PurchasePrice float64            `json:"purchase_price" validate:"required,gt=0"`
		ExpiryBatches []ExpiryBatchDraft `json:"expiry_batches" validate:"omitempty,dive"`
	}

	CreateRequest struct {
		ID             int               `json:"-"`
		UserID         int               `json:"-"`
		InvoiceNumber  string            `json:"invoice_number" validate:"required"`
		SupplierID     *int              `json:"supplier_id" validate:"required,gt=0"`
		PurchaseDate   string            `json:"purchase_date" validate:"required"`
		DiscountAmount float64           `json:"discount_amount"`
		PaymentStatus  string            `json:"payment_status"`
		PaidAmount     float64           `json:"paid_amount"`
		PaymentMethod  string            `json:"payment_method"`
		Notes          string            `json:"notes"`
		Items          []PurchaseRequest `json:"items" validate:"required,min=1,dive"`
	}

	UpdateRequest struct {
		ID             int               `json:"-"`
		UserID         int               `json:"-"`
		InvoiceNumber  string            `json:"invoice_number" validate:"required"`
		SupplierID     *int              `json:"supplier_id" validate:"required,gt=0"`
		PurchaseDate   string            `json:"purchase_date" validate:"required"`
		DiscountAmount float64           `json:"discount_amount"`
		PaymentStatus  string            `json:"payment_status"`
		PaidAmount     float64           `json:"paid_amount"`
		PaymentMethod  string            `json:"payment_method"`
		Notes          string            `json:"notes"`
		Items          []PurchaseRequest `json:"items" validate:"required,min=1,dive"`
	}

	AddItemsUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	AddItemsRequest struct {
		ID     int               `json:"-"`
		UserID int               `json:"-"`
		Items  []PurchaseRequest `json:"items" validate:"required,min=1,dive"`
	}

	PayRequest struct {
		ID            int     `json:"-"`
		UserID        int     `json:"-"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		PaymentDate   string  `json:"payment_date" validate:"required"`
		PaymentMethod string  `json:"payment_method" validate:"required"`
		Notes         string  `json:"notes"`
	}

	// RESPONSE
	PaymentResponse struct {
		ID            int     `json:"id"`
		PaymentDate   string  `json:"payment_date"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		Notes         string  `json:"notes"`
		UserName      string  `json:"user_name"`
		CreatedAt     string  `json:"created_at"`
	}

	PurchaseItemResponse struct {
		ID            int     `json:"id"`
		ProductID     int     `json:"product_id"`
		ProductName   string  `json:"product_name"`
		Quantity      float64 `json:"quantity"`
		Unit          string  `json:"unit"`
		ConversionQty float64 `json:"conversion_qty"`
		PurchasePrice float64 `json:"purchase_price"`
		Subtotal      float64 `json:"subtotal"`
	}

	PurchaseResponse struct {
		ID              int                    `json:"id"`
		PurchaseCode    string                 `json:"purchase_code"`
		InvoiceNumber   string                 `json:"invoice_number"`
		SupplierID      *int                   `json:"supplier_id"`
		SupplierName    string                 `json:"supplier_name"`
		PurchaseDate    string                 `json:"purchase_date"`
		DiscountAmount  float64                `json:"discount_amount"`
		TotalAmount     float64                `json:"total_amount"`
		PaymentStatus   string                 `json:"payment_status"`
		PaidAmount      float64                `json:"paid_amount"`
		RemainingAmount float64                `json:"remaining_amount"`
		Status          string                 `json:"status"`
		UserName        string                 `json:"user_name"`
		Notes           string                 `json:"notes"`
		Items           []PurchaseItemResponse `json:"items,omitempty"`
	}
)
