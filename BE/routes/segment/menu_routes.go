package segment

import (
	menu_handler "pos_api/domain/menu/handler"
	menu_repo "pos_api/domain/menu/repo"
	menu_service "pos_api/domain/menu/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(r *gin.RouterGroup) {
	menuRepo := menu_repo.NewMenuRepo(pkgdatabase.DB)
	menuService := menu_service.NewMenuService(menuRepo)
	menuHandler := menu_handler.NewMenuHandler(menuService)

	g := r.Group("/menus")
	{
		g.POST("/my", menuHandler.GetMyMenus)
		g.POST("/list", middleware.RoleMiddleware("owner", "admin"), menuHandler.GetAll)
		g.POST("/detail/:id", middleware.RoleMiddleware("owner", "admin"), menuHandler.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner"), menuHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner"), menuHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), menuHandler.Delete)
		g.POST("/reorder", middleware.RoleMiddleware("owner"), menuHandler.Reorder)
	}
}
