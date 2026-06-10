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
	cashDrawerSvc := cash_drawer_service.NewCashDrawerService(cashDrawerRepo)
	cashDrawerHand := cash_drawer_handler.NewCashDrawerHandler(cashDrawerSvc)

	g := r.Group("/cash-drawer")
	{
		g.POST("/current", cashDrawerHand.GetCurrent)
		g.POST("/list", cashDrawerHand.GetHistory)
		g.POST("/detail/:id", cashDrawerHand.GetByID)
		g.POST("/open", cashDrawerHand.Open)
		g.POST("/close/:id", cashDrawerHand.Close)
		g.POST("/update-sales/:id", cashDrawerHand.UpdateSales)
		g.POST("/update-expenses/:id", cashDrawerHand.UpdateExpenses)
	}
}
