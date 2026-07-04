package segment

import (
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	product_repo "pos_api/domain/product/repo"
	transaction_handler "pos_api/domain/transaction/handler"
	transaction_repo "pos_api/domain/transaction/repo"
	transaction_service "pos_api/domain/transaction/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	transactionRepo := transaction_repo.NewTransactionRepo(pkgdatabase.DB)
	productRepo := product_repo.NewProductRepo(pkgdatabase.DB)
	transactionService := transaction_service.NewTransactionService(transactionRepo, cashDrawerRepo, productRepo)
	transactionHandler := transaction_handler.NewTransactionHandler(transactionService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "penjualan.transaksi", action)
	}

	g := r.Group("/transactions")
	{
		g.POST("/list", transactionHandler.GetAll)
		g.POST("/detail/:id", transactionHandler.GetByID)
		g.POST("/create", perm("can_create"), transactionHandler.Create)
		g.POST("/void/:id", perm("can_delete"), transactionHandler.Void)
	}
}
