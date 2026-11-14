package dto

// OrderTypeItem: penjualan per jenis pesanan (Coffee, keputusan 117).
type OrderTypeItem struct {
	OrderType string  `json:"order_type"`
	Total     float64 `json:"total"`
	Count     int64   `json:"count"`
}
