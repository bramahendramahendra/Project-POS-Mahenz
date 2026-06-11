package segment

import (
	auth_handler "pos_api/domain/auth/handler"
	auth_repo "pos_api/domain/auth/repo"
	auth_service "pos_api/domain/auth/service"
	pkgdatabase "pos_api/pkg/database"

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
	}
}
