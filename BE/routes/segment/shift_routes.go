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

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "operasional.shift", action)
	}

	g := r.Group("/shifts")
	{
		g.POST("/list", perm("can_view"), shiftHandler.GetAll)
		g.POST("/active", perm("can_view"), shiftHandler.GetOptions)
		g.POST("/summary", perm("can_view"), shiftHandler.GetSummary)
		g.POST("/detail/:id", perm("can_view"), shiftHandler.GetByID)
		g.POST("/create", perm("can_create"), shiftHandler.Create)
		g.POST("/update/:id", perm("can_edit"), shiftHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), shiftHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), shiftHandler.ToggleStatus)
	}
}
