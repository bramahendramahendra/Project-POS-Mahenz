package service

import (
	dto "coffee_api/domain/finance/dto"
	repo "coffee_api/domain/finance/repo"
)

type (
	FinanceServiceInterface interface {
		GetSummary(req *dto.GetSummaryRequest) (*dto.SummaryResponse, error)
		GetCashflow(req *dto.GetCashflowRequest) (data []dto.CashflowItemResponse, total int64, err error)
	}

	financeService struct {
		repo repo.FinanceRepoInterface
	}
)

func NewFinanceService(repo repo.FinanceRepoInterface) *financeService {
	return &financeService{repo: repo}
}
