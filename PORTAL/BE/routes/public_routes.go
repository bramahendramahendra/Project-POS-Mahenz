package routes

import (
	auth_handler "portal_api/domain/auth/handler"
	auth_repo "portal_api/domain/auth/repo"
	auth_service "portal_api/domain/auth/service"
	"portal_api/middleware"
	pkgdatabase "portal_api/pkg/database"
	segment "portal_api/routes/segment"

	"github.com/gin-gonic/gin"
)

func publicRoutes(r *gin.RouterGroup) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authSvc := auth_service.NewAuthService(authRepo)
	authHand := auth_handler.NewAuthHandler(authSvc)

	authGroup := r.Group("/auth")
	authGroup.POST("/login", middleware.LoginRateLimitMiddleware(), authHand.Login)
	authGroup.POST("/refresh", authHand.RefreshToken)
	authGroup.POST("/verify-token", authHand.VerifyToken)

	// Dipanggil BE aplikasi lain (server ke server), diamankan kunci aplikasi.
	segment.SSOExchangeRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "00",
			"status":  true,
			"message": "OK",
			"data":    nil,
		})
	})
}
