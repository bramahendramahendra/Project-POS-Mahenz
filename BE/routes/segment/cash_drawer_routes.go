package segment

import (
	cash_drawer_handler "pos_api/domain/cash_drawer/handler"
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func CashDrawerRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	cashDrawerService := cash_drawer_service.NewCashDrawerService(cashDrawerRepo)
	cashDrawerHandler := cash_drawer_handler.NewCashDrawerHandler(cashDrawerService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "keuangan.kas_harian", action)
	}
	// Aksi kas milik-sendiri (open/close/update-sales/update-expenses) dipakai dari
	// halaman "Kas Saya" — digembok ke menu itu, bukan "Kas Harian" (rekap semua
	// kasir, khusus admin/owner). Kepemilikan tetap divalidasi di service layer
	// (lihat cash_drawer_service.go: Close/UpdateSales/UpdateExpenses menolak user
	// lain kecuali admin/owner).
	myCashPerm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "keuangan.kas_saya", action)
	}

	g := r.Group("/cash-drawer")
	{
		g.POST("/current", cashDrawerHandler.GetCurrent)
		g.POST("/my-cash", cashDrawerHandler.GetMyCash)
		g.POST("/list", perm("can_view"), cashDrawerHandler.GetHistory)
		g.POST("/summary", perm("can_view"), cashDrawerHandler.GetSummary)
		g.POST("/kasir-options", perm("can_view"), cashDrawerHandler.GetKasirOptions)
		g.POST("/detail/:id", cashDrawerHandler.GetByID)
		g.POST("/open", myCashPerm("can_create"), cashDrawerHandler.Open)
		g.POST("/close/:id", myCashPerm("can_edit"), cashDrawerHandler.Close)
		g.POST("/update-sales/:id", myCashPerm("can_edit"), cashDrawerHandler.UpdateSales)
		g.POST("/update-expenses/:id", myCashPerm("can_edit"), cashDrawerHandler.UpdateExpenses)
	}
}
