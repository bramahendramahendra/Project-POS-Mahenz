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

// UserData: data akun dari BE User. ID adalah id akun di BE User (bukan id profil lokal
// Portal), sama dengan id yang dipakai halaman Kelola User.
type UserData struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	RoleID             int    `json:"role_id"`
	RoleName           string `json:"role_name"`
	MustChangePassword bool   `json:"must_change_password"`
}

// VerifyUser: jawaban POST /login/verify dan /account/profile dari BE User.
type VerifyUser struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	RoleID             int    `json:"role_id"`
	RoleName           string `json:"role_name"`
	IsActive           bool   `json:"is_active"`
	MustChangePassword bool   `json:"must_change_password"`
}

// VerifyRequest: body POST /login/verify ke BE User. client_ip = IP pengguna (untuk pembatas login).
type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientIP string `json:"client_ip"`
}

// ChangeMyPasswordRequest: ganti password milik sendiri (wajib menyertakan password lama).
// Divalidasi juga oleh BE User; di sini agar pesan validasi sama seperti sebelumnya.
type ChangeMyPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
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
