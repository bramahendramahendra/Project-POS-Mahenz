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

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pengadaan.pembelian", action)
	}

	g := r.Group("/supplier-purchases")
	{
		g.POST("/generate-code", purchaseHandler.GenerateCode)
		g.POST("/list", purchaseHandler.GetAll)
		g.POST("/detail/:id", purchaseHandler.GetByID)
		g.POST("/:id/payments", purchaseHandler.GetPayments)
		g.POST("/create", perm("can_create"), purchaseHandler.Create)
		g.POST("/update/:id", perm("can_edit"), purchaseHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), purchaseHandler.Delete)
		g.POST("/pay/:id", perm("can_edit"), purchaseHandler.Pay)
	}
}
