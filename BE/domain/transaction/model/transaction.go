package model

import "time"

type Transaction struct {
	ID            int        `gorm:"column:id"`
	TransactionCode string   `gorm:"column:transaction_code"`
	UserID        int        `gorm:"column:user_id"`
	ShiftID       *int       `gorm:"column:shift_id"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	Subtotal      float64    `gorm:"column:subtotal"`
	Discount      float64    `gorm:"column:discount"`
	Tax           float64    `gorm:"column:tax"`
	TotalAmount   float64    `gorm:"column:total_amount"`
	PaymentMethod string     `gorm:"column:payment_method"`
	PaymentAmount float64    `gorm:"column:payment_amount"`
	ChangeAmount  float64    `gorm:"column:change_amount"`
	CustomerID    *int       `gorm:"column:customer_id"`
	IsCredit      bool       `gorm:"column:is_credit"`
	Status        string     `gorm:"column:status"`
	DeviceSource  string     `gorm:"column:device_source"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

type TransactionItem struct {
	ID             int     `gorm:"column:id"`
	TransactionID  int     `gorm:"column:transaction_id"`
	ProductID      int     `gorm:"column:product_id"`
	ProductName    string  `gorm:"column:product_name"`
	Quantity       float64 `gorm:"column:quantity"`
	Unit           string  `gorm:"column:unit"`
	Price          float64 `gorm:"column:price"`
	PurchasePrice  float64 `gorm:"column:purchase_price"`
	Subtotal       float64 `gorm:"column:subtotal"`
	DiscountItem   float64 `gorm:"column:discount_item"`
	ConversionQty  float64 `gorm:"column:conversion_qty"`
	UnitID         *int    `gorm:"column:unit_id"`
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

