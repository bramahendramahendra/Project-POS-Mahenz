package service

import (
	dto "pos_api/domain/stock_mutation/dto"
	repo "pos_api/domain/stock_mutation/repo"
)

type (
	StockMutationServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []*dto.StockMutationResponse, total int64, err error)
		GetByProduct(productID int) (data []*dto.StockMutationByProductResponse, err error)
	}

	stockMutationService struct {
		repo repo.StockMutationRepoInterface
	}
)

func NewStockMutationService(repo repo.StockMutationRepoInterface) *stockMutationService {
	return &stockMutationService{repo: repo}
}
