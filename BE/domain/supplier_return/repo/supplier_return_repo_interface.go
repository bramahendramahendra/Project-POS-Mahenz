package repo

import (
	dto "pos_api/domain/supplier_return/dto"
	model "pos_api/domain/supplier_return/model"

	"gorm.io/gorm"
)

type (
	SupplierReturnRepo interface {
		GetAll(req *dto.SupplierReturnListRequest) ([]*model.SupplierReturnRow, int64, error)
		GetByID(id int) (*model.SupplierReturnRow, error)
		GetStatus(id int) (string, error)
		GetItems(returnID int) ([]model.SupplierReturnItem, error)
		GetPurchaseDate(purchaseID int) (string, error)
		Create(req *dto.CreateSupplierReturnRequest) (*model.SupplierReturnRow, error)
		UpdateStatus(id int, status, notes string) error
		ApproveWithStockReduction(id int, userID int) error
		Delete(req *dto.GetSupplierReturnByIDRequest) error

		GetDB() *gorm.DB
	}

	supplierReturnRepo struct {
		db *gorm.DB
	}
)

func NewSupplierReturnRepo(db *gorm.DB) *supplierReturnRepo {
	return &supplierReturnRepo{db: db}
}

func (r *supplierReturnRepo) GetDB() *gorm.DB {
	return r.db
}
