package dto

// PrintReceiptItem adalah satu baris item pada struk.
type PrintReceiptItem struct {
	Name     string  `json:"name"`
	UnitName string  `json:"unit_name"`
	Qty      float64 `json:"qty"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

// PrintReceiptRequest adalah payload untuk mencetak struk langsung ke printer
// thermal via COM port. FE mengirim data yang sudah siap tampil; BE yang
// menyusun layout ESC/POS agar lebar pas kertas.
type PrintReceiptRequest struct {
	StoreName       string             `json:"store_name"`
	StoreSub        string             `json:"store_sub"`
	Footer          string             `json:"footer"`
	TransactionCode string             `json:"transaction_code"`
	TransactionDate string             `json:"transaction_date"`
	CustomerName    string             `json:"customer_name"`
	PaymentMethod   string             `json:"payment_method"`
	Items           []PrintReceiptItem `json:"items"`
	Subtotal        float64            `json:"subtotal"`
	DiscountAmount  float64            `json:"discount_amount"`
	TaxAmount       float64            `json:"tax_amount"`
	TaxPercent      float64            `json:"tax_percent"`
	GrandTotal      float64            `json:"grand_total"`
	AmountPaid      float64            `json:"amount_paid"`
	Change          float64            `json:"change"`
}
