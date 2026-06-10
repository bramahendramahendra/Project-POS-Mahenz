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
	roleSvc := role_service.NewRoleService(roleRepo)
	roleHand := role_handler.NewRoleHandler(roleSvc)

	g := r.Group("/roles")
	{
		g.POST("/list",              middleware.RoleMiddleware("owner", "admin"), roleHand.GetAll)
		g.POST("/detail/:id",        middleware.RoleMiddleware("owner", "admin"), roleHand.GetByID)
		g.POST("/create",            middleware.RoleMiddleware("owner"),          roleHand.Create)
		g.POST("/update/:id",        middleware.RoleMiddleware("owner"),          roleHand.Update)
		g.POST("/delete/:id",        middleware.RoleMiddleware("owner"),          roleHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner"),          roleHand.ToggleStatus)
	}
}
