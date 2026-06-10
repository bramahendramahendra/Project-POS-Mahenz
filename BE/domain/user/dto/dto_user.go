package dto

import "time"

type GetAllRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	RoleID   *int   `json:"role_id"`
	IsActive *bool  `json:"is_active"`
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
	Password string `json:"password"  validate:"omitempty,min=6"`
}

type DeleteRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type ToggleStatusRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
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
