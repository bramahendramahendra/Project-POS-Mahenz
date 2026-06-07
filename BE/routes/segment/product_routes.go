package segment

import (
	product_category_repo "pos_api/domain/product_category/repo"
	product_handler "pos_api/domain/product/handler"
	product_repo "pos_api/domain/product/repo"
	product_service "pos_api/domain/product/service"
	unit_repo "pos_api/domain/product_unit/repo"
	middleware "pos_api/middleware"
	pkgdatabase "pos_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.RouterGroup) {
	categoryRepo := product_category_repo.NewCategoryRepo(pkgdatabase.DB)
	masterUnitRepo := unit_repo.NewUnitRepo(pkgdatabase.DB)

	productRepo := product_repo.NewProductRepo(pkgdatabase.DB)
	productPackageRepo := product_repo.NewProductPackageRepo(pkgdatabase.DB)
	productSvc := product_service.NewProductService(productRepo, categoryRepo, productPackageRepo, masterUnitRepo)
	productHand := product_handler.NewProductHandler(productSvc)

	productPackageSvc := product_service.NewProductPackageService(productPackageRepo, productRepo)
	productPackageHand := product_handler.NewProductPackageHandler(productPackageSvc)

	productPriceRepo := product_repo.NewProductPriceRepo(pkgdatabase.DB)
	productPriceSvc := product_service.NewProductPriceService(productPriceRepo, productRepo)
	productPriceHand := product_handler.NewProductPriceHandler(productPriceSvc)

	g := r.Group("/products")
	{
		g.POST("/list", productHand.GetAll)
		g.POST("/options", productHand.GetOptions)
		g.POST("/search", productHand.Search)
		g.POST("/detail/:id", productHand.GetByID)
		g.POST("/by-barcode/:barcode", productHand.GetByBarcode)
		g.POST("/generate-barcode", middleware.RoleMiddleware("owner", "admin"), productHand.GenerateBarcode)
		g.POST("/generate-sku", middleware.RoleMiddleware("owner", "admin"), productHand.GenerateSku)
		g.POST("/import-template", middleware.RoleMiddleware("owner", "admin"), productHand.DownloadImportTemplate)
		g.POST("/import", middleware.RoleMiddleware("owner", "admin"), productHand.Import)
		g.POST("/import-preview", middleware.RoleMiddleware("owner", "admin"), productHand.ImportPreview)
		g.POST("/import-bulk", middleware.RoleMiddleware("owner", "admin"), productHand.ImportBulk)
		g.POST("/create", middleware.RoleMiddleware("owner", "admin"), productHand.Create)
		g.POST("/update/:id", middleware.RoleMiddleware("owner", "admin"), productHand.Update)
		g.POST("/delete/:id", middleware.RoleMiddleware("owner"), productHand.Delete)
		g.POST("/toggle-status/:id", middleware.RoleMiddleware("owner", "admin"), productHand.ToggleStatus)

		g.POST("/:id/packages/list", productPackageHand.GetByProduct)
		g.POST("/:id/packages/save", middleware.RoleMiddleware("owner", "admin"), productPackageHand.Save)
		g.POST("/:id/packages/delete/:package_id", middleware.RoleMiddleware("owner", "admin"), productPackageHand.Delete)

		g.POST("/:id/prices/list", productPriceHand.GetByProduct)
		g.POST("/:id/prices/save", middleware.RoleMiddleware("owner", "admin"), productPriceHand.Save)
	}
}
