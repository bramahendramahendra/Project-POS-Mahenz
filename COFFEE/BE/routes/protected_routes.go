package routes

import (
	auth_repo "coffee_api/domain/auth/repo"
	auth_service "coffee_api/domain/auth/service"
	pos_middleware "coffee_api/middleware/auth"
	"coffee_api/pkg/authclient"
	pkgdatabase "coffee_api/pkg/database"
	segment "coffee_api/routes/segment"

	"github.com/gin-gonic/gin"
)

func protectedRoutes(r *gin.RouterGroup, authClient *authclient.Client) {
	authRepo := auth_repo.NewAuthRepo(pkgdatabase.DB)
	authSvc := auth_service.NewAuthService(authRepo, authClient)

	r.Use(pos_middleware.POSBearerAuthMiddleware(authSvc))

	segment.AuthRoutes(r, authClient)
	segment.PinRoutes(r)
	segment.UserRoutes(r)
	segment.ProductCategoryRoutes(r)
	segment.ProductUnitRoutes(r)
	segment.ProductRoutes(r)
	segment.TransactionRoutes(r)
	segment.CashDrawerRoutes(r)
	segment.ExpenseRoutes(r)
	segment.PaymentStatusRoutes(r)
	segment.PaymentMethodRoutes(r)
	segment.PurchaseRoutes(r)
	segment.SupplierRoutes(r)
	segment.CustomerRoutes(r)
	segment.ReceivableRoutes(r)
	segment.ShiftRoutes(r)
	segment.StockMutationRoutes(r)
	segment.StockReconciliationRoutes(r)
	segment.FinanceRoutes(r)
	segment.ReportRoutes(r)
	segment.BusinessSummaryRoutes(r)
	segment.DashboardRoutes(r)
	segment.SettingRoutes(r)
	segment.PrintRoutes(r)
	segment.BackupRoutes(r)
	segment.RoleRoutes(r)
	segment.MenuRoutes(r)
	segment.AccessRoutes(r)
	segment.RouteRegistryRoutes(r)
}
