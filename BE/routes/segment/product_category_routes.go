package segment

import (
	category_handler "pos_api/domain/product_category/handler"
	category_repo "pos_api/domain/product_category/repo"
	category_service "pos_api/domain/product_category/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ProductCategoryRoutes(r *gin.RouterGroup) {
	categoryRepo := category_repo.NewCategoryRepo(pkgdatabase.DB)
	categorySvc := category_service.NewCategoryService(categoryRepo)
	categoryHand := category_handler.NewCategoryHandler(categorySvc)

	g := r.Group("/categories")
	{
		g.POST("/list", categoryHand.GetAll)
		g.POST("/detail/:id", categoryHand.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), categoryHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), categoryHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), categoryHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), categoryHand.ToggleStatus)
	}
}
