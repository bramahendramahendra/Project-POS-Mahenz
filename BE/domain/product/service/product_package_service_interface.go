package service

import (
	dto "pos_api/domain/product/dto"
	repo "pos_api/domain/product/repo"
)

type (
	ProductPackageServiceInterface interface {
		GetByProduct(productID int) (data []*dto.ProductPackageResponse, err error)
		Save(req *dto.SaveProductPackagesRequest) error
		DeleteOne(req *dto.PackageIDUriRequest) error
	}

	productPackageService struct {
		repo     repo.ProductPackageRepo
		prodRepo repo.ProductRepo
	}
)

func NewProductPackageService(repo repo.ProductPackageRepo, prodRepo repo.ProductRepo) *productPackageService {
	return &productPackageService{repo: repo, prodRepo: prodRepo}
}
