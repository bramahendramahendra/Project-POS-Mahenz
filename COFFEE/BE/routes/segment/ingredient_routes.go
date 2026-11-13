package segment

import (
	ingredient_handler "coffee_api/domain/ingredient/handler"
	ingredient_repo "coffee_api/domain/ingredient/repo"
	ingredient_service "coffee_api/domain/ingredient/service"
	middleware "coffee_api/middleware"
	pkgdatabase "coffee_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// IngredientRoutes: bahan baku (tahap 4.3, keputusan 80-84).
//   - bahan.daftar : lihat daftar, kelola bahan (tambah, ubah, aktif/nonaktif, hapus)
//   - bahan.riwayat: lihat riwayat, catat terpakai/terbuang (can_create), batalkan entri (can_delete)
func IngredientRoutes(r *gin.RouterGroup) {
	h := ingredient_handler.NewIngredientHandler(
		ingredient_service.NewIngredientService(ingredient_repo.NewIngredientRepo(pkgdatabase.DB)),
	)

	svc := newAccessService()
	daftar := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "bahan.daftar", action)
	}
	riwayat := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "bahan.riwayat", action)
	}

	g := r.Group("/ingredients")
	{
		g.POST("/list", daftar("can_view"), h.GetAll)
		g.POST("/options", daftar("can_view"), h.GetOptions)
		g.POST("/detail/:id", daftar("can_view"), h.GetByID)
		g.POST("/create", daftar("can_create"), h.Create)
		g.POST("/update/:id", daftar("can_edit"), h.Update)
		g.POST("/toggle-status/:id", daftar("can_edit"), h.ToggleStatus)
		g.POST("/delete/:id", daftar("can_delete"), h.Delete)

		g.POST("/usage", riwayat("can_create"), h.RecordUsage)
		g.POST("/waste", riwayat("can_create"), h.RecordWaste)
		g.POST("/mutations/list", riwayat("can_view"), h.GetMutations)
		g.POST("/mutations/cancel/:id", riwayat("can_delete"), h.CancelMutation)
	}
}
