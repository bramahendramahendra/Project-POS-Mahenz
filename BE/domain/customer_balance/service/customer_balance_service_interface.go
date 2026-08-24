package service

import (
	customer_repo "pos_api/domain/customer/repo"
	"pos_api/domain/customer_balance/dto"
	"pos_api/domain/customer_balance/repo"
)

type (
	CustomerBalanceServiceInterface interface {
		Topup(req *dto.TopupRequest) (*dto.BalanceResponse, error)
		Refund(req *dto.RefundRequest) (*dto.BalanceResponse, error)
		GetHistory(req *dto.HistoryRequest) ([]*dto.MutationResponse, int64, error)
		SaveToBalance(req *dto.SaveToBalanceRequest) (*dto.BalanceResponse, error)
	}

	customerBalanceService struct {
		repo         repo.CustomerBalanceRepoInterface
		customerRepo customer_repo.CustomerRepoInterface
	}
)

func NewCustomerBalanceService(repo repo.CustomerBalanceRepoInterface, customerRepo customer_repo.CustomerRepoInterface) CustomerBalanceServiceInterface {
	return &customerBalanceService{repo: repo, customerRepo: customerRepo}
}
