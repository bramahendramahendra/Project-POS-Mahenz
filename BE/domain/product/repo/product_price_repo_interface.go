package repo_product

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"

	"gorm.io/gorm"
)

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
