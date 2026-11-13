package segment

import (
	cash_drawer_repo "coffee_api/domain/cash_drawer/repo"
	expense_handler "coffee_api/domain/expense/handler"
	expense_repo "coffee_api/domain/expense/repo"
	expense_service "coffee_api/domain/expense/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(r *gin.RouterGroup) {
	cashDrawerRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	expenseRepo := expense_repo.NewExpenseRepo(pkgdatabase.DB)
	expenseService := expense_service.NewExpenseService(expenseRepo, cashDrawerRepo)
	expenseHandler := expense_handler.NewExpenseHandler(expenseService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "keuangan.pengeluaran", action)
	}

	g := r.Group("/expenses")
	{
		g.POST("/list", expenseHandler.GetAll)
		g.POST("/detail/:id", expenseHandler.GetByID)
		g.POST("/create", perm("can_create"), expenseHandler.Create)
		g.POST("/update/:id", perm("can_edit"), expenseHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), expenseHandler.Delete)
	}
}
