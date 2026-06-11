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

	g := r.Group("/roles")
	{
		g.POST("/list", middleware.RoleMiddleware("owner", "admin"), roleHandler.GetAll)
		g.POST("/detail/:id", middleware.RoleMiddleware("owner", "admin"), roleHandler.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner"), roleHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner"), roleHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), roleHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner"), roleHandler.ToggleStatus)
	}
}
