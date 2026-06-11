package service

import (
	"pos_api/domain/transaction/dto"
	repo "pos_api/domain/transaction/repo"
)

type (
	TransactionServiceInterface interface {
		GetAll(filter *dto.TransactionFilter) ([]*dto.TransactionResponse, int, error)
		GetByID(id int) (*dto.TransactionResponse, error)
		Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error)
		Void(id, userID int) error
	}

	transactionService struct {
		repo repo.TransactionRepoInterface
	}
)

func NewTransactionService(r repo.TransactionRepoInterface) *transactionService {
	return &transactionService{repo: r}
}
