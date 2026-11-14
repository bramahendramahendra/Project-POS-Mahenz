package model

import "time"

type Transaction struct {
	ID              int       `gorm:"column:id"`
	TransactionCode string    `gorm:"column:transaction_code"`
	UserID          int       `gorm:"column:user_id"`
	ShiftID         *int      `gorm:"column:shift_id"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	Subtotal        float64   `gorm:"column:subtotal"`
	Discount        float64   `gorm:"column:discount"`
	Tax             float64   `gorm:"column:tax"`
	TotalAmount     float64   `gorm:"column:total_amount"`
	PaymentMethod   string    `gorm:"column:payment_method"`
	PaymentAmount   float64   `gorm:"column:payment_amount"`
	ChangeAmount    float64   `gorm:"column:change_amount"`
	BalanceUsed     float64   `gorm:"column:balance_used"`
	CustomerID      *int      `gorm:"column:customer_id"`
	IsCredit        bool      `gorm:"column:is_credit"`
	Status          string    `gorm:"column:status"`
	DeviceSource    string    `gorm:"column:device_source"`
	CashDrawerID    *int      `gorm:"column:cash_drawer_id"`
	CreatedBy       *int      `gorm:"column:created_by"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

// TransactionItem: baris pesanan. ProductID > 0 hanya untuk menu barang jadi (stok dipotong);
// menu racikan ProductID = 0.
type TransactionItem struct {
	ID            int     `gorm:"column:id"`
	TransactionID int     `gorm:"column:transaction_id"`
	MenuItemID    *int    `gorm:"column:menu_item_id"`
	VariantID     *int    `gorm:"column:variant_id"`
	ProductID     int     `gorm:"column:product_id"`
	ProductName   string  `gorm:"column:product_name"`
	VariantName   string  `gorm:"column:variant_name"`
	Notes         string  `gorm:"column:notes"`
	Quantity      float64 `gorm:"column:quantity"`
	Unit          string  `gorm:"column:unit"`
	Price         float64 `gorm:"column:price"`
	PurchasePrice float64 `gorm:"column:purchase_price"`
	Subtotal      float64 `gorm:"column:subtotal"`
	DiscountItem  float64 `gorm:"column:discount_item"`
	ConversionQty float64 `gorm:"column:conversion_qty"`
	UnitID        *int    `gorm:"column:unit_id"`
}

// TransactionItemModifier: tambahan yang dipilih pada satu baris pesanan.
type TransactionItemModifier struct {
	TransactionItemID int     `gorm:"column:transaction_item_id"`
	OptionID          *int    `gorm:"column:option_id"`
	GroupName         string  `gorm:"column:group_name"`
	OptionName        string  `gorm:"column:option_name"`
	Price             float64 `gorm:"column:price"`
}

// SaleVariant: data kombinasi varian beserta menunya untuk menghitung dan memeriksa pesanan.
type SaleVariant struct {
	VariantID     int     `gorm:"column:variant_id"`
	MenuItemID    int     `gorm:"column:menu_item_id"`
	MenuName      string  `gorm:"column:menu_name"`
	StockMode     string  `gorm:"column:stock_mode"`
	ProductID     *int    `gorm:"column:product_id"`
	MenuActive    bool    `gorm:"column:menu_active"`
	SoldOut       bool    `gorm:"column:sold_out"`
	VariantActive bool    `gorm:"column:variant_active"`
	Price         float64 `gorm:"column:price"`
	SizeName      string  `gorm:"column:size_name"`
	Temperature   *string `gorm:"column:temperature"`
}

// SaleModifierOption: pilihan tambahan dari kelompok yang terpasang di menu.
type SaleModifierOption struct {
	OptionID      int     `gorm:"column:option_id"`
	OptionName    string  `gorm:"column:option_name"`
	Price         float64 `gorm:"column:price"`
	OptionActive  bool    `gorm:"column:option_active"`
	GroupID       int     `gorm:"column:group_id"`
	GroupName     string  `gorm:"column:group_name"`
	SelectionType string  `gorm:"column:selection_type"`
	IsRequired    bool    `gorm:"column:is_required"`
	GroupActive   bool    `gorm:"column:group_active"`
}

type StockMutation struct {
	ID            int     `gorm:"column:id"`
	ProductID     int     `gorm:"column:product_id"`
	MutationType  string  `gorm:"column:mutation_type"`
	Quantity      float64 `gorm:"column:quantity"`
	StockBefore   float64 `gorm:"column:stock_before"`
	StockAfter    float64 `gorm:"column:stock_after"`
	ReferenceType string  `gorm:"column:reference_type"`
	ReferenceID   int     `gorm:"column:reference_id"`
	Notes         string  `gorm:"column:notes"`
	UserID        int     `gorm:"column:user_id"`
}
