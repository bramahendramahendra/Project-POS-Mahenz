package repo_product

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"
)

type ProductRepo interface {
	GetAll(req *dto_product.ProductListRequest) ([]*dto_product.ProductResponse, int, error)
	GetOptions() ([]*dto_product.ProductOption, error)
	GetByID(id int) (*model_product.Product, error)
	GetByBarcode(barcode string) (*model_product.Product, error)
	Search(keyword string, limit int) ([]*dto_product.ProductSearchResult, error)
	GetLowStock() ([]*dto_product.LowStockProduct, error)
	CheckBarcodeExists(barcode string, excludeID int) (bool, error)
	CheckSkuExists(sku string, excludeID int) (bool, error)
	CountSkuByCategory(categoryID int) (int, error)
	CountTransactionItems(productID int) (int, error)
	Create(req *dto_product.ProductRequest) (int64, error)
	Update(req *dto_product.UpdateProductRequest) error
	Delete(id int) error
	ToggleStatus(id int) error
	UpdateStock(id int, delta float64) error
}
