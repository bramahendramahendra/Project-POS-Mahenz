package segment

import (
	menu_handler "pos_api/domain/menu/handler"
	menu_repo "pos_api/domain/menu/repo"
	menu_service "pos_api/domain/menu/service"
	route_registry_repo "pos_api/domain/route_registry/repo"
	route_registry_service "pos_api/domain/route_registry/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(r *gin.RouterGroup) {
	menuRepo := menu_repo.NewMenuRepo(pkgdatabase.DB)
	routeRegistryRepo := route_registry_repo.NewRouteRegistryRepo(pkgdatabase.DB)
	routeRegistryService := route_registry_service.NewRouteRegistryService(routeRegistryRepo)
	menuService := menu_service.NewMenuService(menuRepo, routeRegistryService)
	menuHandler := menu_handler.NewMenuHandler(menuService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "sistem.menus", action)
	}

	g := r.Group("/menus")
	{
		g.POST("/my", menuHandler.GetMyMenus)
		g.POST("/options", perm("can_view"), menuHandler.GetOptions)
		g.POST("/list", perm("can_view"), menuHandler.GetAll)
		g.POST("/detail/:id", perm("can_view"), menuHandler.GetByID)
		g.POST("/create", perm("can_create"), menuHandler.Create)
		g.POST("/update/:id", perm("can_edit"), menuHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), menuHandler.Delete)
		g.POST("/reorder", perm("can_edit"), menuHandler.Reorder)
	}
}
