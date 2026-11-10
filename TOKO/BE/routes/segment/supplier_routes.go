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

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pengadaan.supplier", action)
	}

	g := r.Group("/suppliers")
	{
		g.POST("/list", supplierHandler.GetAll)
		g.POST("/options", supplierHandler.GetOptions)
		g.POST("/detail/:id", supplierHandler.GetDetail)
		g.POST("/create", perm("can_create"), supplierHandler.Create)
		g.POST("/update/:id", perm("can_edit"), supplierHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), supplierHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), supplierHandler.ToggleStatus)
	}
}
