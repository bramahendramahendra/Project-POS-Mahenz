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
	customerSvc := customer_service.NewCustomerService(customerRepo)
	customerHand := customer_handler.NewCustomerHandler(customerSvc)

	g := r.Group("/customers")
	{
		g.POST("/list", customerHand.GetAll)
		g.POST("/active", customerHand.GetOptions)
		g.POST("/detail/:id", customerHand.GetByID)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), customerHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), customerHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), customerHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), customerHand.ToggleStatus)
	}
}
