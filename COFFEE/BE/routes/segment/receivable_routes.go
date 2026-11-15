package segment

import (
	receivable_handler "coffee_api/domain/receivable/handler"
	receivable_repo "coffee_api/domain/receivable/repo"
	receivable_service "coffee_api/domain/receivable/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ReceivableRoutes(r *gin.RouterGroup) {
	receivableRepo := receivable_repo.NewReceivableRepo(pkgdatabase.DB)
	receivableService := receivable_service.NewReceivableService(receivableRepo)
	receivableHandler := receivable_handler.NewReceivableHandler(receivableService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelanggan.piutang", action)
	}

	g := r.Group("/receivables")
	{
		g.POST("/list", receivableHandler.GetAll)
		g.POST("/stats", perm("can_view"), receivableHandler.GetStats)
		g.POST("/detail/:id", receivableHandler.GetByID)
		g.POST("/pay/:id", perm("can_edit"), receivableHandler.Pay)
	}
}
