package service

import (
	"mime/multipart"

	dto "pos_api/domain/product/dto"
	repo "pos_api/domain/product/repo"
	repo_category "pos_api/domain/product_category/repo"
	repo_unit "pos_api/domain/product_unit/repo"
)

type (
	ProductServiceInterface interface {
		GetAll(req *dto.ProductListRequest) (data []dto.ProductResponse, total int64, err error)
		GetOptions() (data []*dto.ProductOption, err error)
		Search(req *dto.SearchProductRequest) (data []*dto.ProductSearchResult, err error)
		GetByID(id int) (data dto.ProductResponse, err error)
		GetByBarcode(barcode string) (data dto.ProductResponse, err error)
		Create(req *dto.ProductRequest) (data dto.ProductResponse, err error)
		Update(req *dto.UpdateProductRequest) (data dto.ProductResponse, err error)
		Delete(req *dto.DeleteProductRequest) error
		ToggleStatus(req *dto.ToggleStatusProductRequest) error

		ImportFromFile(file *multipart.FileHeader) (data dto.ImportResult, err error)
		ImportPreview(file *multipart.FileHeader) (data dto.ImportPreviewResponse, err error)
		ImportBulk(req dto.BulkImportRequest) (data dto.BulkImportResult, err error)

		GenerateBarcode() (data dto.GenerateBarcodeResponse, err error)
		GenerateSku(categoryID int) (data dto.GenerateSkuResponse, err error)

		GetPackagesByProduct(productID int) (data []*dto.ProductPackageResponse, err error)
		SavePackages(req *dto.SaveProductPackagesRequest) error
		DeletePackage(req *dto.PackageIDUriRequest) error

		GetPricesByProduct(productID int) (data []*dto.ProductPriceResponse, err error)
		SavePrices(req *dto.SaveProductPricesRequest) error

		GetLowStock() (data []*dto.LowStockProduct, err error)
		GetCategoryNames() (data []string, err error)
		GetUnitInfos() (data []*dto.UnitInfo, err error)
	}

	productService struct {
		repo         repo.ProductRepo
		repoCategory repo_category.CategoryRepoInterface
		repoUnit     repo_unit.UnitRepoInterface
	}
)

func NewProductService(
	repo repo.ProductRepo,
	repoCategory repo_category.CategoryRepoInterface,
	repoUnit repo_unit.UnitRepoInterface,
) *productService {
	return &productService{repo: repo, repoCategory: repoCategory, repoUnit: repoUnit}
}
