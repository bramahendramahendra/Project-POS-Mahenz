package segment

import (
	app_handler "auth_api/domain/app/handler"
	app_repo "auth_api/domain/app/repo"
	app_service "auth_api/domain/app/service"
	pkgdatabase "auth_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func newAppHandler() *app_handler.AppHandler {
	appRepo := app_repo.NewAppRepo(pkgdatabase.DB)
	return app_handler.NewAppHandler(app_service.NewAppService(appRepo))
}

// MyAppRoutes: aplikasi milik user yang bertindak (kartu di halaman Aplikasi Portal).
func MyAppRoutes(r *gin.RouterGroup) {
	h := newAppHandler()

	g := r.Group("/apps")
	{
		g.POST("/my", h.GetMy)
	}
}

// AppAdminRoutes: kelola aplikasi dan akses user (grup sudah dibatasi Portal + admin).
func AppAdminRoutes(r *gin.RouterGroup) {
	h := newAppHandler()

	g := r.Group("/apps")
	{
		g.POST("/list", h.GetAll)
		g.POST("/create", h.Create)
		g.POST("/update/:id", h.Update)
		g.POST("/toggle-status/:id", h.ToggleStatus)
		g.POST("/rotate-key/:id", h.RotateKey)
		g.POST("/access/:user_id", h.GetUserAccess)
		g.POST("/access/:user_id/set", h.SetUserAccess)
	}
}
