package repo

import (
	dto "pos_api/domain/stock_mutation/dto"

	"gorm.io/gorm"
)

type (
	StockMutationRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.StockMutationResponse, int64, error)
		GetByProduct(productID int) ([]*dto.StockMutationByProductResponse, error)

		GetDB() *gorm.DB
	}

	stockMutationRepo struct {
		db *gorm.DB
	}
)

func NewStockMutationRepo(db *gorm.DB) *stockMutationRepo {
	return &stockMutationRepo{db: db}
}

func (r *stockMutationRepo) GetDB() *gorm.DB {
	return r.db
}
