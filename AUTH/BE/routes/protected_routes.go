package routes

import (
	app_repo "auth_api/domain/app/repo"
	auth_repo "auth_api/domain/auth/repo"
	"auth_api/middleware"
	pkgdatabase "auth_api/pkg/database"
	segment "auth_api/routes/segment"

	"github.com/gin-gonic/gin"
)

// protectedRoutes: API internal, setiap permintaan wajib membawa kunci aplikasi (X-App-Key).
// Rancangan: bagian 5.0d dokumen RANCANGAN_PORTAL_TOKO_COFFEE.md (keputusan 63-73).
func protectedRoutes(r *gin.RouterGroup) {
	apps := appKeyLookup{repo: app_repo.NewAppRepo(pkgdatabase.DB)}
	users := actingUserLookup{repo: auth_repo.NewAuthRepo(pkgdatabase.DB)}

	r.Use(middleware.AppKeyMiddleware(apps))

	// Semua aplikasi: verifikasi login dan penukaran kode SSO.
	segment.LoginRoutes(r)
	segment.SSOExchangeRoutes(r)

	// Khusus Portal, atas nama user yang bertindak (X-Acting-User-Id).
	portal := r.Group("", middleware.PortalOnlyMiddleware())

	self := portal.Group("", middleware.ActingUserMiddleware(users, false))
	segment.AccountRoutes(self)
	segment.MyAppRoutes(self)
	segment.SSOLaunchRoutes(self)

	admin := portal.Group("", middleware.ActingUserMiddleware(users, true))
	segment.UserRoutes(admin)
	segment.RoleRoutes(admin)
	segment.AppAdminRoutes(admin)
}
