package segment

import (
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
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
	transactionService := transaction_service.NewTransactionService(transactionRepo, cashDrawerRepo)
	transactionHandler := transaction_handler.NewTransactionHandler(transactionService)

	g := r.Group("/transactions")
	{
		g.POST("/list", transactionHandler.GetAll)
		g.POST("/detail/:id", transactionHandler.GetByID)
		g.POST("/create", transactionHandler.Create)
		g.POST("/void/:id", middleware.RoleMiddleware("owner", "admin"), transactionHandler.Void)
	}
}
