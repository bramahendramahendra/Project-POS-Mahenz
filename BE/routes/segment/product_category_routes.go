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
	categoryService := category_service.NewCategoryService(categoryRepo)
	categoryHandler := category_handler.NewCategoryHandler(categoryService)

	g := r.Group("/categories")
	{
		g.POST("/list", categoryHandler.GetAll)
		g.POST("/options", categoryHandler.GetOptions)
		g.POST("/detail/:id", categoryHandler.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), categoryHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), categoryHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), categoryHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), categoryHandler.ToggleStatus)
	}
}
