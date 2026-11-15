package segment

import (
	setting_handler "coffee_api/domain/setting/handler"
	setting_repo "coffee_api/domain/setting/repo"
	setting_service "coffee_api/domain/setting/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func SettingRoutes(r *gin.RouterGroup) {
	settingRepo := setting_repo.NewSettingRepo(pkgdatabase.DB)
	settingService := setting_service.NewSettingService(settingRepo)
	settingHandler := setting_handler.NewSettingHandler(settingService)

	svc := newAccessService()

	g := r.Group("/settings")
	{
		// Pengaturan umum & reset tidak punya representasi menu tersendiri di role_menu_access,
		// tetap pakai RoleMiddleware.
		g.GET("", settingHandler.GetAll)
		g.POST("", middleware.RoleMiddleware("owner", "admin"), settingHandler.Save)

		g.GET("/store", settingHandler.GetStoreProfile)
		g.POST("/store", middleware.PermissionMiddleware(svc, "sistem.profil_toko", "can_edit"), settingHandler.UpdateStoreProfile)

		g.GET("/printer", settingHandler.GetPrinterSettings)
		g.POST("/printer", middleware.PermissionMiddleware(svc, "sistem.printer", "can_edit"), settingHandler.UpdatePrinterSettings)

		g.GET("/:key", settingHandler.GetByKey)
	}
}
