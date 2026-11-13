package segment

import (
	payment_status_handler "coffee_api/domain/payment_status/handler"
	payment_status_repo "coffee_api/domain/payment_status/repo"
	payment_status_service "coffee_api/domain/payment_status/service"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func PaymentStatusRoutes(r *gin.RouterGroup) {
	paymentStatusRepo := payment_status_repo.NewPaymentStatusRepo(pkgdatabase.DB)
	paymentStatusService := payment_status_service.NewPaymentStatusService(paymentStatusRepo)
	paymentStatusHandler := payment_status_handler.NewPaymentStatusHandler(paymentStatusService)

	g := r.Group("/payment-statuses")
	{
		g.POST("/list", paymentStatusHandler.GetAll)
	}
}
