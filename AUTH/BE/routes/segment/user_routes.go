package segment

import (
	role_repo "auth_api/domain/role/repo"
	user_handler "auth_api/domain/user/handler"
	user_repo "auth_api/domain/user/repo"
	user_service "auth_api/domain/user/service"
	middleware "auth_api/middleware"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// UserRoutes: pengelolaan user Portal. Hanya role admin yang boleh mengakses.
// (Portal tidak memakai sistem permission per menu seperti Toko.)
func UserRoutes(r *gin.RouterGroup) {
	userRepo := user_repo.NewUserRepo(pkgdatabase.DB)
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	userService := user_service.NewUserService(userRepo, roleRepo)
	userHandler := user_handler.NewUserHandler(userService)

	g := r.Group("/users", middleware.RoleMiddleware("admin"))
	{
		g.POST("/list", userHandler.GetAll)
		g.POST("/detail/:id", userHandler.GetByID)
		g.POST("/create", userHandler.Create)
		g.POST("/update/:id", userHandler.Update)
		g.POST("/reset-password/:id", userHandler.ResetPassword)
		g.POST("/delete/:id", userHandler.Delete)
		g.POST("/toggle-status/:id", userHandler.ToggleStatus)
	}
}
