package segment

import (
	category_handler "coffee_api/domain/product_category/handler"
	category_repo "coffee_api/domain/product_category/repo"
	category_service "coffee_api/domain/product_category/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ProductCategoryRoutes(r *gin.RouterGroup) {
	categoryRepo := category_repo.NewCategoryRepo(pkgdatabase.DB)
	categoryService := category_service.NewCategoryService(categoryRepo)
	categoryHandler := category_handler.NewCategoryHandler(categoryService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "produk.kategori", action)
	}

	g := r.Group("/categories")
	{
		// Category Handler
		g.POST("/list", categoryHandler.GetAll)
		g.POST("/options", categoryHandler.GetOptions)
		g.POST("/create", perm("can_create"), categoryHandler.Create)
		g.POST("/update/:id", perm("can_edit"), categoryHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), categoryHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), categoryHandler.ToggleStatus)
	}
}
