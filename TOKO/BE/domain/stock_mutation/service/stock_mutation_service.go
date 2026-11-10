package service

import (
	dto "pos_api/domain/stock_mutation/dto"
)

func (s *stockMutationService) GetAll(req *dto.GetAllRequest) (data []*dto.StockMutationResponse, total int64, err error) {
	return s.repo.GetAll(req)
}

func (s *stockMutationService) GetByProduct(productID int) (data []*dto.StockMutationByProductResponse, err error) {
	return s.repo.GetByProduct(productID)
}
