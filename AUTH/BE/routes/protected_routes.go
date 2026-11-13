package routes

import (
	"github.com/gin-gonic/gin"
)

// protectedRoutes: API internal yang wajib membawa kunci aplikasi (X-App-Key).
// Diisi di tahap 3B.2: verifikasi login, SSO, akun sendiri, kelola akun/aplikasi/akses
// (rancangan bagian 5.0d dokumen RANCANGAN_PORTAL_TOKO_COFFEE.md).
func protectedRoutes(_ *gin.RouterGroup) {}
