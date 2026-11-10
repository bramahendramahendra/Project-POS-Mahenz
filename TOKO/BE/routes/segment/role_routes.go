package segment

import (
	role_handler "pos_api/domain/role/handler"
	role_repo "pos_api/domain/role/repo"
	role_service "pos_api/domain/role/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.RouterGroup) {
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	roleService := role_service.NewRoleService(roleRepo)
	roleHandler := role_handler.NewRoleHandler(roleService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "sistem.roles", action)
	}

	g := r.Group("/roles")
	{
		g.POST("/list", perm("can_view"), roleHandler.GetAll)
		g.POST("/options", roleHandler.GetOptions)
		g.POST("/detail/:id", perm("can_view"), roleHandler.GetByID)
		g.POST("/create", perm("can_create"), roleHandler.Create)
		g.POST("/update/:id", perm("can_edit"), roleHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), roleHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), roleHandler.ToggleStatus)
	}
}
