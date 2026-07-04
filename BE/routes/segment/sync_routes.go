package segment

import (
	expense_repo "pos_api/domain/expense/repo"
	product_repo "pos_api/domain/product/repo"
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
	productRepo := product_repo.NewProductRepo(pkgdatabase.DB)
	syncService := sync_service.NewSyncService(syncRepo, transactionRepo, expenseRepo, productRepo)
	syncHandler := sync_handler.NewSyncHandler(syncService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "operasional.sync", action)
	}

	g := r.Group("/sync")
	{
		g.GET("/conflicts", perm("can_view"), syncHandler.GetConflicts)
		g.GET("/conflicts/count", perm("can_view"), syncHandler.GetConflictCount)
		g.POST("/conflicts/:id/resolve", perm("can_edit"), syncHandler.ResolveConflict)
		g.GET("/queue", perm("can_view"), syncHandler.GetQueue)
		g.GET("/history", perm("can_view"), syncHandler.GetHistory)
		g.POST("/push", perm("can_create"), syncHandler.PushSync)
	}
}
