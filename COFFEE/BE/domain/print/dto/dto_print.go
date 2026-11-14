package dto

// PrintReceiptItem adalah satu baris item pada struk.
type PrintReceiptItem struct {
	Name     string  `json:"name"`
	UnitName string  `json:"unit_name"`
	Qty      float64 `json:"qty"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
	// Details: baris rincian di bawah item (Coffee: kombinasi, tambahan, catatan).
	Details []string `json:"details"`
}

// PrintReceiptRequest adalah payload untuk mencetak struk langsung ke printer
// thermal. FE mengirim data yang sudah siap tampil; BE yang menyusun layout
// ESC/POS agar lebar pas kertas.
//
// Target menentukan tujuan cetak:
//
//	"" atau "server" : printer di komputer server via COM port (default)
//	"lan"            : printer jaringan via TCP ke LanHost:LanPort. Bila LanHost
//	                   kosong, dipakai nilai dari config (Printer.LanHost/LanPort).
type PrintReceiptRequest struct {
	Target  string `json:"target"`
	LanHost string `json:"lan_host"`
	LanPort int    `json:"lan_port"`

	StoreName       string `json:"store_name"`
	StoreSub        string `json:"store_sub"`
	Footer          string `json:"footer"`
	TransactionCode string `json:"transaction_code"`
	TransactionDate string `json:"transaction_date"`
	CustomerName    string `json:"customer_name"`
	// Info pesanan Coffee (tahap 4.7b): nomor antrean dicetak besar, jenis pesanan, meja, nama pemesan.
	QueueLabel string `json:"queue_label"`
	OrderType  string `json:"order_type"`
	TableName  string `json:"table_name"`
	OrderName  string `json:"order_name"`
	// Unpaid: tagihan sementara pesanan terbuka (keputusan 108), tanpa baris pembayaran.
	Unpaid         bool               `json:"unpaid"`
	PaymentMethod  string             `json:"payment_method"`
	Items          []PrintReceiptItem `json:"items"`
	Subtotal       float64            `json:"subtotal"`
	DiscountAmount float64            `json:"discount_amount"`
	TaxAmount      float64            `json:"tax_amount"`
	TaxPercent     float64            `json:"tax_percent"`
	GrandTotal     float64            `json:"grand_total"`
	AmountPaid     float64            `json:"amount_paid"`
	Change         float64            `json:"change"`
}
