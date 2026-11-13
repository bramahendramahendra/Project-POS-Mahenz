package segment

import (
	auth_handler "portal_api/domain/auth/handler"
	auth_repo "portal_api/domain/auth/repo"
	auth_service "portal_api/domain/auth/service"
	"portal_api/pkg/authclient"
	pkgdatabase "portal_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authService := auth_service.NewAuthService(authRepo, authClient)
	authHandler := auth_handler.NewAuthHandler(authService)

	g := r.Group("/auth")
	{
		g.POST("/me", authHandler.GetMe)
		g.POST("/logout", authHandler.Logout)
		g.POST("/change-password", authHandler.ChangePassword)
	}
}
