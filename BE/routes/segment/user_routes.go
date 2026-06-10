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
	userSvc := user_service.NewUserService(userRepo)
	userHand := user_handler.NewUserHandler(userSvc)

	g := r.Group("/users", middleware.RoleMiddleware("owner", "admin"))
	{
		g.POST("/list",              userHand.GetAll)
		g.POST("/detail/:id",        userHand.GetByID)
		g.POST("/create",            userHand.Create)
		g.POST("/update/:id",        userHand.Update)
		g.POST("/delete/:id",        userHand.Delete)
		g.POST("/toggle-status/:id", userHand.ToggleStatus)
	}
}
