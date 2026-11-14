package segment

import (
	cash_drawer_repo "coffee_api/domain/cash_drawer/repo"
	product_repo "coffee_api/domain/product/repo"
	transaction_handler "coffee_api/domain/transaction/handler"
	transaction_repo "coffee_api/domain/transaction/repo"
	transaction_service "coffee_api/domain/transaction/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	transactionRepo := transaction_repo.NewTransactionRepo(pkgdatabase.DB)
	productRepo := product_repo.NewProductRepo(pkgdatabase.DB)
	transactionService := transaction_service.NewTransactionService(transactionRepo, cashDrawerRepo, productRepo)
	transactionHandler := transaction_handler.NewTransactionHandler(transactionService)

	svc := newAccessService()
	permTransaksi := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "penjualan.transaksi", action)
	}
	permKasir := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "penjualan.kasir", action)
	}

	g := r.Group("/transactions")
	{
		g.POST("/list", permTransaksi("can_view"), transactionHandler.GetAll)
		g.POST("/detail/:id", transactionHandler.GetByID)
		g.POST("/create", permKasir("can_create"), transactionHandler.Create)
		g.POST("/void/:id", permTransaksi("can_delete"), transactionHandler.Void)

		// Open bill (tahap 4.8): kasir, admin, owner (izin tambah di halaman Kasir, keputusan 105, 106).
		o := g.Group("/open")
		o.POST("/list", permKasir("can_view"), transactionHandler.ListOpen)
		o.POST("/:id/add-items", permKasir("can_create"), transactionHandler.AddOpenItems)
		o.POST("/:id/remove-item", permKasir("can_create"), transactionHandler.RemoveOpenItem)
		o.POST("/:id/move-table", permKasir("can_create"), transactionHandler.MoveOpenTable)
		o.POST("/:id/merge", permKasir("can_create"), transactionHandler.MergeOpen)
		o.POST("/:id/cancel", permKasir("can_create"), transactionHandler.CancelOpen)
		o.POST("/:id/pay", permKasir("can_create"), transactionHandler.PayOpen)
	}
}
