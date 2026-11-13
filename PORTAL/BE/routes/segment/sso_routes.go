package segment

import (
	sso_handler "portal_api/domain/sso/handler"
	sso_service "portal_api/domain/sso/service"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// SSOLaunchRoutes: butuh login Portal (dipasang di grup protected). Penukaran kode
// (exchange) tidak lagi di Portal: aplikasi tujuan menukarnya langsung ke BE User.
func SSOLaunchRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	h := sso_handler.NewSSOHandler(sso_service.NewSSOService(authClient))

	g := r.Group("/sso")
	{
		g.POST("/launch", h.Launch)
	}
}
