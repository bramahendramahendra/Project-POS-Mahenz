package routes

import (
	auth_repo "portal_api/domain/auth/repo"
	auth_service "portal_api/domain/auth/service"
	pos_middleware "portal_api/middleware/auth"
	"portal_api/pkg/authclient"
	pkgdatabase "portal_api/pkg/database"
	segment "portal_api/routes/segment"

	"github.com/gin-gonic/gin"
)

// protectedRoutes: butuh sesi Portal. Hak akses (admin atau bukan) diperiksa BE User
// pada setiap permintaan yang diteruskan, memakai role terbaru akun tersebut.
func protectedRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authSvc := auth_service.NewAuthService(authRepo, authClient)

	r.Use(pos_middleware.POSBearerAuthMiddleware(authSvc))

	segment.AuthRoutes(r, authClient)
	segment.UserRoutes(r, authClient)
	segment.RoleRoutes(r, authClient)
	segment.AppRoutes(r, authClient)
	segment.SSOLaunchRoutes(r, authClient)
}
