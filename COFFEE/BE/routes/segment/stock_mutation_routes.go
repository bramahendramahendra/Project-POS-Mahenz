package segment

import (
	stock_mutation_handler "coffee_api/domain/stock_mutation/handler"
	stock_mutation_repo "coffee_api/domain/stock_mutation/repo"
	stock_mutation_service "coffee_api/domain/stock_mutation/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func StockMutationRoutes(r *gin.RouterGroup) {
	stockMutationRepo := stock_mutation_repo.NewStockMutationRepo(pkgdatabase.DB)
	stockMutationService := stock_mutation_service.NewStockMutationService(stockMutationRepo)
	stockMutationHandler := stock_mutation_handler.NewStockMutationHandler(stockMutationService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelaporan.stok", action)
	}

	g := r.Group("/stock-mutations")
	{
		g.POST("/list", perm("can_view"), stockMutationHandler.GetAll)
		g.POST("/product/:product_id", perm("can_view"), stockMutationHandler.GetByProduct)
	}
}
