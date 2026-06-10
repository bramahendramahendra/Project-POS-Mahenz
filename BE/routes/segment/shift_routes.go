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
	shiftSvc := shift_service.NewShiftService(shiftRepo)
	shiftHand := shift_handler.NewShiftHandler(shiftSvc)

	g := r.Group("/shifts")
	{
		g.POST("/list", shiftHand.GetAll)
		g.POST("/active", shiftHand.GetActive)
		g.POST("/summary", middleware.RoleMiddleware("owner", "admin"), shiftHand.GetSummary)
		g.POST("/detail/:id", shiftHand.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), shiftHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), shiftHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), shiftHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), shiftHand.ToggleStatus)
	}
}
