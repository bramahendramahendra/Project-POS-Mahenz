package dto

import "time"

// Open bill (tahap 4.8, keputusan 27-29, 59, 103-109).
type (
	OpenBillURIRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// OpenAddItemsRequest: tambah item ke pesanan terbuka (stok barang jadi dipotong saat itu, keputusan 40).
	OpenAddItemsRequest struct {
		Items []CreateTransactionItemRequest `json:"items" validate:"required,min=1,max=100,dive"`
	}

	// OpenRemoveItemRequest: hapus item atau kurangi jumlahnya. Quantity 0 = hapus seluruh baris.
	// Barang jadi wajib memilih Action (return/waste); waste wajib Reason (keputusan 100, 105).
	OpenRemoveItemRequest struct {
		ItemID   int     `json:"item_id" validate:"required,min=1"`
		Quantity float64 `json:"quantity" validate:"gte=0,lte=999"`
		Action   string  `json:"action" validate:"omitempty,oneof=return waste"`
		Reason   string  `json:"reason" validate:"max=255"`
	}

	OpenMoveTableRequest struct {
		TableID int `json:"table_id" validate:"required,min=1"`
	}

	// OpenMergeRequest: pesanan terbuka SourceID digabung ke pesanan pada URI (keputusan 107).
	OpenMergeRequest struct {
		SourceID int `json:"source_id" validate:"required,min=1"`
	}

	// OpenPayRequest: bayar pesanan terbuka atau jadikan utang (is_credit + pelanggan). Diskon nota dan
	// pajak diisi saat ini (keputusan 104).
	OpenPayRequest struct {
		ShiftID       *int    `json:"shift_id"`
		Discount      float64 `json:"discount" validate:"gte=0"`
		Tax           float64 `json:"tax" validate:"gte=0"`
		PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash transfer qris card kredit"`
		PaymentAmount float64 `json:"payment_amount" validate:"min=0"`
		CustomerID    *int    `json:"customer_id"`
		IsCredit      bool    `json:"is_credit"`

		Subtotal     float64 `json:"-"`
		TotalAmount  float64 `json:"-"`
		ChangeAmount float64 `json:"-"`
	}

	// OpenCancelRequest: batal pesanan terbuka, alasan wajib; pilihan stok per item barang jadi.
	OpenCancelRequest struct {
		Reason string           `json:"reason" validate:"required,max=255,notblank"`
		Items  []VoidItemAction `json:"items" validate:"max=100,dive"`
	}

	// OpenBillSummary: baris daftar pesanan terbuka di kasir.
	OpenBillSummary struct {
		ID              int       `json:"id"`
		TransactionCode string    `json:"transaction_code"`
		OrderType       string    `json:"order_type"`
		TableID         *int      `json:"table_id"`
		TableName       string    `json:"table_name"`
		QueueNumber     int       `json:"queue_number"`
		QueueLabel      string    `json:"queue_label"`
		OrderName       string    `json:"order_name"`
		Subtotal        float64   `json:"subtotal"`
		ItemCount       float64   `json:"item_count"`
		CreatedByName   string    `json:"created_by_name"`
		CreatedAt       time.Time `json:"created_at"`
		// Tanggal pesanan dan tanda tertunda (dari hari sebelumnya, tahap 4.9).
		BillDate  string `json:"bill_date"`
		IsOverdue bool   `json:"is_overdue"`
	}

	ItemCancellationResponse struct {
		ItemName    string    `json:"item_name"`
		VariantName string    `json:"variant_name"`
		Quantity    float64   `json:"quantity"`
		Price       float64   `json:"price"`
		StockAction string    `json:"stock_action"`
		Reason      string    `json:"reason"`
		UserName    string    `json:"user_name"`
		CreatedAt   time.Time `json:"created_at"`
	}
)
