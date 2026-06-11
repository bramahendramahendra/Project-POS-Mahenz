package segment

import (
	shift_handler "pos_api/domain/shift/handler"
	shift_repo "pos_api/domain/shift/repo"
	shift_service "pos_api/domain/shift/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ShiftRoutes(r *gin.RouterGroup) {
	shiftRepo := shift_repo.NewShiftRepo(pkgdatabase.DB)
	shiftService := shift_service.NewShiftService(shiftRepo)
	shiftHandler := shift_handler.NewShiftHandler(shiftService)

	g := r.Group("/shifts")
	{
		g.POST("/list", shiftHandler.GetAll)
		g.POST("/active", shiftHandler.GetOptions)
		g.POST("/summary", middleware.RoleMiddleware("owner", "admin"), shiftHandler.GetSummary)
		g.POST("/detail/:id", shiftHandler.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), shiftHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), shiftHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), shiftHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), shiftHandler.ToggleStatus)
	}
}
