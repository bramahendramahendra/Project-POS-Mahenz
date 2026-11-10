package repo

import (
	"pos_api/domain/customer_balance/dto"
	"pos_api/domain/customer_balance/model"

	"gorm.io/gorm"
)

type (
	CustomerBalanceRepoInterface interface {
		GetHistory(req *dto.HistoryRequest) ([]*dto.MutationResponse, int64, error)
		GetMutationByRef(refType string, refID int, mutationType string) (*model.CustomerBalanceMutation, error)
		Topup(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error)
		Deduct(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error)
		Refund(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error)
		Adjust(customerID int, amount float64, notes string, userID int) (float64, error)
		GetBalance(customerID int) (float64, error)
		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) CustomerBalanceRepoInterface
	}

	customerBalanceRepo struct {
		db *gorm.DB
	}
)

func NewCustomerBalanceRepo(db *gorm.DB) CustomerBalanceRepoInterface {
	return &customerBalanceRepo{db: db}
}

func (r *customerBalanceRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *customerBalanceRepo) WithTx(tx *gorm.DB) CustomerBalanceRepoInterface {
	return &customerBalanceRepo{db: tx}
}
