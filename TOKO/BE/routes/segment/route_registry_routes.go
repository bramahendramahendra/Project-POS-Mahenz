package segment

import (
	route_registry_handler "pos_api/domain/route_registry/handler"
	route_registry_repo "pos_api/domain/route_registry/repo"
	route_registry_service "pos_api/domain/route_registry/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func RouteRegistryRoutes(r *gin.RouterGroup) {
	routeRegistryRepo := route_registry_repo.NewRouteRegistryRepo(pkgdatabase.DB)
	routeRegistryService := route_registry_service.NewRouteRegistryService(routeRegistryRepo)
	routeRegistryHandler := route_registry_handler.NewRouteRegistryHandler(routeRegistryService)

	svc := newAccessService()
	perm := middleware.PermissionMiddleware(svc, "sistem.menus", "can_view")

	g := r.Group("/route-registry")
	{
		// Opsi path buat dropdown di form Manajemen Menu. Read-only — tidak ada
		// create/update/delete di sini, tabel ini hanya diisi lewat migration developer.
		g.POST("/options", perm, routeRegistryHandler.GetOptions)
	}
}
