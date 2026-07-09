package dto

import "time"

type GetAllRequest struct {
	Page      int    `json:"page" validate:"required,min=1"`
	Limit     int    `json:"limit" validate:"required,min=1"`
	Search    string `json:"search" validate:"max=100"`
	RoleID    *int   `json:"role_id"`
	IsActive  *bool  `json:"is_active"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type GetByIDRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type CreateRequest struct {
	Username string `json:"username"  validate:"required,min=3,alphanum"`
	Password string `json:"password"  validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   int    `json:"role_id"   validate:"required,gt=0"`
}

type UpdateUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type UpdateRequest struct {
	ID       int    `json:"-"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   int    `json:"role_id"   validate:"required,gt=0"`
}

type DeleteRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ToggleStatusRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ChangePasswordUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ChangePasswordRequest struct {
	ID       int    `json:"-"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	RoleID    int       `json:"role_id"`
	RoleName  string    `json:"role_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
