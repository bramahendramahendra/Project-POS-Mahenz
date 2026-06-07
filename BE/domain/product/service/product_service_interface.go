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
		GetAll(req *dto.ProductListRequest) (data []*dto.ProductResponse, total int64, err error)
		GetOptions() (data []*dto.ProductOption, err error)
		GetCategoryNames() (data []string, err error)
		GetUnitNames() (data []string, err error)
		GetUnitInfos() (data []*dto.UnitInfo, err error)
		GetByID(id int) (data dto.ProductResponse, err error)
		GetByBarcode(barcode string) (data dto.ProductResponse, err error)
		Search(keyword string, limit int) (data []*dto.ProductSearchResult, err error)
		GetLowStock() (data []*dto.LowStockProduct, err error)
		GenerateBarcode() (data dto.GenerateBarcodeResponse, err error)
		GenerateSku(categoryID int) (data dto.GenerateSkuResponse, err error)
		Create(req *dto.ProductRequest) (data dto.ProductResponse, err error)
		Update(req *dto.UpdateProductRequest) (data dto.ProductResponse, err error)
		Delete(id int) error
		ToggleStatus(id int) error
		ImportFromFile(file *multipart.FileHeader) (data dto.ImportResult, err error)
		ImportBulk(req dto.BulkImportRequest) (data dto.BulkImportResult, err error)
		ImportPreview(file *multipart.FileHeader) (data dto.ImportPreviewResponse, err error)
	}

	productService struct {
		repo           repo.ProductRepo
		catRepo        repo_category.CategoryRepoInterface
		packageRepo    repo.ProductPackageRepo
		masterUnitRepo repo_unit.UnitRepoInterface
	}
)

func NewProductService(
	repo repo.ProductRepo,
	catRepo repo_category.CategoryRepoInterface,
	packageRepo repo.ProductPackageRepo,
	masterUnitRepo repo_unit.UnitRepoInterface,
) *productService {
	return &productService{repo: repo, catRepo: catRepo, packageRepo: packageRepo, masterUnitRepo: masterUnitRepo}
}
