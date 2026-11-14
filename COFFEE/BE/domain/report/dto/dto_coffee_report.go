package dto

// Laporan Coffee (tahap 4.10, keputusan 114-117).

type CoffeeReportRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

// ─── Rincian penjualan (keputusan 115) ─────────────────────────────
// Pendapatan item = subtotal item dikurangi bagian diskon nota secara proporsional, tanpa pajak
// (sama dengan Laba Rugi).

type SalesVariantRow struct {
	VariantName string  `json:"variant_name"`
	Qty         float64 `json:"qty"`
	Revenue     float64 `json:"revenue"`
}

type SalesMenuRow struct {
	MenuItemID   int               `json:"menu_item_id"`
	MenuName     string            `json:"menu_name"`
	CategoryName string            `json:"category_name"`
	Qty          float64           `json:"qty"`
	Revenue      float64           `json:"revenue"`
	Share        float64           `json:"share"` // persen dari total pendapatan
	Variants     []SalesVariantRow `json:"variants"`
}

type SalesCategoryRow struct {
	CategoryName string  `json:"category_name"`
	Qty          float64 `json:"qty"`
	Revenue      float64 `json:"revenue"`
	Share        float64 `json:"share"`
}

// SalesGroupRow: per metode bayar atau per jenis pesanan (total transaksi termasuk pajak).
type SalesGroupRow struct {
	Key   string  `json:"key"`
	Count int     `json:"count"`
	Total float64 `json:"total"`
}

type SalesBreakdownResponse struct {
	TotalRevenue  float64            `json:"total_revenue"`
	TotalQty      float64            `json:"total_qty"`
	Menus         []SalesMenuRow     `json:"menus"`
	Categories    []SalesCategoryRow `json:"categories"`
	PaymentMethod []SalesGroupRow    `json:"payment_methods"`
	OrderTypes    []SalesGroupRow    `json:"order_types"`
}

// ─── Bahan baku (keputusan 114, 116) ────────────────────────────────

type IngredientReportRow struct {
	IngredientID  int     `json:"ingredient_id"`
	Name          string  `json:"name"`
	UnitAbbr      string  `json:"unit_abbr"`
	PurchasedQty  float64 `json:"purchased_qty"`
	PurchaseValue float64 `json:"purchase_value"`
	UsedQty       float64 `json:"used_qty"`
	WastedQty     float64 `json:"wasted_qty"`
	Stock         float64 `json:"stock"`
	AvgPrice      float64 `json:"avg_price"`
	UsedValue     float64 `json:"used_value"`
	WastedValue   float64 `json:"wasted_value"`
}

type IngredientReportResponse struct {
	TotalPurchaseValue float64               `json:"total_purchase_value"`
	TotalUsedValue     float64               `json:"total_used_value"`
	TotalWastedValue   float64               `json:"total_wasted_value"`
	Items              []IngredientReportRow `json:"items"`
}

// ─── Barang terbuang (keputusan 100, 116) ───────────────────────────

type WasteRow struct {
	CreatedAt string  `json:"created_at"`
	Kind      string  `json:"kind"` // ingredient = bahan baku, product = barang jadi
	ItemName  string  `json:"item_name"`
	Qty       float64 `json:"qty"`
	Unit      string  `json:"unit"`
	Value     float64 `json:"value"`
	Reason    string  `json:"reason"`
	UserName  string  `json:"user_name"`
}

type WasteReportResponse struct {
	TotalValue           float64    `json:"total_value"`
	TotalIngredientValue float64    `json:"total_ingredient_value"`
	TotalProductValue    float64    `json:"total_product_value"`
	Items                []WasteRow `json:"items"`
}
