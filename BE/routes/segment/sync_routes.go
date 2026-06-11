package segment

import (
	expense_repo "pos_api/domain/expense/repo"
	sync_handler "pos_api/domain/sync/handler"
	sync_repo "pos_api/domain/sync/repo"
	sync_service "pos_api/domain/sync/service"
	transaction_repo "pos_api/domain/transaction/repo"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func SyncRoutes(r *gin.RouterGroup) {
	syncRepo := sync_repo.NewSyncRepo(pkgdatabase.DB)
	transactionRepo := transaction_repo.NewTransactionRepo(pkgdatabase.DB)
	expenseRepo := expense_repo.NewExpenseRepo(pkgdatabase.DB)
	syncService := sync_service.NewSyncService(syncRepo, transactionRepo, expenseRepo)
	syncHandler := sync_handler.NewSyncHandler(syncService)

	g := r.Group("/sync")
	{
		g.GET("/conflicts", middleware.RoleMiddleware("owner", "admin"), syncHandler.GetConflicts)
		g.GET("/conflicts/count", middleware.RoleMiddleware("owner", "admin"), syncHandler.GetConflictCount)
		g.POST("/conflicts/:id/resolve", middleware.RoleMiddleware("owner", "admin"), syncHandler.ResolveConflict)
		g.GET("/queue", middleware.RoleMiddleware("owner", "admin"), syncHandler.GetQueue)
		g.GET("/history", middleware.RoleMiddleware("owner", "admin"), syncHandler.GetHistory)
		g.POST("/push", syncHandler.PushSync)
	}
}
