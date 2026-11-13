package dto

// VerifyRequest: BE aplikasi meneruskan username + password dari form login-nya (keputusan 63).
// client_ip = IP browser pengguna, dipakai pembatas login (koneksi ke BE User selalu dari server).
type VerifyRequest struct {
	Username string `json:"username"  validate:"required,max=50"`
	Password string `json:"password"  validate:"required,max=72"`
	ClientIP string `json:"client_ip" validate:"max=45"`
}

// UserData: identitas akun. Aplikasi menautkan profil lokalnya lewat id (auth_user_id).
// role_name = role pengelola akun (admin/user), bukan role di aplikasi.
type UserData struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	RoleID             int    `json:"role_id"`
	RoleName           string `json:"role_name"`
	IsActive           bool   `json:"is_active"`
	MustChangePassword bool   `json:"must_change_password"`
}

// ChangeMyPasswordRequest: ganti password milik sendiri (wajib menyertakan password lama).
type ChangeMyPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}
