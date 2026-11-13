package segment

import (
	finance_handler "coffee_api/domain/finance/handler"
	finance_repo "coffee_api/domain/finance/repo"
	finance_service "coffee_api/domain/finance/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func FinanceRoutes(r *gin.RouterGroup) {
	financeRepo := finance_repo.NewFinanceRepo(pkgdatabase.DB)
	financeService := finance_service.NewFinanceService(financeRepo)
	financeHandler := finance_handler.NewFinanceHandler(financeService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "keuangan.dashboard", action)
	}

	g := r.Group("/finance")
	{
		g.POST("/summary", perm("can_view"), financeHandler.GetSummary)
		g.POST("/cashflow", perm("can_view"), financeHandler.GetCashflow)
	}
}
