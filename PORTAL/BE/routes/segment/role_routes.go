package segment

import (
	role_handler "portal_api/domain/role/handler"
	role_service "portal_api/domain/role/service"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// RoleRoutes: daftar pilihan role untuk form user (admin, diperiksa BE User).
func RoleRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	roleHandler := role_handler.NewRoleHandler(role_service.NewRoleService(authClient))

	g := r.Group("/roles")
	{
		g.POST("/options", roleHandler.GetOptions)
	}
}
