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
	Name        string `json:"name"         validate:"required,min=2,alphanum"`
	DisplayName string `json:"display_name" validate:"required"`
	Description string `json:"description"`
}

type UpdateUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type UpdateRequest struct {
	ID          int    `json:"-"`
	DisplayName string `json:"display_name" validate:"required"`
	Description string `json:"description"`
}

type DeleteRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ToggleStatusRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

// RoleOptionResponse digunakan untuk dropdown pemilihan role (hanya role aktif)
type RoleOptionResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

type RoleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}
