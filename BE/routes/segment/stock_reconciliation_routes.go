package segment

import (
	stock_reconciliation_handler "pos_api/domain/stock_reconciliation/handler"
	stock_reconciliation_repo "pos_api/domain/stock_reconciliation/repo"
	stock_reconciliation_service "pos_api/domain/stock_reconciliation/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// StockReconciliationRoutes — menu Rekonsiliasi Stok (stok lama vs baru +
// koreksi manual). HANYA untuk admin: dijaga dua lapis —
//  1. RoleMiddleware("admin") di level grup (hard lock, tidak bergantung seed)
//  2. PermissionMiddleware per aksi (mengikuti role_menu_access, konsisten
//     dengan segment lain; menu hanya di-seed untuk admin di migrasi 009)
func StockReconciliationRoutes(r *gin.RouterGroup) {
	repo := stock_reconciliation_repo.NewStockReconciliationRepo(pkgdatabase.DB)
	svc := stock_reconciliation_service.NewStockReconciliationService(repo)
	handler := stock_reconciliation_handler.NewStockReconciliationHandler(svc)

	accessSvc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(accessSvc, "pelaporan.rekonsiliasi_stok", action)
	}

	g := r.Group("/stock-reconciliation", middleware.RoleMiddleware("admin"))
	{
		g.POST("/list", perm("can_view"), handler.GetList)
		g.POST("/summary", perm("can_view"), handler.GetSummary)
		g.POST("/detail/:product_id", perm("can_view"), handler.GetDetail)
		g.POST("/adjust", perm("can_edit"), handler.Adjust)
		g.POST("/mark-reviewed", perm("can_edit"), handler.MarkReviewed)
	}
}
