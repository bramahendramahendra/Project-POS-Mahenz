package segment

import (
	business_summary_handler "pos_api/domain/business_summary/handler"
	business_summary_repo "pos_api/domain/business_summary/repo"
	business_summary_service "pos_api/domain/business_summary/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func BusinessSummaryRoutes(r *gin.RouterGroup) {
	businessSummaryRepo := business_summary_repo.NewBusinessSummaryRepo(pkgdatabase.DB)
	businessSummaryService := business_summary_service.NewBusinessSummaryService(businessSummaryRepo)
	businessSummaryHandler := business_summary_handler.NewBusinessSummaryHandler(businessSummaryService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.ringkasan_bisnis", action)
	}

	g := r.Group("/reports/business-summary")
	{
		g.GET("/stats", perm("can_view"), businessSummaryHandler.GetStats)
		g.GET("/sales-trend", perm("can_view"), businessSummaryHandler.GetSalesTrend)
		g.GET("/top-products", perm("can_view"), businessSummaryHandler.GetTopProducts)
		g.GET("/top-categories", perm("can_view"), businessSummaryHandler.GetTopCategories)
		g.GET("/payment-methods", perm("can_view"), businessSummaryHandler.GetPaymentMethods)
		g.GET("/summary-extra", perm("can_view"), businessSummaryHandler.GetSummaryExtra)
	}
}
