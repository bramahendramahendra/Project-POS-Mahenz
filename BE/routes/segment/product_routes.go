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

	g := r.Group("/products")
	{
		g.POST("/list", productHandler.GetAll)
		g.POST("/options", productHandler.GetOptions)
		g.POST("/search", productHandler.Search)
		g.POST("/detail/:id", productHandler.GetByID)
		g.POST("/by-barcode/:barcode", productHandler.GetByBarcode)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), productHandler.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), productHandler.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), productHandler.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), productHandler.ToggleStatus)

		g.POST("/import", middleware.RoleMiddleware("owner", "admin"), productImportHandler.Import)
		g.POST("/import-preview", middleware.RoleMiddleware("owner", "admin"), productImportHandler.ImportPreview)
		g.POST("/import-bulk", middleware.RoleMiddleware("owner", "admin"), productImportHandler.ImportBulk)
		g.POST("/import-template", middleware.RoleMiddleware("owner", "admin"), productImportHandler.DownloadImportTemplate)

		g.POST("/generate-barcode", middleware.RoleMiddleware("owner", "admin"), productGenerateHandler.GenerateBarcode)
		g.POST("/generate-sku", middleware.RoleMiddleware("owner", "admin"), productGenerateHandler.GenerateSku)

		g.POST("/:id/packages/list", productPackageHandler.GetPackagesByProduct)
		g.POST("/:id/packages/save", middleware.RoleMiddleware("owner", "admin"), productPackageHandler.SavePackages)
		g.POST("/:id/packages/delete/:package_id", middleware.RoleMiddleware("owner", "admin"), productPackageHandler.DeletePackage)

		g.POST("/:id/prices/list", productPriceHandler.GetPricesByProduct)
		g.POST("/:id/prices/save", middleware.RoleMiddleware("owner", "admin"), productPriceHandler.SavePrices)
	}
}
