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
