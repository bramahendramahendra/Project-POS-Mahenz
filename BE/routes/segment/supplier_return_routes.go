package segment

import (
	supplier_return_handler "pos_api/domain/supplier_return/handler"
	supplier_return_repo "pos_api/domain/supplier_return/repo"
	supplier_return_service "pos_api/domain/supplier_return/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func SupplierReturnRoutes(r *gin.RouterGroup) {
	supplierReturnRepo := supplier_return_repo.NewSupplierReturnRepo(pkgdatabase.DB)
	supplierReturnSvc := supplier_return_service.NewSupplierReturnService(supplierReturnRepo)
	supplierReturnHand := supplier_return_handler.NewSupplierReturnHandler(supplierReturnSvc)

	g := r.Group("/supplier-returns")
	{
		g.POST("/list", supplierReturnHand.GetAll)
		g.POST("/detail/:id", supplierReturnHand.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), supplierReturnHand.Create)
		g.POST("/update-status/:id", middleware.RoleMiddleware("owner", "admin"), supplierReturnHand.UpdateStatus)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), supplierReturnHand.Delete)
	}
}
