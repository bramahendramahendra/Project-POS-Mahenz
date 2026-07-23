package segment

import (
	expiry_batch_handler "pos_api/domain/expiry_batch/handler"
	expiry_batch_repo "pos_api/domain/expiry_batch/repo"
	expiry_batch_service "pos_api/domain/expiry_batch/service"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// ExpiryBatchRoutes: nebeng permission module "produk.produk" (menu Produk) — bukan
// menu terpisah, karena warning/aksi expired ini bagian dari halaman Produk yang sama.
func ExpiryBatchRoutes(r *gin.RouterGroup) {
	expiryBatchRepo := expiry_batch_repo.NewExpiryBatchRepo(pkgdatabase.DB)
	expiryBatchService := expiry_batch_service.NewExpiryBatchService(expiryBatchRepo)
	expiryBatchHandler := expiry_batch_handler.NewExpiryBatchHandler(expiryBatchService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "produk.produk", action)
	}

	g := r.Group("/expiry-batches")
	{
		g.POST("/warnings", perm("can_view"), expiryBatchHandler.GetWarnings)
		g.POST("/product/:id", perm("can_view"), expiryBatchHandler.GetByProduct)
		g.POST("/confirm/:id", perm("can_edit"), expiryBatchHandler.Confirm)
		g.POST("/write-off/:id", perm("can_delete"), expiryBatchHandler.WriteOff)
	}
}
