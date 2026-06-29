package segment

import (
	product_unit_handler "pos_api/domain/product_unit/handler"
	product_unit_repo "pos_api/domain/product_unit/repo"
	product_unit_service "pos_api/domain/product_unit/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ProductUnitRoutes(r *gin.RouterGroup) {
	unitRepo := product_unit_repo.NewUnitRepo(pkgdatabase.DB)
	unitService := product_unit_service.NewUnitService(unitRepo)
	unitHandler := product_unit_handler.NewUnitHandler(unitService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "produk.unit", action)
	}

	g := r.Group("/units")
	{
		g.POST("/list", unitHandler.GetAll)
		g.POST("/options", unitHandler.GetOptions)
		g.POST("/detail/:id", unitHandler.GetByID)
		g.POST("/create", perm("can_create"), unitHandler.Create)
		g.POST("/update/:id", perm("can_edit"), unitHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), unitHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), unitHandler.ToggleStatus)
	}
}
