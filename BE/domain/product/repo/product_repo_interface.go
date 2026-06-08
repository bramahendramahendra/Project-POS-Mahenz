package repo_product

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"

	"gorm.io/gorm"
)

type ProductRepo interface {
	GetAll(req *dto_product.ProductListRequest) ([]*model_product.Product, int64, error)
	GetOptions() ([]*model_product.ProductOption, error)
	GetByID(id int) (*model_product.Product, error)
	GetByBarcode(barcode string) (*model_product.Product, error)
	Search(keyword string, limit int) ([]*model_product.ProductSearchResult, error)
	GetLowStock() ([]*model_product.LowStockProduct, error)
	CheckBarcodeExists(barcode string, excludeID int) (bool, error)
	CheckSkuExists(sku string, excludeID int) (bool, error)
	CountSkuByCategory(categoryID int) (int, error)
	CountTransactionItems(productID int) (int, error)
	Create(req *dto_product.ProductRequest) (int64, error)
	Update(req *dto_product.UpdateProductRequest) error
	Delete(req *dto_product.DeleteProductRequest) error
	ToggleStatus(req *dto_product.ToggleStatusProductRequest) error
	UpdateStock(id int, delta float64) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{db: db}
}

type ProductPriceRepo interface {
	GetByProduct(productID int) ([]*model_product.ProductPrice, error)
	Save(productID int, prices []dto_product.ProductPriceRequest) error
}

type productPriceRepo struct {
	db *gorm.DB
}

func NewProductPriceRepo(db *gorm.DB) ProductPriceRepo {
	return &productPriceRepo{db: db}
}

type ProductPackageRepo interface {
	GetByProduct(productID int) ([]*model_product.ProductPackage, error)
	Save(productID int, packages []dto_product.ProductPackageRequest) error
	DeleteOne(id, productID int) error
}

type productPackageRepo struct {
	db *gorm.DB
}

func NewProductPackageRepo(db *gorm.DB) ProductPackageRepo {
	return &productPackageRepo{db: db}
}
