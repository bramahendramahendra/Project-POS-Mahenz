package dto

import "time"

type LoginRequest struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	DeviceInfo string `json:"device_info" validate:"required,oneof=desktop web android"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserData  `json:"user"`
}

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// LoginOptionsResponse: pilihan cara login yang tersedia (endpoint publik untuk halaman login).
type LoginOptionsResponse struct {
	// DirectLoginEnabled: true = form username + password boleh dipakai.
	DirectLoginEnabled bool `json:"direct_login_enabled"`
	// PortalLoginEnabled: true = login lewat Portal sudah diatur di server ini.
	PortalLoginEnabled bool `json:"portal_login_enabled"`
	// DirectLoginForced: true = login langsung dipaksa aktif oleh flag darurat di config server.
	DirectLoginForced bool `json:"direct_login_forced"`
}

type SetDirectLoginRequest struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

// SSORequest: kode sekali pakai dari Portal untuk ditukar menjadi sesi Toko.
type SSORequest struct {
	Code       string `json:"code" validate:"required,max=128"`
	DeviceInfo string `json:"device_info" validate:"omitempty,oneof=desktop web android"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyTokenResponse struct {
	Valid       bool           `json:"valid"`
	Claims      map[string]any `json:"claims,omitempty"`
	ExpReadable string         `json:"exp_readable,omitempty"`
	Error       string         `json:"error,omitempty"`
}

type RefreshResponse struct {
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
