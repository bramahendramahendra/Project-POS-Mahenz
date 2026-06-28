package segment

import (
	version_handler "pos_api/domain/version/handler"
	version_repo "pos_api/domain/version/repo"
	version_service "pos_api/domain/version/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func VersionAdminRoutes(r *gin.RouterGroup) {
	versionRepo := version_repo.NewVersionRepo(pkgdatabase.DB)
	versionService := version_service.NewVersionService(versionRepo)
	versionHandler := version_handler.NewVersionHandler(versionService)

	g := r.Group("/version")
	{
		g.POST("/list", middleware.RoleMiddleware("owner", "admin"), versionHandler.GetAll)
		g.POST("/android", middleware.RoleMiddleware("owner", "admin"), versionHandler.UpdateAndroidVersion)
	}
}
