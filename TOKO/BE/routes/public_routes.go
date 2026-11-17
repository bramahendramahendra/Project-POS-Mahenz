package routes

import (
	auth_handler "pos_api/domain/auth/handler"
	auth_repo "pos_api/domain/auth/repo"
	auth_service "pos_api/domain/auth/service"
	version_handler "pos_api/domain/version/handler"
	version_repo "pos_api/domain/version/repo"
	version_service "pos_api/domain/version/service"
	"pos_api/middleware"
	"pos_api/pkg/authclient"
	pkgdatabase "pos_api/pkg/database"

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

	versionRepoInst := version_repo.NewVersionRepo(pkgdatabase.DB)
	versionSvc := version_service.NewVersionService(versionRepoInst)
	versionHand := version_handler.NewVersionHandler(versionSvc)

	r.GET("/version/android", versionHand.CheckAndroid)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "00",
			"status":  true,
			"message": "OK",
			"data":    nil,
		})
	})
}
