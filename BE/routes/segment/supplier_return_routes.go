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
	supplierReturnService := supplier_return_service.NewSupplierReturnService(supplierReturnRepo)
	supplierReturnHandler := supplier_return_handler.NewSupplierReturnHandler(supplierReturnService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pengadaan.retur", action)
	}

	g := r.Group("/supplier-returns")
	{
		g.POST("/list", supplierReturnHandler.GetAll)
		g.POST("/detail/:id", supplierReturnHandler.GetByID)
		g.POST("/create", perm("can_create"), supplierReturnHandler.Create)
		g.POST("/update-status/:id", perm("can_edit"), supplierReturnHandler.UpdateStatus)
		g.POST("/delete/:id", perm("can_delete"), supplierReturnHandler.Delete)
	}
}
