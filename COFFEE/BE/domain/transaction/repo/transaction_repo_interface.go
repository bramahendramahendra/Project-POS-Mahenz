package repo

import (
	"time"

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

	// Open bill (tahap 4.8).
	LockBill(id int) (*LockedBill, error)
	LockTableForOrder(tableID, excludeBillID int) (name string, active, found bool, occupiedBy int, err error)
	ListOpen() ([]dto.OpenBillSummary, error)
	GetItem(billID, itemID int) (*model.TransactionItem, error)
	AddOpenItems(bill *LockedBill, items []dto.CreateTransactionItemRequest, userID int) error
	RemoveOpenItem(bill *LockedBill, item *model.TransactionItem, qty float64, action, reason string, userID int) error
	MoveOpenTable(billID, tableID int, tableName string) error
	MergeOpen(targetID, sourceID int) error
	CancelOpen(bill *LockedBill, wasteItemIDs map[int]bool, reason string, userID int) error
	PayOpen(billID int, req *dto.OpenPayRequest, ownerID, settledBy int, txDate, now time.Time) error

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
