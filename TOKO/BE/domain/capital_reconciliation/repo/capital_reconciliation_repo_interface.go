package repo

import (
	dto "pos_api/domain/capital_reconciliation/dto"

	"gorm.io/gorm"
)

type (
	CapitalReconciliationRepoInterface interface {
		GetList(req *dto.ListRequest) ([]*dto.ListItem, int64, error)
		GetSummary() (*dto.SummaryResponse, error)
		Resolve(req *dto.ResolveRequest) error
		Skip(req *dto.SkipRequest) error
	}

	capitalReconciliationRepo struct {
		db *gorm.DB
	}
)

func NewCapitalReconciliationRepo(db *gorm.DB) *capitalReconciliationRepo {
	return &capitalReconciliationRepo{db: db}
}
