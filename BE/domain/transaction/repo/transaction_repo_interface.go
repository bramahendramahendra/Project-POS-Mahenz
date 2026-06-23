package repo

import (
	"pos_api/domain/transaction/dto"
	"pos_api/domain/transaction/model"

	"gorm.io/gorm"
)

type (
	TransactionRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error)
		GetByID(id int) (*dto.TransactionResponse, error)
		Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error)
		Void(id, userID int) error
		GetItems(transactionID int) ([]model.TransactionItem, error)
		UpdateFromSync(id int, data map[string]interface{}) error
		ReturnStockForRejectSync(transactionID, resolvedBy int) error
		ApplySyncTransaction(payload string, localID string) (int, error)
	}

	transactionRepo struct {
		db *gorm.DB
	}
)

func NewTransactionRepo(db *gorm.DB) *transactionRepo {
	return &transactionRepo{db: db}
}
