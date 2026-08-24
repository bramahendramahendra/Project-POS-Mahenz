package segment

import (
	customer_handler "pos_api/domain/customer/handler"
	customer_repo "pos_api/domain/customer/repo"
	customer_service "pos_api/domain/customer/service"
	cb_handler "pos_api/domain/customer_balance/handler"
	cb_repo "pos_api/domain/customer_balance/repo"
	cb_service "pos_api/domain/customer_balance/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(r *gin.RouterGroup) {
	customerRepo := customer_repo.NewCustomerRepo(pkgdatabase.DB)
	customerService := customer_service.NewCustomerService(customerRepo)
	customerHandler := customer_handler.NewCustomerHandler(customerService)

	// Customer Balance
	balanceRepo := cb_repo.NewCustomerBalanceRepo(pkgdatabase.DB)
	balanceService := cb_service.NewCustomerBalanceService(balanceRepo, customerRepo)
	balanceHandler := cb_handler.NewCustomerBalanceHandler(balanceService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelanggan.pelanggan", action)
	}

	g := r.Group("/customers")
	{
		g.POST("/list", customerHandler.GetAll)
		g.POST("/active", customerHandler.GetOptions)
		g.POST("/detail/:id", customerHandler.GetByID)
		g.POST("/create", perm("can_create"), customerHandler.Create)
		g.POST("/update/:id", perm("can_edit"), customerHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), customerHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), customerHandler.ToggleStatus)

		// Balance (Saldo Pelanggan)
		g.POST("/:id/balance/topup", perm("can_edit"), balanceHandler.Topup)
		g.POST("/:id/balance/refund", perm("can_edit"), balanceHandler.Refund)
		g.POST("/:id/balance/history", balanceHandler.History)
	}

	// Save-to-balance (dari halaman struk, accessible by all authenticated users)
	r.POST("/transactions/save-to-balance", balanceHandler.SaveToBalance)
}
