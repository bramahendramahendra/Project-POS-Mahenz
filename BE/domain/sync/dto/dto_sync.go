package dto

import "time"

type SyncItem struct {
	LocalID     string `json:"local_id"`
	ServerID    int    `json:"server_id"`
	EntityType  string `json:"entity_type"`
	EntityID    int    `json:"entity_id"`
	Action      string `json:"action"`
	Payload     string `json:"payload"`
	DesktopTime string `json:"desktop_time"`
	UpdatedAt   string `json:"updated_at"`
}

type PushSyncRequest struct {
	DeviceID   string     `json:"device_id" validate:"required"`
	DeviceType string     `json:"device_type"`
	Items      []SyncItem `json:"items" validate:"required"`
}

type SyncItemResult struct {
	LocalID    string `json:"local_id"`
	Status     string `json:"status"`
	ServerID   int    `json:"server_id,omitempty"`
	ConflictID int    `json:"conflict_id,omitempty"`
	Message    string `json:"message,omitempty"`
}

type PushSyncResponse struct {
	Processed int              `json:"processed"`
	Conflicts int              `json:"conflicts"`
	Failed    int              `json:"failed"`
	Pending   int              `json:"pending"`
	Results   []SyncItemResult `json:"results"`
}

type ConflictFilter struct {
	Status string
	Page   int
	Limit  int
}

type ConflictResponse struct {
	ID          int       `json:"id"`
	EntityType  string    `json:"entity_type"`
	EntityID    int       `json:"entity_id"`
	LocalID     string    `json:"local_id"`
	DeviceID    string    `json:"device_id"`
	DesktopData string    `json:"desktop_data"`
	OnlineData  string    `json:"online_data"`
	DesktopTime time.Time `json:"desktop_time"`
	OnlineTime  time.Time `json:"online_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConflictListResponse struct {
	Data  []ConflictResponse `json:"data"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type ResolveConflictRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
}

type QueueFilter struct {
	DeviceID   string
	Status     string
	EntityType string
	Page       int
	Limit      int
}

type QueueResponse struct {
	ID         int       `json:"id"`
	DeviceID   string    `json:"device_id"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	Action     string    `json:"action"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type QueueListResponse struct {
	Data  []QueueResponse `json:"data"`
	Total int             `json:"total"`
}

type SyncTransactionItemPayload struct {
	ProductID     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
	Price         float64 `json:"price"`
	Subtotal      float64 `json:"subtotal"`
	DiscountItem  float64 `json:"discount_item"`
	ConversionQty float64 `json:"conversion_qty"`
	UnitID        *int    `json:"unit_id"`
}

type SyncTransactionPayload struct {
	UserID        int                          `json:"user_id"`
	ShiftID       *int                         `json:"shift_id"`
	Subtotal      float64                      `json:"subtotal"`
	Discount      float64                      `json:"discount"`
	Tax           float64                      `json:"tax"`
	TotalAmount   float64                      `json:"total_amount"`
	PaymentMethod string                       `json:"payment_method"`
	PaymentAmount float64                      `json:"payment_amount"`
	ChangeAmount  float64                      `json:"change_amount"`
	CustomerID    *int                         `json:"customer_id"`
	IsCredit      bool                         `json:"is_credit"`
	DeviceSource  string                       `json:"device_source"`
	Items         []SyncTransactionItemPayload `json:"items"`
}

type HistoryFilter struct {
	DeviceID  string
	StartDate string
	EndDate   string
	Page      int
	Limit     int
}

type SyncHistoryResponse struct {
	ID            int64   `json:"id"`
	DeviceID      string  `json:"device_id"`
	DeviceType    string  `json:"device_type"`
	TotalItems    int     `json:"total_items"`
	SyncedItems   int     `json:"synced_items"`
	ConflictItems int     `json:"conflict_items"`
	FailedItems   int     `json:"failed_items"`
	PendingItems  int     `json:"pending_items"`
	DurationMs    int     `json:"duration_ms"`
	Status        string  `json:"status"`
	StartedAt     string  `json:"started_at"`
	FinishedAt    *string `json:"finished_at"`
}

type SyncHistoryListResponse struct {
	Data  []SyncHistoryResponse `json:"data"`
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}
