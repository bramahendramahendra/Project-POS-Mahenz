package segment

import (
	auth_handler "pos_api/domain/auth/handler"
	auth_repo "pos_api/domain/auth/repo"
	auth_service "pos_api/domain/auth/service"
	middleware "pos_api/middleware"
	"pos_api/pkg/authclient"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authService := auth_service.NewAuthService(authRepo, authClient)
	authHandler := auth_handler.NewAuthHandler(authService)

	g := r.Group("/auth")
	{
		g.POST("/logout", authHandler.Logout)
		// Hanya role admin (bukan owner) yang boleh menyalakan/mematikan login langsung.
		g.POST("/direct-login", middleware.RoleMiddleware("admin"), authHandler.SetDirectLogin)
	}
}
