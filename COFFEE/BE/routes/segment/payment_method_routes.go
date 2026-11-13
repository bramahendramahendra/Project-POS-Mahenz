package segment

import (
	payment_method_handler "coffee_api/domain/payment_method/handler"
	payment_method_repo "coffee_api/domain/payment_method/repo"
	payment_method_service "coffee_api/domain/payment_method/service"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func PaymentMethodRoutes(r *gin.RouterGroup) {
	paymentMethodRepo := payment_method_repo.NewPaymentMethodRepo(pkgdatabase.DB)
	paymentMethodService := payment_method_service.NewPaymentMethodService(paymentMethodRepo)
	paymentMethodHandler := payment_method_handler.NewPaymentMethodHandler(paymentMethodService)

	g := r.Group("/payment-methods")
	{
		g.POST("/list", paymentMethodHandler.GetAll)
	}
}
