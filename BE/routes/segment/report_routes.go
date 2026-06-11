package segment

import (
	report_handler "pos_api/domain/report/handler"
	report_repo "pos_api/domain/report/repo"
	report_service "pos_api/domain/report/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(r *gin.RouterGroup) {
	reportRepo := report_repo.NewReportRepo(pkgdatabase.DB)
	reportService := report_service.NewReportService(reportRepo)
	reportHandler := report_handler.NewReportHandler(reportService)

	g := r.Group("/reports")
	{
		g.GET("/sales", reportHandler.GetSalesReport)
		g.GET("/sales/chart", reportHandler.GetSalesChart)
		g.GET("/sales/export", middleware.RoleMiddleware("owner", "admin"), reportHandler.ExportSalesReport)
		g.GET("/profit-loss", middleware.RoleMiddleware("owner", "admin"), reportHandler.GetProfitLoss)
		g.GET("/profit-loss/export", middleware.RoleMiddleware("owner", "admin"), reportHandler.ExportProfitLoss)
		g.GET("/stock", reportHandler.GetStockReport)
		g.GET("/stock/export", middleware.RoleMiddleware("owner", "admin"), reportHandler.ExportStockReport)
		g.GET("/cashier", middleware.RoleMiddleware("owner", "admin"), reportHandler.GetCashierReport)
		g.GET("/cashier/export", middleware.RoleMiddleware("owner", "admin"), reportHandler.ExportCashierReport)
	}
}
