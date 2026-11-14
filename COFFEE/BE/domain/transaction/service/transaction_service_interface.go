package service

import (
	cash_drawer_repo "coffee_api/domain/cash_drawer/repo"
	product_repo "coffee_api/domain/product/repo"
	"coffee_api/domain/transaction/dto"
	repo "coffee_api/domain/transaction/repo"
)

type (
	TransactionServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error)
		GetByID(id int, requestingUserID int, role string) (*dto.TransactionResponse, error)
		Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error)
		Void(req *dto.VoidRequest, userID int) error

		// Open bill (tahap 4.8).
		ListOpen() ([]dto.OpenBillSummary, error)
		AddOpenItems(id int, req *dto.OpenAddItemsRequest, userID int) (*dto.TransactionResponse, error)
		RemoveOpenItem(id int, req *dto.OpenRemoveItemRequest, userID int) (*dto.TransactionResponse, error)
		MoveOpenTable(id int, req *dto.OpenMoveTableRequest) (*dto.TransactionResponse, error)
		MergeOpen(targetID int, req *dto.OpenMergeRequest) (*dto.TransactionResponse, error)
		CancelOpen(id int, req *dto.OpenCancelRequest, userID int) error
		PayOpen(id int, req *dto.OpenPayRequest, userID int) (*dto.TransactionResponse, error)
	}

	transactionService struct {
		repo           repo.TransactionRepoInterface
		cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface
		productRepo    product_repo.ProductRepoInterface
	}
)

func NewTransactionService(r repo.TransactionRepoInterface, cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface, productRepo product_repo.ProductRepoInterface) *transactionService {
	return &transactionService{repo: r, cashDrawerRepo: cashDrawerRepo, productRepo: productRepo}
}
