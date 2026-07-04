package service

import (
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	product_repo "pos_api/domain/product/repo"
	"pos_api/domain/transaction/dto"
	repo "pos_api/domain/transaction/repo"
)

type (
	TransactionServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error)
		GetByID(id int) (*dto.TransactionResponse, error)
		Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error)
		Void(req *dto.VoidRequest, userID int) error
	}

	transactionService struct {
		repo           repo.TransactionRepoInterface
		cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface
		productRepo    product_repo.ProductRepo
	}
)

func NewTransactionService(r repo.TransactionRepoInterface, cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface, productRepo product_repo.ProductRepo) *transactionService {
	return &transactionService{repo: r, cashDrawerRepo: cashDrawerRepo, productRepo: productRepo}
}
