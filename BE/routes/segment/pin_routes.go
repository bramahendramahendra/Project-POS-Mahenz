package segment

import (
	pin_handler "pos_api/domain/pin/handler"
	pin_repo "pos_api/domain/pin/repo"
	pin_service "pos_api/domain/pin/service"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func PinRoutes(r *gin.RouterGroup) {
	pinRepo := pin_repo.NewPinRepo(pkgdatabase.DB)
	pinService := pin_service.NewPinService(pinRepo)
	pinHandler := pin_handler.NewPinHandler(pinService)

	g := r.Group("/pin")
	{
		g.GET("/check", pinHandler.CheckPin)
		g.POST("/set", pinHandler.SetPin)
		g.POST("/verify", pinHandler.VerifyPin)
		g.POST("/change", pinHandler.ChangePin)
	}
}
