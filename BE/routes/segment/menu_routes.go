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
	menuSvc := menu_service.NewMenuService(menuRepo)
	menuHand := menu_handler.NewMenuHandler(menuSvc)

	g := r.Group("/menus")
	{
		g.POST("/my",            menuHand.GetMyMenus)
		g.POST("/list",          middleware.RoleMiddleware("owner", "admin"), menuHand.GetAll)
		g.POST("/detail/:id",    middleware.RoleMiddleware("owner", "admin"), menuHand.GetByID)
		g.POST("/create",        middleware.RoleMiddleware("owner"),          menuHand.Create)
		g.POST("/update/:id",    middleware.RoleMiddleware("owner"),          menuHand.Update)
		g.POST("/delete/:id",    middleware.RoleMiddleware("owner"),          menuHand.Delete)
		g.POST("/reorder",       middleware.RoleMiddleware("owner"),          menuHand.Reorder)
	}
}
