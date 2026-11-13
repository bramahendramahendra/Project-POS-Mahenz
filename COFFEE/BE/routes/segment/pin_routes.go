package segment

import (
	pin_handler "coffee_api/domain/pin/handler"
	pin_repo "coffee_api/domain/pin/repo"
	pin_service "coffee_api/domain/pin/service"
	pkgdatabase "coffee_api/pkg/database"

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
