package segment

import (
	finance_handler "pos_api/domain/finance/handler"
	finance_repo "pos_api/domain/finance/repo"
	finance_service "pos_api/domain/finance/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func FinanceRoutes(r *gin.RouterGroup) {
	financeRepo := finance_repo.NewFinanceRepo(pkgdatabase.DB)
	financeService := finance_service.NewFinanceService(financeRepo)
	financeHandler := finance_handler.NewFinanceHandler(financeService)

	g := r.Group("/finance")
	{
		g.POST("/summary", middleware.RoleMiddleware("owner", "admin"), financeHandler.GetSummary)
		g.POST("/cashflow", middleware.RoleMiddleware("owner", "admin"), financeHandler.GetCashflow)
	}
}
