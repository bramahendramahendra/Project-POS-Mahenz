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
	// DirectLoginEnabled: true = form username + password boleh dipakai (diperiksa BE User).
	DirectLoginEnabled bool `json:"direct_login_enabled"`
	// PortalLoginEnabled: true = layanan akun (BE User) sudah diatur di server ini, sehingga
	// login lewat Portal (SSO) bisa dipakai.
	PortalLoginEnabled bool `json:"portal_login_enabled"`
}

type SetDirectLoginRequest struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

// SSORequest: kode sekali pakai (dibuat BE User saat kartu Coffee dibuka di Portal) untuk
// ditukar menjadi sesi Coffee.
type SSORequest struct {
	Code       string `json:"code" validate:"required,max=128"`
	DeviceInfo string `json:"device_info" validate:"omitempty,oneof=desktop web android"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	Token        string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// VerifyRequest: body POST /login/verify ke BE User. client_ip = IP pengguna (pembatas login).
type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientIP string `json:"client_ip"`
}

// AccountData: akun dari BE User (jawaban /login/verify dan bagian user dari /sso/exchange).
type AccountData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

// ExchangeRequest: body POST /sso/exchange ke BE User.
type ExchangeRequest struct {
	Code string `json:"code"`
}

// ExchangeResult: jawaban /sso/exchange dari BE User.
type ExchangeResult struct {
	Valid bool        `json:"valid"`
	AppID string      `json:"app_id"` // kode aplikasi, mis. "coffee"
	User  AccountData `json:"user"`
}
