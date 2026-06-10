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
	receivableSvc := receivable_service.NewReceivableService(receivableRepo)
	receivableHand := receivable_handler.NewReceivableHandler(receivableSvc)

	g := r.Group("/receivables")
	{
		g.POST("/list", receivableHand.GetAll)
		g.POST("/summary", middleware.RoleMiddleware("owner", "admin"), receivableHand.GetSummary)
		g.POST("/detail/:id", receivableHand.GetByID)
		g.POST("/payments/:id", receivableHand.GetPayments)
		g.POST("/pay/:id", receivableHand.Pay)
	}
}
