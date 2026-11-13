package dto

// LaunchRequest: user (sudah login di Portal) ingin membuka sebuah aplikasi.
// Dikirim BE Portal atas nama user (header X-Acting-User-Id).
type LaunchRequest struct {
	AppCode string `json:"app_code" validate:"required,max=50"`
}

type LaunchResponse struct {
	RedirectURL      string `json:"redirect_url"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

// ExchangeRequest: dikirim BE aplikasi (server ke server) untuk menukar kode SSO.
// Aplikasi pemanggil dikenali dari kunci aplikasi (header X-App-Key), bukan dari body.
type ExchangeRequest struct {
	Code string `json:"code" validate:"required,max=128"`
}

// ExchangeUser: identitas akun di BE User. Aplikasi menautkan profil lokalnya lewat id (auth_user_id).
type ExchangeUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

type ExchangeResponse struct {
	Valid bool         `json:"valid"`
	AppID string       `json:"app_id"`
	User  ExchangeUser `json:"user"`
}
