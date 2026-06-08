package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetPricesByProduct(productID int) (data []*dto.ProductPriceResponse, err error) {
	exists, err := s.repo.GetByID(productID)
	if err != nil {
		return data, err
	}
	if exists == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	prices, err := s.repo.GetPricesByProduct(productID)
	if err != nil {
		return data, err
	}

	for _, v := range prices {
		data = append(data, &dto.ProductPriceResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			TierName:  v.TierName,
			MinQty:    v.MinQty,
			Price:     v.Price,
		})
	}

	return data, nil
}

func (s *productService) SavePrices(req *dto.SaveProductPricesRequest) (err error) {
	exists, err := s.repo.GetByID(req.ProductID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	return s.repo.SavePrices(req.ProductID, req.Prices)
}
