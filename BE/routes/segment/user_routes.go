package segment

import (
	user_handler "pos_api/domain/user/handler"
	user_repo "pos_api/domain/user/repo"
	user_service "pos_api/domain/user/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	userRepo := user_repo.NewUserRepo(pkgdatabase.DB)
	userService := user_service.NewUserService(userRepo)
	userHandler := user_handler.NewUserHandler(userService)

	g := r.Group("/users", middleware.RoleMiddleware("owner", "admin"))
	{
		g.POST("/list", userHandler.GetAll)
		g.POST("/detail/:id", userHandler.GetByID)
		g.POST("/create", userHandler.Create)
		g.POST("/update/:id", userHandler.Update)
		g.POST("/delete/:id", userHandler.Delete)
		g.POST("/toggle-status/:id", userHandler.ToggleStatus)
	}
}
