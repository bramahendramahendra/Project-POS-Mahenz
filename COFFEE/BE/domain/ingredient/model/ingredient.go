package model

import "time"

// Ingredient: bahan baku. Stok dalam satu satuan (master Unit), boleh desimal.
type Ingredient struct {
	ID          int       `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	UnitID      int       `gorm:"column:unit_id"`
	UnitName    string    `gorm:"column:unit_name"`
	UnitAbbr    string    `gorm:"column:unit_abbr"`
	Stock       float64   `gorm:"column:stock"`
	TotalUsed   float64   `gorm:"column:total_used"`
	TotalWasted float64   `gorm:"column:total_wasted"`
	MinStock    *float64  `gorm:"column:min_stock"`
	IsActive    bool      `gorm:"column:is_active"`
	HasHistory  bool      `gorm:"column:has_history"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

// IngredientMutation: satu baris riwayat stok bahan. Tidak pernah diubah atau dihapus,
// kecuali penanda pembatalan (cancelled_*).
type IngredientMutation struct {
	ID              int        `gorm:"column:id"`
	IngredientID    int        `gorm:"column:ingredient_id"`
	IngredientName  string     `gorm:"column:ingredient_name"`
	UnitAbbr        string     `gorm:"column:unit_abbr"`
	Type            string     `gorm:"column:type"`
	QtyChange       float64    `gorm:"column:qty_change"`
	StockBefore     float64    `gorm:"column:stock_before"`
	StockAfter      float64    `gorm:"column:stock_after"`
	Notes           string     `gorm:"column:notes"`
	ReferenceType   string     `gorm:"column:reference_type"`
	ReferenceID     *int       `gorm:"column:reference_id"`
	CancelledAt     *time.Time `gorm:"column:cancelled_at"`
	CancelledBy     *int       `gorm:"column:cancelled_by"`
	CancelledByName string     `gorm:"column:cancelled_by_name"`
	CancelReason    string     `gorm:"column:cancel_reason"`
	CreatedBy       int        `gorm:"column:created_by"`
	CreatedByName   string     `gorm:"column:created_by_name"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
}
