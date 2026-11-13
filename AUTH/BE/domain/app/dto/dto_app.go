package dto

import "time"

// MyAppResponse: aplikasi yang boleh dibuka oleh user yang sedang login.
type MyAppResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Icon        string `json:"icon"`
}

// ---- Khusus admin ----

type IDUriRequest struct {
	ID int `uri:"id" validate:"required,gt=0"`
}

type CreateRequest struct {
	Code        string `json:"code"        validate:"required,min=2,max=50"`
	Name        string `json:"name"        validate:"required,max=100"`
	Description string `json:"description" validate:"max=255"`
	Path        string `json:"path"        validate:"required,max=255"`
	Icon        string `json:"icon"        validate:"max=50"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateRequest struct {
	ID          int    `json:"-"`
	Name        string `json:"name"        validate:"required,max=100"`
	Description string `json:"description" validate:"max=255"`
	Path        string `json:"path"        validate:"required,max=255"`
	Icon        string `json:"icon"        validate:"max=50"`
	SortOrder   int    `json:"sort_order"`
}

type AppAdminResponse struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Icon        string    `json:"icon"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	HasKey      bool      `json:"has_key"`
	UserCount   int       `json:"user_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAppResponse: kunci aplikasi asli hanya muncul di sini (dan saat rotasi), tidak disimpan di Portal.
type CreateAppResponse struct {
	App    AppAdminResponse `json:"app"`
	AppKey string           `json:"app_key"`
}

type RotateKeyResponse struct {
	AppKey string `json:"app_key"`
}

type UserAccessItem struct {
	AppID    int    `json:"app_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Granted  bool   `json:"granted"`
}

type UserUriRequest struct {
	UserID int `uri:"user_id" validate:"required,gt=0"`
}

type SetUserAccessRequest struct {
	UserID int   `json:"-"`
	AppIDs []int `json:"app_ids" validate:"dive,gt=0"`
}
