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
	unitSvc := product_unit_service.NewUnitService(unitRepo)
	unitHand := product_unit_handler.NewUnitHandler(unitSvc)

	g := r.Group("/units")
	{
		g.POST("/list", unitHand.GetAll)
		g.POST("/options", unitHand.GetOptions)
		g.POST("/detail/:id", unitHand.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), unitHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), unitHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), unitHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), unitHand.ToggleStatus)
	}
}
