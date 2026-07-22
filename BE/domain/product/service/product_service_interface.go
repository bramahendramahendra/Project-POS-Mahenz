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
		GetAll(req *dto.GetAllRequest) (data []dto.ProductResponse, total int64, err error)
		GetOptions() (data []*dto.GetOptionResponse, err error)
		Search(req *dto.SearchRequest) (data []*dto.SearchResponse, err error)
		GetByID(id int) (data dto.ProductResponse, err error)
		GetByBarcode(barcode string) (data dto.ProductResponse, err error)
		Create(req *dto.CreateRequest, role string) (data dto.ProductResponse, err error)
		Update(req *dto.UpdateRequest, role string) (data dto.ProductResponse, err error)
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error

		ImportPreview(file *multipart.FileHeader) (data dto.ImportPreviewResponse, err error)
		ImportBulk(req dto.BulkImportRequest) (data dto.BulkImportResult, err error)

		GenerateBarcode() (data dto.GenerateBarcodeResponse, err error)
		GenerateSku(categoryID int) (data dto.GenerateSkuResponse, err error)

		GetPackagesByProduct(productID int) (data []*dto.PackageResponse, err error)
		CreatePackage(productID int, req *dto.CreatePackageRequest) (data *dto.PackageResponse, err error)
		UpdatePackage(id, productID int, req *dto.UpdatePackageRequest) (data *dto.PackageResponse, err error)
		DeletePackage(req *dto.DeletePackageRequest) (err error)

		GetPricesByProduct(productID int) (data []*dto.PriceResponse, err error)
		SavePrices(req *dto.SavePriceRequest) (err error)

		GetLowStock() (data []*dto.GetLowStockResponse, err error)
		GetCategoryNames() (data []string, err error)
		GetUnitInfos() (data []*dto.GetUnitInfoResponse, err error)
	}

	productService struct {
		repo         repo.ProductRepoInterface
		repoCategory repo_category.CategoryRepoInterface
		repoUnit     repo_unit.UnitRepoInterface
	}
)

func NewProductService(
	repo repo.ProductRepoInterface,
	repoCategory repo_category.CategoryRepoInterface,
	repoUnit repo_unit.UnitRepoInterface,
) *productService {
	return &productService{repo: repo, repoCategory: repoCategory, repoUnit: repoUnit}
}
