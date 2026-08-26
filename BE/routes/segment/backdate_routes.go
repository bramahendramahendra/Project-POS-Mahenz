package segment

import (
	backdate_handler "pos_api/domain/backdate/handler"
	backdate_repo "pos_api/domain/backdate/repo"
	backdate_service "pos_api/domain/backdate/service"
	cb_repo "pos_api/domain/customer_balance/repo"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func BackdateRoutes(r *gin.RouterGroup) {
	backdateRepo := backdate_repo.NewBackdateRepo(pkgdatabase.DB)
	customerBalanceRepo := cb_repo.NewCustomerBalanceRepo(pkgdatabase.DB)
	backdateSvc := backdate_service.NewBackdateService(backdateRepo, customerBalanceRepo)
	handler := backdate_handler.NewBackdateHandler(backdateSvc)

	svc := newAccessService()
	permKas := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "historis.kas_historis", action)
	}
	permKasir := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "historis.kasir_historis", action)
	}

	g := r.Group("/backdate")
	{
		g.POST("/cash-drawer/open", permKas("can_create"), handler.OpenCashDrawer)
		g.POST("/cash-drawer/current", permKas("can_view"), handler.GetCurrentCashDrawer)
		g.POST("/cash-drawer/close/:id", permKas("can_edit"), handler.CloseCashDrawer)
		g.POST("/transactions/create", permKasir("can_create"), handler.CreateTransaction)
	}
}
