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

	g := r.Group("/customers")
	{
		g.POST("/list", customerHandler.GetAll)
		g.POST("/active", customerHandler.GetOptions)
		g.POST("/detail/:id", customerHandler.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), customerHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), customerHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), customerHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), customerHandler.ToggleStatus)
	}
}
