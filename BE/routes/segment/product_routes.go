package segment

import (
	product_handler "pos_api/domain/product/handler"
	product_repo "pos_api/domain/product/repo"
	product_service "pos_api/domain/product/service"
	product_category_repo "pos_api/domain/product_category/repo"
	unit_repo "pos_api/domain/product_unit/repo"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.RouterGroup) {
	categoryRepo := product_category_repo.NewCategoryRepo(pkgdatabase.DB)
	masterUnitRepo := unit_repo.NewUnitRepo(pkgdatabase.DB)

	productRepo := product_repo.NewProductRepo(pkgdatabase.DB)
	productService := product_service.NewProductService(productRepo, categoryRepo, masterUnitRepo)
	productHandler := product_handler.NewProductHandler(productService)
	productImportHandler := product_handler.NewProductImportHandler(productService)
	productGenerateHandler := product_handler.NewProductGenerateHandler(productService)
	productPackageHandler := product_handler.NewProductPackageHandler(productService)
	productPriceHandler := product_handler.NewProductPriceHandler(productService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "produk.produk", action)
	}

	g := r.Group("/products")
	{
		g.POST("/list", productHandler.GetAll)
		g.POST("/options", productHandler.GetOptions)
		g.POST("/search", productHandler.Search)
		g.POST("/detail/:id", productHandler.GetByID)
		g.POST("/by-barcode/:barcode", productHandler.GetByBarcode)
		g.POST("/create", perm("can_create"), productHandler.Create)
		g.POST("/update/:id", perm("can_edit"), productHandler.Update)
		g.POST("/delete/:id", perm("can_delete"), productHandler.Delete)
		g.POST("/toggle-status/:id", perm("can_edit"), productHandler.ToggleStatus)

		g.POST("/import", perm("can_create"), productImportHandler.Import)
		g.POST("/import-preview", perm("can_create"), productImportHandler.ImportPreview)
		g.POST("/import-bulk", perm("can_create"), productImportHandler.ImportBulk)
		g.POST("/import-template", perm("can_view"), productImportHandler.DownloadImportTemplate)

		g.POST("/generate-barcode", perm("can_create"), productGenerateHandler.GenerateBarcode)
		g.POST("/generate-sku", perm("can_create"), productGenerateHandler.GenerateSku)

		g.POST("/:id/packages/list", productPackageHandler.GetPackagesByProduct)
		g.POST("/:id/packages/save", perm("can_edit"), productPackageHandler.SavePackages)
		g.POST("/:id/packages/delete/:package_id", perm("can_delete"), productPackageHandler.DeletePackage)

		g.POST("/:id/prices/list", productPriceHandler.GetPricesByProduct)
		g.POST("/:id/prices/save", perm("can_edit"), productPriceHandler.SavePrices)
	}
}
