package segment

import (
	setting_handler "pos_api/domain/setting/handler"
	setting_repo "pos_api/domain/setting/repo"
	setting_service "pos_api/domain/setting/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func SettingRoutes(r *gin.RouterGroup) {
	settingRepo := setting_repo.NewSettingRepo(pkgdatabase.DB)
	settingService := setting_service.NewSettingService(settingRepo)
	settingHandler := setting_handler.NewSettingHandler(settingService)

	g := r.Group("/settings")
	{
		g.GET("", settingHandler.GetAll)
		g.GET("/:key", settingHandler.GetByKey)
		g.POST("", middleware.RoleMiddleware("owner", "admin"), settingHandler.Save)
		g.POST("/reset", middleware.RoleMiddleware("admin"), settingHandler.Reset)
	}
}
