package segment

import (
	payment_method_handler "pos_api/domain/payment_method/handler"
	payment_method_repo "pos_api/domain/payment_method/repo"
	payment_method_service "pos_api/domain/payment_method/service"
	pkgdatabase "pos_api/pkg/database"

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
