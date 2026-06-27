package repo

import (
	dto "pos_api/domain/finance/dto"

	"gorm.io/gorm"
)

type (
	FinanceRepo interface {
		GetSummary(req *dto.GetSummaryRequest) (*dto.SummaryResponse, error)
		GetCashflow(req *dto.GetCashflowRequest) ([]dto.CashflowItemResponse, int64, error)

		GetDB() *gorm.DB
	}

	financeRepo struct {
		db *gorm.DB
	}
)

func NewFinanceRepo(db *gorm.DB) *financeRepo {
	return &financeRepo{db: db}
}

func (r *financeRepo) GetDB() *gorm.DB {
	return r.db
}
