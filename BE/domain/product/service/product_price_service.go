package service

import (
	dto "pos_api/domain/product/dto"
	repo "pos_api/domain/product/repo"
	"pos_api/errors"
)

type (
	ProductPriceServiceInterface interface {
		GetByProduct(productID int) (data []*dto.ProductPriceResponse, err error)
		Save(req *dto.SaveProductPricesRequest) error
	}

	productPriceService struct {
		repo     repo.ProductPriceRepo
		prodRepo repo.ProductRepo
	}
)

func NewProductPriceService(repo repo.ProductPriceRepo, prodRepo repo.ProductRepo) *productPriceService {
	return &productPriceService{repo: repo, prodRepo: prodRepo}
}

func (s *productPriceService) GetByProduct(productID int) (data []*dto.ProductPriceResponse, err error) {
	if err = s.checkProductExists(productID); err != nil {
		return
	}
	data, err = s.repo.GetByProduct(productID)
	return
}

func (s *productPriceService) Save(req *dto.SaveProductPricesRequest) error {
	if err := s.checkProductExists(req.ProductID); err != nil {
		return err
	}
	return s.repo.Save(req.ProductID, req.Prices)
}

func (s *productPriceService) checkProductExists(productID int) error {
	p, err := s.prodRepo.GetByID(productID)
	if err != nil {
		return err
	}
	if p == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}
	return nil
}
