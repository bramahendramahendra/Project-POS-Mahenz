package service

import (
	dto "pos_api/domain/finance/dto"
	repo "pos_api/domain/finance/repo"
)

type (
	FinanceServiceInterface interface {
		GetSummary(req *dto.GetSummaryRequest) (*dto.SummaryResponse, error)
		GetCashflow(req *dto.GetCashflowRequest) (data []dto.CashflowItemResponse, total int64, err error)
	}

	financeService struct {
		repo repo.FinanceRepo
	}
)

func NewFinanceService(repo repo.FinanceRepo) *financeService {
	return &financeService{repo: repo}
}
