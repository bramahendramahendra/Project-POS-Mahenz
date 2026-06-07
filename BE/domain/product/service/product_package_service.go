package service

import (
	dto "pos_api/domain/product/dto"
	repo "pos_api/domain/product/repo"
	"pos_api/errors"
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

func (s *productPackageService) GetByProduct(productID int) (data []*dto.ProductPackageResponse, err error) {
	if err = s.checkProductExists(productID); err != nil {
		return
	}
	data, err = s.repo.GetByProduct(productID)
	return
}

func (s *productPackageService) Save(req *dto.SaveProductPackagesRequest) error {
	if err := s.checkProductExists(req.ProductID); err != nil {
		return err
	}
	return s.repo.Save(req.ProductID, req.Packages)
}

func (s *productPackageService) DeleteOne(req *dto.PackageIDUriRequest) error {
	if err := s.checkProductExists(req.ID); err != nil {
		return err
	}
	return s.repo.DeleteOne(req.PackageID, req.ID)
}

func (s *productPackageService) checkProductExists(productID int) error {
	p, err := s.prodRepo.GetByID(productID)
	if err != nil {
		return err
	}
	if p == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}
	return nil
}
