package dto

// LaunchRequest: user (sudah login di Portal) ingin membuka sebuah aplikasi.
type LaunchRequest struct {
	AppCode string `json:"app_code" validate:"required,max=50"`
}

type LaunchResponse struct {
	RedirectURL      string `json:"redirect_url"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

// ExchangeRequest: dikirim BE aplikasi (server ke server) untuk menukar kode SSO.
// Kunci aplikasi dikirim lewat header X-App-Key, bukan di body.
type ExchangeRequest struct {
	AppID string `json:"app_id" validate:"required,max=50"`
	Code  string `json:"code"   validate:"required,max=128"`
}

type ExchangeUser struct {
	PortalUserID int    `json:"portal_user_id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	IsActive     bool   `json:"is_active"`
}

type ExchangeResponse struct {
	Valid bool         `json:"valid"`
	AppID string       `json:"app_id"`
	User  ExchangeUser `json:"user"`
}
