package segment

import (
	app_handler "portal_api/domain/app/handler"
	app_service "portal_api/domain/app/service"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

func AppRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	appHandler := app_handler.NewAppHandler(app_service.NewAppService(authClient))

	g := r.Group("/apps")
	{
		// Semua user login: aplikasi miliknya.
		g.POST("/my", appHandler.GetMy)

		// Khusus admin (diperiksa BE User): kelola aplikasi dan akses user.
		g.POST("/list", appHandler.GetAll)
		g.POST("/create", appHandler.Create)
		g.POST("/update/:id", appHandler.Update)
		g.POST("/toggle-status/:id", appHandler.ToggleStatus)
		g.POST("/rotate-key/:id", appHandler.RotateKey)
		g.POST("/access/:user_id", appHandler.GetUserAccess)
		g.POST("/access/:user_id/set", appHandler.SetUserAccess)
	}
}
