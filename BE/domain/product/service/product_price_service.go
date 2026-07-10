package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetPricesByProduct(productID int) (data []*dto.PriceResponse, err error) {
	exists, err := s.repo.GetByID(productID)
	if err != nil {
		return data, err
	}
	if exists == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	dataDB, err := s.repo.GetPricesByProduct(productID)
	if err != nil {
		return data, err
	}

	data = make([]*dto.PriceResponse, 0, len(dataDB))
	for _, v := range dataDB {
		data = append(data, &dto.PriceResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			TierName:  v.TierName,
			MinQty:    v.MinQty,
			Price:     v.Price,
		})
	}

	return data, nil
}

func (s *productService) SavePrices(req *dto.SavePriceRequest) (err error) {
	exists, err := s.repo.GetByID(req.ProductID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	err = s.repo.SavePrices(req.ProductID, req.Prices)
	return err
}
