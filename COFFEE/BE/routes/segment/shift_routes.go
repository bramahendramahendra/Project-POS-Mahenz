package segment

import (
	shift_handler "coffee_api/domain/shift/handler"
	shift_repo "coffee_api/domain/shift/repo"
	shift_service "coffee_api/domain/shift/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

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
		// /active tanpa perm() disengaja — pola sama seperti endpoint options/lookup lain
		// (product-categories/options, suppliers/options, customers/active, dst): data referensi
		// buat dropdown, dibutuhkan lintas fitur (di sini oleh OpenCashDrawerModal), bukan cuma
		// oleh menu Shift itu sendiri. Jangan tempel perm() lagi di sini.
		g.POST("/active", shiftHandler.GetOptions)
		g.POST("/create", perm("can_create"), shiftHandler.Create)
		g.POST("/update/:id", perm("can_edit"), shiftHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), shiftHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), shiftHandler.ToggleStatus)
	}
}
