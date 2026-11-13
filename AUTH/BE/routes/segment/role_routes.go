package segment

import (
	role_handler "auth_api/domain/role/handler"
	role_repo "auth_api/domain/role/repo"
	role_service "auth_api/domain/role/service"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// RoleRoutes: pilihan role pengelola akun (grup sudah dibatasi Portal + admin).
func RoleRoutes(r *gin.RouterGroup) {
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	roleService := role_service.NewRoleService(roleRepo)
	roleHandler := role_handler.NewRoleHandler(roleService)

	g := r.Group("/roles")
	{
		g.POST("/options", roleHandler.GetOptions)
	}
}
