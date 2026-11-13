package segment

import (
	auth_handler "auth_api/domain/auth/handler"
	auth_repo "auth_api/domain/auth/repo"
	auth_service "auth_api/domain/auth/service"
	"auth_api/middleware"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func newAuthHandler() *auth_handler.AuthHandler {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	return auth_handler.NewAuthHandler(auth_service.NewAuthService(authRepo))
}

// LoginRoutes: verifikasi username + password untuk semua aplikasi (Portal, Toko, Coffee).
func LoginRoutes(r *gin.RouterGroup) {
	h := newAuthHandler()

	g := r.Group("/login")
	{
		g.POST("/verify", middleware.LoginRateLimitMiddleware(), h.Verify)
	}
}

// AccountRoutes: akun milik user yang bertindak (Portal atas nama user).
func AccountRoutes(r *gin.RouterGroup) {
	h := newAuthHandler()

	g := r.Group("/account")
	{
		g.POST("/profile", h.Me)
		g.POST("/change-password", h.ChangePassword)
	}
}
