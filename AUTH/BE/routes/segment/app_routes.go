package segment

import (
	app_handler "auth_api/domain/app/handler"
	app_repo "auth_api/domain/app/repo"
	app_service "auth_api/domain/app/service"
	middleware "auth_api/middleware"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func AppRoutes(r *gin.RouterGroup) {
	appRepo := app_repo.NewAppRepo(pkgdatabase.DB)
	appService := app_service.NewAppService(appRepo)
	appHandler := app_handler.NewAppHandler(appService)

	g := r.Group("/apps")
	{
		// Semua user login: aplikasi miliknya.
		g.POST("/my", appHandler.GetMy)

		// Khusus admin: kelola aplikasi dan akses user.
		admin := g.Group("", middleware.RoleMiddleware("admin"))
		{
			admin.POST("/list", appHandler.GetAll)
			admin.POST("/create", appHandler.Create)
			admin.POST("/update/:id", appHandler.Update)
			admin.POST("/toggle-status/:id", appHandler.ToggleStatus)
			admin.POST("/rotate-key/:id", appHandler.RotateKey)
			admin.POST("/access/:user_id", appHandler.GetUserAccess)
			admin.POST("/access/:user_id/set", appHandler.SetUserAccess)
		}
	}
}
