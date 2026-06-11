package segment

import (
	receivable_handler "pos_api/domain/receivable/handler"
	receivable_repo "pos_api/domain/receivable/repo"
	receivable_service "pos_api/domain/receivable/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ReceivableRoutes(r *gin.RouterGroup) {
	receivableRepo := receivable_repo.NewReceivableRepo(pkgdatabase.DB)
	receivableService := receivable_service.NewReceivableService(receivableRepo)
	receivableHandler := receivable_handler.NewReceivableHandler(receivableService)

	g := r.Group("/receivables")
	{
		g.POST("/list", receivableHandler.GetAll)
		g.POST("/summary", middleware.RoleMiddleware("owner", "admin"), receivableHandler.GetSummary)
		g.POST("/detail/:id", receivableHandler.GetByID)
		g.POST("/payments/:id", receivableHandler.GetPayments)
		g.POST("/pay/:id", receivableHandler.Pay)
	}
}
