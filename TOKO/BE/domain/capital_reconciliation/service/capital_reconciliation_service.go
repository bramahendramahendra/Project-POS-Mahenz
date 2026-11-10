package service

import (
	dto "pos_api/domain/capital_reconciliation/dto"
)

func (s *capitalReconciliationService) GetList(req *dto.ListRequest) (data []*dto.ListItem, total int64, err error) {
	return s.repo.GetList(req)
}

func (s *capitalReconciliationService) GetSummary() (data *dto.SummaryResponse, err error) {
	return s.repo.GetSummary()
}

func (s *capitalReconciliationService) Resolve(req *dto.ResolveRequest) error {
	return s.repo.Resolve(req)
}

func (s *capitalReconciliationService) Skip(req *dto.SkipRequest) error {
	return s.repo.Skip(req)
}
