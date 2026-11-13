package segment

import (
	sso_handler "auth_api/domain/sso/handler"
	sso_repo "auth_api/domain/sso/repo"
	sso_service "auth_api/domain/sso/service"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func newSSOHandler() *sso_handler.SSOHandler {
	repo := sso_repo.NewSSORepo(pkgdatabase.DB)
	return sso_handler.NewSSOHandler(sso_service.NewSSOService(repo))
}

// SSOLaunchRoutes: butuh login Portal (dipasang di grup protected).
func SSOLaunchRoutes(r *gin.RouterGroup) {
	h := newSSOHandler()

	g := r.Group("/sso")
	{
		g.POST("/launch", h.Launch)
	}
}

// SSOExchangeRoutes: dipanggil BE aplikasi, diamankan kunci aplikasi (dipasang di grup public).
func SSOExchangeRoutes(r *gin.RouterGroup) {
	h := newSSOHandler()

	g := r.Group("/sso")
	{
		g.POST("/exchange", h.Exchange)
	}
}
