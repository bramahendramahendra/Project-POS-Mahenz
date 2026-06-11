package segment

import (
	transaction_handler "pos_api/domain/transaction/handler"
	transaction_repo "pos_api/domain/transaction/repo"
	transaction_service "pos_api/domain/transaction/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup) {
	transactionRepo := transaction_repo.NewTransactionRepo(pkgdatabase.DB)
	transactionService := transaction_service.NewTransactionService(transactionRepo)
	transactionHandler := transaction_handler.NewTransactionHandler(transactionService)

	g := r.Group("/transactions")
	{
		g.GET("", transactionHandler.GetAll)
		g.GET("/:id", transactionHandler.GetByID)
		g.POST("", transactionHandler.Create)
		g.PATCH("/:id/void", middleware.RoleMiddleware("owner", "admin"), transactionHandler.Void)
	}
}
