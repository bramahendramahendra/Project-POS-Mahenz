package segment

import (
	customer_handler "pos_api/domain/customer/handler"
	customer_repo "pos_api/domain/customer/repo"
	customer_service "pos_api/domain/customer/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(r *gin.RouterGroup) {
	customerRepo := customer_repo.NewCustomerRepo(pkgdatabase.DB)
	customerService := customer_service.NewCustomerService(customerRepo)
	customerHandler := customer_handler.NewCustomerHandler(customerService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "pelanggan.pelanggan", action)
	}

	g := r.Group("/customers")
	{
		g.POST("/list", customerHandler.GetAll)
		g.POST("/active", customerHandler.GetOptions)
		g.POST("/detail/:id", customerHandler.GetByID)
		g.POST("/create", perm("can_create"), customerHandler.Create)
		g.POST("/update/:id", perm("can_edit"), customerHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), customerHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), customerHandler.ToggleStatus)
	}
}
