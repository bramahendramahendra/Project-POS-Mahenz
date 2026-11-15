package routes

import (
	auth_handler "coffee_api/domain/auth/handler"
	auth_repo "coffee_api/domain/auth/repo"
	auth_service "coffee_api/domain/auth/service"
	"coffee_api/middleware"
	"coffee_api/pkg/authclient"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func publicRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authSvc := auth_service.NewAuthService(authRepo, authClient)
	authHand := auth_handler.NewAuthHandler(authSvc)

	authGroup := r.Group("/auth")
	authGroup.POST("/login", middleware.LoginRateLimitMiddleware(), authHand.Login)
	// Login lewat Portal: kode sekali pakai (berentropi tinggi), jadi tanpa pembatas percobaan.
	authGroup.POST("/sso", authHand.SSO)
	authGroup.GET("/login-options", authHand.GetLoginOptions)
	authGroup.POST("/refresh", authHand.RefreshToken)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "00",
			"status":  true,
			"message": "OK",
			"data":    nil,
		})
	})
}
