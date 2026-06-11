package segment

import (
	dashboard_handler "pos_api/domain/dashboard/handler"
	dashboard_repo "pos_api/domain/dashboard/repo"
	dashboard_service "pos_api/domain/dashboard/service"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func DashboardRoutes(r *gin.RouterGroup) {
	dashboardRepo := dashboard_repo.NewDashboardRepo(pkgdatabase.DB)
	dashboardService := dashboard_service.NewDashboardService(dashboardRepo)
	dashboardHandler := dashboard_handler.NewDashboardHandler(dashboardService)

	g := r.Group("/dashboard")
	{
		g.GET("/stats", dashboardHandler.GetStats)
		g.GET("/sales-trend", dashboardHandler.GetSalesTrend)
		g.GET("/top-products", dashboardHandler.GetTopProducts)
		g.GET("/top-categories", dashboardHandler.GetTopCategories)
		g.GET("/payment-methods", dashboardHandler.GetPaymentMethods)
		g.GET("/summary-extra", dashboardHandler.GetSummaryExtra)
	}
}
