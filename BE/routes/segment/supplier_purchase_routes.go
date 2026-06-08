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
	purchaseSvc := purchase_service.NewPurchaseService(purchaseRepo)
	purchaseHand := purchase_handler.NewPurchaseHandler(purchaseSvc)

	g := r.Group("/supplier-purchases")
	{
		g.POST("/generate-code", purchaseHand.GenerateCode)
		g.POST("/list", purchaseHand.GetAll)
		g.POST("/detail/:id", purchaseHand.GetByID)
		g.POST("/detail/:id/items", purchaseHand.GetItems)
		g.POST("/detail/:id/payments", purchaseHand.GetPayments)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), purchaseHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHand.Delete)
		g.POST("/pay/:id", middleware.RoleMiddleware("owner", "admin"), purchaseHand.Pay)
	}
}
