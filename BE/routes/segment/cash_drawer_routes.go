package segment

import (
	cash_drawer_handler "pos_api/domain/cash_drawer/handler"
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func CashDrawerRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	cashDrawerService := cash_drawer_service.NewCashDrawerService(cashDrawerRepo)
	cashDrawerHandler := cash_drawer_handler.NewCashDrawerHandler(cashDrawerService)

	g := r.Group("/cash-drawer")
	{
		g.POST("/current", cashDrawerHandler.GetCurrent)
		g.POST("/list", cashDrawerHandler.GetHistory)
		g.POST("/detail/:id", cashDrawerHandler.GetByID)
		g.POST("/open", cashDrawerHandler.Open)
		g.POST("/close/:id", cashDrawerHandler.Close)
		g.POST("/update-sales/:id", cashDrawerHandler.UpdateSales)
		g.POST("/update-expenses/:id", cashDrawerHandler.UpdateExpenses)
	}
}
