package dto

type SupplierReturnItemRequest struct {
	PurchaseItemID int     `json:"purchase_item_id" validate:"required,gt=0"`
	ProductID      int     `json:"product_id" validate:"required,gt=0"`
	ProductName    string  `json:"product_name" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	Unit           string  `json:"unit" validate:"required"`
	PurchasePrice  float64 `json:"purchase_price" validate:"required,gt=0"`
}

type GetSupplierReturnByIDRequest struct {
	ID int `uri:"id" validate:"required,min=1"`
}

type CreateSupplierReturnRequest struct {
	UserID       int                         `json:"-"`
	PurchaseID   int                         `json:"purchase_id" validate:"required,gt=0"`
	SupplierID   *int                        `json:"supplier_id"`
	SupplierName string                      `json:"supplier_name" validate:"required"`
	ReturnDate   string                      `json:"return_date" validate:"required,dateformat"`
	Reason       string                      `json:"reason" validate:"required"`
	Notes        string                      `json:"notes" validate:"omitempty"`
	Items        []SupplierReturnItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateStatusRequest struct {
	ID     int    `json:"-"`
	UserID int    `json:"-"`
	Status string `json:"status" validate:"required,oneof=approved rejected"`
	Notes  string `json:"notes"`
}

type SupplierReturnItemResponse struct {
	ID            int     `json:"id"`
	ProductID     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
	PurchasePrice float64 `json:"purchase_price"`
	Subtotal      float64 `json:"subtotal"`
}

type SupplierReturnResponse struct {
	ID                int                          `json:"id"`
	ReturnCode        string                       `json:"return_code"`
	PurchaseID        int                          `json:"purchase_id"`
	SupplierID        *int                         `json:"supplier_id"`
	SupplierName      string                       `json:"supplier_name"`
	ReturnDate        string                       `json:"return_date"`
	TotalReturnAmount float64                      `json:"total_return_amount"`
	Reason            string                       `json:"reason"`
	Status            string                       `json:"status"`
	UserName          string                       `json:"user_name"`
	Notes             string                       `json:"notes"`
	Items             []SupplierReturnItemResponse `json:"items,omitempty"`
}

type SupplierReturnListRequest struct {
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SupplierID *int   `json:"supplier_id"`
	Status     string `json:"status"`
	SortBy     string `json:"sort_by"`
	SortOrder  string `json:"sort_order"`
}
