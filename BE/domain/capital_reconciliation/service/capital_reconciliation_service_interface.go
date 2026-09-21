package service

import (
	dto "pos_api/domain/capital_reconciliation/dto"
	repo "pos_api/domain/capital_reconciliation/repo"
)

type (
	CapitalReconciliationServiceInterface interface {
		GetList(req *dto.ListRequest) (data []*dto.ListItem, total int64, err error)
		GetSummary() (data *dto.SummaryResponse, err error)
		Resolve(req *dto.ResolveRequest) error
		Skip(req *dto.SkipRequest) error
	}

	capitalReconciliationService struct {
		repo repo.CapitalReconciliationRepoInterface
	}
)

func NewCapitalReconciliationService(repo repo.CapitalReconciliationRepoInterface) *capitalReconciliationService {
	return &capitalReconciliationService{repo: repo}
}
