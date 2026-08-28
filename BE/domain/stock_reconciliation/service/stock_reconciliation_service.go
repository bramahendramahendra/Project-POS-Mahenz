package service

import (
	dto "pos_api/domain/stock_reconciliation/dto"
)

func (s *stockReconciliationService) GetList(req *dto.ListRequest) (data []*dto.ListItem, total int64, err error) {
	return s.repo.GetList(req)
}

func (s *stockReconciliationService) GetSummary() (data *dto.SummaryResponse, err error) {
	return s.repo.GetSummary()
}

func (s *stockReconciliationService) GetDetail(productID int) (data *dto.DetailResponse, err error) {
	return s.repo.GetDetail(productID)
}

func (s *stockReconciliationService) Adjust(req *dto.AdjustRequest) (data *dto.DetailResponse, err error) {
	return s.repo.Adjust(req)
}

func (s *stockReconciliationService) MarkReviewed(req *dto.MarkReviewedRequest) (data *dto.DetailResponse, err error) {
	return s.repo.MarkReviewed(req)
}
