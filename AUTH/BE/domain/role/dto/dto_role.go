package dto

// RoleOptionResponse digunakan untuk dropdown pemilihan role (hanya role aktif)
type RoleOptionResponse struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
}
