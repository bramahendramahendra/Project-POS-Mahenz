package segment

import (
	dashboard_handler "pos_api/domain/dashboard/handler"
	dashboard_repo "pos_api/domain/dashboard/repo"
	dashboard_service "pos_api/domain/dashboard/service"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// DashboardRoutes — landing page operasional (/dashboard), dipakai semua role yang
// punya akses menu beranda.dashboard (termasuk Kasir). Datanya sudah otomatis
// ter-scope ke user yang login (lewat user_id di context), jadi tidak perlu
// permission middleware per-menu seperti domain pelaporan.
func DashboardRoutes(r *gin.RouterGroup) {
	dashboardRepo := dashboard_repo.NewDashboardRepo(pkgdatabase.DB)
	dashboardService := dashboard_service.NewDashboardService(dashboardRepo)
	dashboardHandler := dashboard_handler.NewDashboardHandler(dashboardService)

	g := r.Group("/dashboard")
	{
		g.GET("/recent-transactions", dashboardHandler.GetRecentTransactions)
		g.GET("/today-summary", dashboardHandler.GetTodaySummary)
	}
}
