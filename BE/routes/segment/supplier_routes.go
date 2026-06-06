package segment

import (
	supplier_handler "pos_api/domain/supplier/handler"
	supplier_repo "pos_api/domain/supplier/repo"
	supplier_service "pos_api/domain/supplier/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func SupplierRoutes(r *gin.RouterGroup) {
	supplierRepo := supplier_repo.NewSupplierRepo(pkgdatabase.DB)
	supplierSvc := supplier_service.NewSupplierService(supplierRepo)
	supplierHand := supplier_handler.NewSupplierHandler(supplierSvc)

	g := r.Group("/suppliers")
	{
		g.POST("/list", supplierHand.GetAll)
		g.POST("/options", supplierHand.GetOptions)
		g.POST("/detail/:id", supplierHand.GetDetail)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), supplierHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), supplierHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), supplierHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), supplierHand.ToggleStatus)
	}
}
