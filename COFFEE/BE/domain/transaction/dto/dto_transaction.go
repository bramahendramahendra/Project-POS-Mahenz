package dto

import "time"

type (
	// REQUEST
	GetAllRequest struct {
		Page          int    `json:"page" validate:"required,min=1"`
		Limit         int    `json:"limit" validate:"required,min=1"`
		Search        string `json:"search"`
		Status        string `json:"status"`
		PaymentMethod string `json:"payment_method"`
		OrderType     string `json:"order_type" validate:"omitempty,oneof=dine_in take_away"`
		DateFrom      string `json:"date_from"`
		DateTo        string `json:"date_to"`
		UserID        *int   `json:"user_id"`
		SortBy        string `json:"sort_by"`
		SortOrder     string `json:"sort_order"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	VoidRequest struct {
		ID int `uri:"id" validate:"required,min=1"`

		Body VoidBodyRequest `json:"-"`
	}

	// VoidItemAction: pilihan per item barang jadi saat void (keputusan 100): return = kembali ke
	// stok (bawaan bila item tidak disebut), waste = dibuang (alasan wajib).
	VoidItemAction struct {
		ItemID int    `json:"item_id" validate:"required,min=1"`
		Action string `json:"action" validate:"required,oneof=return waste"`
	}

	VoidBodyRequest struct {
		Items  []VoidItemAction `json:"items" validate:"max=100,dive"`
		Reason string           `json:"reason" validate:"max=255"`
	}

	// CreateTransactionItemRequest: satu baris pesanan Coffee (tahap 4.7a). Harga dihitung server dari
	// kombinasi varian + tambahan; klien hanya mengirim pilihan. Kolom ber-tag json:"-" diisi service.
	CreateTransactionItemRequest struct {
		VariantID         int     `json:"variant_id" validate:"required,min=1"`
		Quantity          float64 `json:"quantity" validate:"required,min=1,max=999"`
		ModifierOptionIDs []int   `json:"modifier_option_ids" validate:"max=30,dive,min=1"`
		Notes             string  `json:"notes" validate:"max=200"`
		DiscountItem      float64 `json:"discount_item" validate:"gte=0"`

		MenuItemID  int            `json:"-"`
		MenuName    string         `json:"-"`
		VariantName string         `json:"-"`
		ProductID   *int           `json:"-"` // terisi hanya untuk menu barang jadi (stok dipotong)
		Price       float64        `json:"-"` // harga satuan = varian + tambahan
		Subtotal    float64        `json:"-"`
		Modifiers   []LineModifier `json:"-"`
	}

	// LineModifier: tambahan terpilih dengan nama dan harga saat dijual.
	LineModifier struct {
		OptionID   int
		GroupName  string
		OptionName string
		Price      float64
	}

	// CreateTransactionRequest: pesanan bayar sekarang (4.7a). Subtotal, total, dan kembalian dihitung server.
	CreateTransactionRequest struct {
		ShiftID *int `json:"shift_id"`
		// OpenBill: simpan sebagai pesanan terbuka (bayar belakangan, keputusan 27); pembayaran diabaikan.
		OpenBill      bool                           `json:"open_bill"`
		OrderType     string                         `json:"order_type" validate:"required,oneof=dine_in take_away"`
		TableID       *int                           `json:"table_id" validate:"omitempty,min=1"`
		OrderName     string                         `json:"order_name" validate:"max=50"`
		Discount      float64                        `json:"discount" validate:"gte=0"`
		Tax           float64                        `json:"tax" validate:"gte=0"`
		PaymentMethod string                         `json:"payment_method" validate:"omitempty,oneof=cash transfer qris card kredit"`
		PaymentAmount float64                        `json:"payment_amount" validate:"min=0"`
		CustomerID    *int                           `json:"customer_id"`
		IsCredit      bool                           `json:"is_credit"`
		DeviceSource  string                         `json:"device_source" validate:"required,oneof=desktop web android"`
		Items         []CreateTransactionItemRequest `json:"items" validate:"required,min=1,max=100,dive"`

		TableName    string  `json:"-"`
		Subtotal     float64 `json:"-"`
		TotalAmount  float64 `json:"-"`
		ChangeAmount float64 `json:"-"`
	}

	// RESPONSE
	ItemModifierResponse struct {
		OptionID   *int    `json:"option_id"`
		GroupName  string  `json:"group_name"`
		OptionName string  `json:"option_name"`
		Price      float64 `json:"price"`
	}

	TransactionItemResponse struct {
		ID            int                    `json:"id"`
		MenuItemID    *int                   `json:"menu_item_id"`
		VariantID     *int                   `json:"variant_id"`
		ProductID     *int                   `json:"product_id"`   // terisi = barang jadi (stok dipotong)
		ProductName   string                 `json:"product_name"` // nama menu saat dijual
		VariantName   string                 `json:"variant_name"`
		Notes         string                 `json:"notes"`
		Quantity      float64                `json:"quantity"`
		Unit          string                 `json:"unit"`
		Price         float64                `json:"price"`
		Subtotal      float64                `json:"subtotal"`
		DiscountItem  float64                `json:"discount_item"`
		ConversionQty float64                `json:"conversion_qty"`
		UnitID        *int                   `json:"unit_id"`
		Modifiers     []ItemModifierResponse `json:"modifiers"`
	}

	TransactionResponse struct {
		ID              int       `json:"id"`
		TransactionCode string    `json:"transaction_code"`
		UserID          int       `json:"user_id"`
		KasirName       string    `json:"kasir_name"`
		ShiftID         *int      `json:"shift_id"`
		TransactionDate time.Time `json:"transaction_date"`
		Subtotal        float64   `json:"subtotal"`
		Discount        float64   `json:"discount"`
		Tax             float64   `json:"tax"`
		TotalAmount     float64   `json:"total_amount"`
		PaymentMethod   string    `json:"payment_method"`
		PaymentAmount   float64   `json:"payment_amount"`
		ChangeAmount    float64   `json:"change_amount"`
		BalanceUsed     float64   `json:"balance_used"`
		CustomerID      *int      `json:"customer_id"`
		CustomerName    string    `json:"customer_name"`
		IsCredit        bool      `json:"is_credit"`
		Status          string    `json:"status"`
		DeviceSource    string    `json:"device_source"`
		CashDrawerID    *int      `json:"cash_drawer_id"`
		OrderType       string    `json:"order_type"`
		TableID         *int      `json:"table_id"`
		TableName       string    `json:"table_name"`
		QueueNumber     int       `json:"queue_number"`
		QueueLabel      string    `json:"queue_label"` // 001, 002, ...
		OrderName       string    `json:"order_name"`
		// Open bill (tahap 4.8): pembuat pesanan, gabung, batal, dan catatan item yang dihapus.
		CreatedBy     *int   `json:"created_by"`
		CreatedByName string `json:"created_by_name"`
		MergedIntoID  *int   `json:"merged_into_id"`
		CancelReason  string `json:"cancel_reason"`
		// Pelaku penyelesaian pembayaran dan tanda pesanan tertunda (tahap 4.9).
		SettledByName string                     `json:"settled_by_name"`
		IsOverdue     bool                       `json:"is_overdue"`
		Cancellations []ItemCancellationResponse `json:"cancellations,omitempty"`
		Items         []TransactionItemResponse  `json:"items,omitempty"`
	}

	CreateTransactionResponse struct {
		TransactionResponse
	}
)
