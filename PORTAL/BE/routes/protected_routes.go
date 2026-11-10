package routes

import (
	auth_repo "portal_api/domain/auth/repo"
	auth_service "portal_api/domain/auth/service"
	pos_middleware "portal_api/middleware/auth"
	pkgdatabase "portal_api/pkg/database"
	segment "portal_api/routes/segment"

	"github.com/gin-gonic/gin"
)

func protectedRoutes(r *gin.RouterGroup) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authSvc := auth_service.NewAuthService(authRepo)

	r.Use(pos_middleware.POSBearerAuthMiddleware(authSvc))

	segment.AuthRoutes(r)
	segment.UserRoutes(r)
	segment.RoleRoutes(r)
}
