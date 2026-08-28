package repo

import (
	dto "pos_api/domain/stock_reconciliation/dto"

	"gorm.io/gorm"
)

type (
	StockReconciliationRepoInterface interface {
		GetList(req *dto.ListRequest) ([]*dto.ListItem, int64, error)
		GetSummary() (*dto.SummaryResponse, error)
		GetDetail(productID int) (*dto.DetailResponse, error)
		Adjust(req *dto.AdjustRequest) (*dto.DetailResponse, error)
		MarkReviewed(req *dto.MarkReviewedRequest) (*dto.DetailResponse, error)

		GetDB() *gorm.DB
	}

	stockReconciliationRepo struct {
		db *gorm.DB
	}
)

func NewStockReconciliationRepo(db *gorm.DB) *stockReconciliationRepo {
	return &stockReconciliationRepo{db: db}
}

func (r *stockReconciliationRepo) GetDB() *gorm.DB {
	return r.db
}
