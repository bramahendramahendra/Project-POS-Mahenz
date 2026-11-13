package segment

import (
	auth_handler "coffee_api/domain/auth/handler"
	auth_repo "coffee_api/domain/auth/repo"
	auth_service "coffee_api/domain/auth/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authService := auth_service.NewAuthService(authRepo)
	authHandler := auth_handler.NewAuthHandler(authService)

	g := r.Group("/auth")
	{
		g.POST("/me", authHandler.GetMe)
		g.POST("/logout", authHandler.Logout)
		// Hanya role admin (bukan owner) yang boleh menyalakan/mematikan login langsung.
		g.POST("/direct-login", middleware.RoleMiddleware("admin"), authHandler.SetDirectLogin)
	}
}
