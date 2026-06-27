package service

import (
	dto "pos_api/domain/finance/dto"
)

func (s *financeService) GetSummary(req *dto.GetSummaryRequest) (*dto.SummaryResponse, error) {
	return s.repo.GetSummary(req)
}

func (s *financeService) GetCashflow(req *dto.GetCashflowRequest) ([]dto.CashflowItemResponse, int64, error) {
	return s.repo.GetCashflow(req)
}
