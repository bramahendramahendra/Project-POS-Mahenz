package segment

import (
	role_handler "auth_api/domain/role/handler"
	role_repo "auth_api/domain/role/repo"
	role_service "auth_api/domain/role/service"
	middleware "auth_api/middleware"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// RoleRoutes: Portal hanya punya 2 role tetap (admin, user), jadi tidak ada CRUD role.
// Endpoint ini hanya menyediakan daftar pilihan role untuk form user (khusus admin).
func RoleRoutes(r *gin.RouterGroup) {
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	roleService := role_service.NewRoleService(roleRepo)
	roleHandler := role_handler.NewRoleHandler(roleService)

	g := r.Group("/roles", middleware.RoleMiddleware("admin"))
	{
		g.POST("/options", roleHandler.GetOptions)
	}
}
