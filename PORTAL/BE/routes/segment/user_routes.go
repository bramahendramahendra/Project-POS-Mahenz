package segment

import (
	user_handler "portal_api/domain/user/handler"
	user_repo "portal_api/domain/user/repo"
	user_service "portal_api/domain/user/service"
	"portal_api/pkg/authclient"
	pkgdatabase "portal_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// UserRoutes: pengelolaan akun, diteruskan ke BE User. Hanya admin (diperiksa BE User).
func UserRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	userRepo := user_repo.NewUserRepo(pkgdatabase.DB)
	userService := user_service.NewUserService(userRepo, authClient)
	userHandler := user_handler.NewUserHandler(userService)

	g := r.Group("/users")
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
