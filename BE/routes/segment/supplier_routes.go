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
	supplierService := supplier_service.NewSupplierService(supplierRepo)
	supplierHandler := supplier_handler.NewSupplierHandler(supplierService)

	g := r.Group("/suppliers")
	{
		g.POST("/list", supplierHandler.GetAll)
		g.POST("/options", supplierHandler.GetOptions)
		g.POST("/detail/:id", supplierHandler.GetDetail)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), supplierHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), supplierHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), supplierHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), supplierHandler.ToggleStatus)
	}
}
