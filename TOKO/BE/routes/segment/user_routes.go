package segment

import (
	role_repo "pos_api/domain/role/repo"
	user_handler "pos_api/domain/user/handler"
	user_repo "pos_api/domain/user/repo"
	user_service "pos_api/domain/user/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	userRepo := user_repo.NewUserRepo(pkgdatabase.DB)
	roleRepo := role_repo.NewRoleRepo(pkgdatabase.DB)
	userService := user_service.NewUserService(userRepo, roleRepo)
	userHandler := user_handler.NewUserHandler(userService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "sistem.users", action)
	}

	g := r.Group("/users")
	{
		g.POST("/list", perm("can_view"), userHandler.GetAll)
		g.POST("/detail/:id", perm("can_view"), userHandler.GetByID)
		g.POST("/create", perm("can_create"), userHandler.Create)
		g.POST("/update/:id", perm("can_edit"), userHandler.Update)
		g.POST("/change-password/:id", perm("can_edit"), userHandler.ChangePassword)
		g.POST("/delete/:id", perm("can_delete"), userHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), userHandler.ToggleStatus)
	}
}
