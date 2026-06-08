package service

import (
	dto "pos_api/domain/product/dto"
	repo "pos_api/domain/product/repo"
)

type (
	ProductPriceServiceInterface interface {
		GetByProduct(productID int) (data []*dto.ProductPriceResponse, err error)
		Save(req *dto.SaveProductPricesRequest) error
	}

	productPriceService struct {
		repo     repo.ProductRepo
		prodRepo repo.ProductRepo
	}
)

func NewProductPriceService(repo repo.ProductRepo, prodRepo repo.ProductRepo) *productPriceService {
	return &productPriceService{repo: repo, prodRepo: prodRepo}
}
