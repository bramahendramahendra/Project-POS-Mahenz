package service

import (
	dto "pos_api/domain/stock_reconciliation/dto"
	repo "pos_api/domain/stock_reconciliation/repo"
)

type (
	StockReconciliationServiceInterface interface {
		GetList(req *dto.ListRequest) (data []*dto.ListItem, total int64, err error)
		GetSummary() (data *dto.SummaryResponse, err error)
		GetDetail(productID int) (data *dto.DetailResponse, err error)
		Adjust(req *dto.AdjustRequest) (data *dto.DetailResponse, err error)
		MarkReviewed(req *dto.MarkReviewedRequest) (data *dto.DetailResponse, err error)
	}

	stockReconciliationService struct {
		repo repo.StockReconciliationRepoInterface
	}
)

func NewStockReconciliationService(repo repo.StockReconciliationRepoInterface) *stockReconciliationService {
	return &stockReconciliationService{repo: repo}
}
