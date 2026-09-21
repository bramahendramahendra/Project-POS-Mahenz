package segment

import (
	capital_reconciliation_handler "pos_api/domain/capital_reconciliation/handler"
	capital_reconciliation_repo "pos_api/domain/capital_reconciliation/repo"
	capital_reconciliation_service "pos_api/domain/capital_reconciliation/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// CapitalReconciliationRoutes — menu Rekonsiliasi Modal (koreksi manual HPP per
// baris penjualan yang tidak bisa dibetulkan otomatis). HANYA untuk admin:
//  1. RoleMiddleware("admin") di level grup (hard lock)
//  2. PermissionMiddleware per aksi (mengikuti role_menu_access; menu di-seed
//     hanya untuk admin di migrasi 010)
func CapitalReconciliationRoutes(r *gin.RouterGroup) {
	repo := capital_reconciliation_repo.NewCapitalReconciliationRepo(pkgdatabase.DB)
	svc := capital_reconciliation_service.NewCapitalReconciliationService(repo)
	handler := capital_reconciliation_handler.NewCapitalReconciliationHandler(svc)

	accessSvc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(accessSvc, "pelaporan.rekonsiliasi_modal", action)
	}

	g := r.Group("/capital-reconciliation", middleware.RoleMiddleware("admin"))
	{
		g.POST("/list", perm("can_view"), handler.GetList)
		g.POST("/summary", perm("can_view"), handler.GetSummary)
		g.POST("/resolve", perm("can_edit"), handler.Resolve)
		g.POST("/skip", perm("can_edit"), handler.Skip)
	}
}
