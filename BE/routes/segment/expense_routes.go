package segment

import (
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	expense_handler "pos_api/domain/expense/handler"
	expense_repo "pos_api/domain/expense/repo"
	expense_service "pos_api/domain/expense/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	expenseRepo := expense_repo.NewExpenseRepo(pkgdatabase.DB)
	expenseService := expense_service.NewExpenseService(expenseRepo, cashDrawerRepo)
	expenseHandler := expense_handler.NewExpenseHandler(expenseService)

	g := r.Group("/expenses")
	{
		g.POST("/list", expenseHandler.GetAll)
		g.POST("/detail/:id", expenseHandler.GetByID)
		g.POST("/create", expenseHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), expenseHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner", "admin"), expenseHandler.Delete)
	}
}
