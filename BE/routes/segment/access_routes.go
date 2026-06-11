package segment

import (
	access_handler "pos_api/domain/access/handler"
	access_repo "pos_api/domain/access/repo"
	access_service "pos_api/domain/access/service"
	role_repo "pos_api/domain/role/repo"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AccessRoutes(r *gin.RouterGroup) {
	accessRepo := access_repo.NewAccessRepo(pkgdatabase.DB)
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	accessService := access_service.NewAccessService(accessRepo, roleRepo)
	accessHandler := access_handler.NewAccessHandler(accessService)

	g := r.Group("/roles/:id/menus")
	{
		g.POST("/list", middleware.RoleMiddleware("owner", "admin"), accessHandler.GetByRoleID)
		g.POST("/set", middleware.RoleMiddleware("owner"), accessHandler.SetRoleAccess)
	}
}
