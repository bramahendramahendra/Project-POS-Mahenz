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

	svc := newAccessService()
	permSales := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.penjualan", action)
	}
	permPL := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.laba_rugi", action)
	}
	permStock := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.stok", action)
	}
	permCashier := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.kinerja_kasir", action)
	}

	g := r.Group("/reports")
	{
		g.GET("/sales", permSales("can_view"), reportHandler.GetSalesReport)
		g.GET("/sales/chart", permSales("can_view"), reportHandler.GetSalesChart)
		g.GET("/sales/export", permSales("can_view"), reportHandler.ExportSalesReport)
		g.POST("/sales/list", permSales("can_view"), reportHandler.GetSalesList)
		g.POST("/sales/summary", permSales("can_view"), reportHandler.GetSalesSummaryData)

		g.GET("/profit-loss", permPL("can_view"), reportHandler.GetProfitLoss)
		g.POST("/profit-loss/data", permPL("can_view"), reportHandler.GetProfitLossData)
		g.GET("/profit-loss/export", permPL("can_view"), reportHandler.ExportProfitLoss)

		g.GET("/stock", permStock("can_view"), reportHandler.GetStockReport)
		g.POST("/stock/list", permStock("can_view"), reportHandler.GetStockList)
		g.POST("/stock/summary", permStock("can_view"), reportHandler.GetStockSummaryData)
		g.GET("/stock/export", permStock("can_view"), reportHandler.ExportStockReport)

		g.GET("/cashier", permCashier("can_view"), reportHandler.GetCashierReport)
		g.POST("/cashier/list", permCashier("can_view"), reportHandler.GetCashierList)
		g.GET("/cashier/export", permCashier("can_view"), reportHandler.ExportCashierReport)
	}
}
