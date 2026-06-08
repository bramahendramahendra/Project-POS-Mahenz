package repo_product

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"

	"gorm.io/gorm"
)

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
