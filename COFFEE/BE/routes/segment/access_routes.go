package segment

import (
	access_handler "coffee_api/domain/access/handler"
	access_repo "coffee_api/domain/access/repo"
	access_service "coffee_api/domain/access/service"
	role_repo "coffee_api/domain/role/repo"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AccessRoutes(r *gin.RouterGroup) {
	accessRepo := access_repo.NewAccessRepo(pkgdatabase.DB)
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	accessService := access_service.NewAccessService(accessRepo, roleRepo)
	accessHandler := access_handler.NewAccessHandler(accessService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "sistem.roles", action)
	}

	g := r.Group("/roles/:id/menus")
	{
		g.POST("/list", perm("can_view"), accessHandler.GetByRoleID)
		g.POST("/set", perm("can_edit"), accessHandler.SetRoleAccess)
	}
}
