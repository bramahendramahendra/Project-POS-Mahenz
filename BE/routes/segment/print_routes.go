package segment

import (
	print_handler "pos_api/domain/print/handler"
	print_service "pos_api/domain/print/service"

	"github.com/gin-gonic/gin"
)

func PrintRoutes(r *gin.RouterGroup) {
	printService := print_service.NewPrintService()
	printHandler := print_handler.NewPrintHandler(printService)

	g := r.Group("/print")
	{
		g.POST("/receipt", printHandler.PrintReceipt)
	}
}
