package dto

import "time"

type GetAllRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	IsActive *bool  `json:"is_active"`
}

type GetByIDRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type CreateRequest struct {
	ParentID   *int   `json:"parent_id"`
	KeyName    string `json:"key_name"    validate:"required,min=2"`
	Label      string `json:"label"       validate:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	OrderIndex int    `json:"order_index"`
}

type UpdateUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type UpdateRequest struct {
	ID         int    `json:"-"`
	ParentID   *int   `json:"parent_id"`
	Label      string `json:"label" validate:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	OrderIndex int    `json:"order_index"`
}

type DeleteRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ReorderItem struct {
	ID         int `json:"id"          validate:"required,gt=0"`
	OrderIndex int `json:"order_index"`
}

type ReorderRequest struct {
	Items []ReorderItem `json:"items" validate:"required,min=1,dive"`
}

type MenuResponse struct {
	ID         int       `json:"id"`
	ParentID   *int      `json:"parent_id"`
	KeyName    string    `json:"key_name"`
	Label      string    `json:"label"`
	Icon       *string   `json:"icon"`
	Path       *string   `json:"path"`
	OrderIndex int       `json:"order_index"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type MenuPermission struct {
	CanView   bool `json:"can_view"`
	CanCreate bool `json:"can_create"`
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}

type MyMenuItem struct {
	KeyName    string         `json:"key_name"`
	Label      string         `json:"label"`
	Icon       *string        `json:"icon"`
	Path       *string        `json:"path"`
	OrderIndex int            `json:"order_index"`
	Permission MenuPermission `json:"permission"`
	Children   []MyMenuItem   `json:"children"`
}
