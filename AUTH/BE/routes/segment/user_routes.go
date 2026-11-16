package segment

import (
	role_repo "auth_api/domain/role/repo"
	user_handler "auth_api/domain/user/handler"
	user_repo "auth_api/domain/user/repo"
	user_service "auth_api/domain/user/service"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// UserRoutes: kelola akun (grup sudah dibatasi Portal + admin).
func UserRoutes(r *gin.RouterGroup) {
	userRepo := user_repo.NewUserRepo(pkgdatabase.DB)
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	userService := user_service.NewUserService(userRepo, roleRepo)
	userHandler := user_handler.NewUserHandler(userService)

	g := r.Group("/users")
	{
		g.POST("/list", userHandler.GetAll)
		g.POST("/create", userHandler.Create)
		g.POST("/update/:id", userHandler.Update)
		g.POST("/reset-password/:id", userHandler.ResetPassword)
		g.POST("/delete/:id", userHandler.Delete)
		g.POST("/toggle-status/:id", userHandler.ToggleStatus)
	}
}
