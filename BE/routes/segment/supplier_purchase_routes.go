package segment

import (
	purchase_handler "pos_api/domain/supplier_purchase/handler"
	purchase_repo "pos_api/domain/supplier_purchase/repo"
	purchase_service "pos_api/domain/supplier_purchase/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func PurchaseRoutes(r *gin.RouterGroup) {
	purchaseRepo := purchase_repo.NewPurchaseRepo(pkgdatabase.DB)
	purchaseService := purchase_service.NewPurchaseService(purchaseRepo)
	purchaseHandler := purchase_handler.NewPurchaseHandler(purchaseService)

	g := r.Group("/supplier-purchases")
	{
		g.POST("/generate-code", purchaseHandler.GenerateCode)
		g.POST("/list", purchaseHandler.GetAll)
		g.POST("/detail/:id", purchaseHandler.GetByID)
		g.POST("/detail/:id/items", purchaseHandler.GetItems)
		g.POST("/detail/:id/payments", purchaseHandler.GetPayments)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), purchaseHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHandler.Delete)
		g.POST("/pay/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHandler.Pay)
	}
}
