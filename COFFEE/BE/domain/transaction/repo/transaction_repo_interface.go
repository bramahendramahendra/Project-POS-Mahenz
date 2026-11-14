package repo

import (
	"coffee_api/domain/transaction/dto"
	"coffee_api/domain/transaction/model"

	"gorm.io/gorm"
)

type TransactionRepoInterface interface {
	GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error)
	GetByID(id int) (*dto.TransactionResponse, error)
	Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error)
	Void(id, userID int, wasteItemIDs map[int]bool, reason string) error
	GetItems(transactionID int) ([]model.TransactionItem, error)

	// Data menu dan meja untuk pesanan (tahap 4.7a).
	GetSaleVariant(variantID int) (*model.SaleVariant, error)
	GetMenuModifierOptions(menuItemID int) ([]model.SaleModifierOption, error)
	GetTableForOrder(id int) (name string, active bool, found bool, err error)

	GetDB() *gorm.DB
	WithTx(tx *gorm.DB) TransactionRepoInterface
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *transactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *transactionRepo) WithTx(tx *gorm.DB) TransactionRepoInterface {
	return &transactionRepo{db: tx}
}
