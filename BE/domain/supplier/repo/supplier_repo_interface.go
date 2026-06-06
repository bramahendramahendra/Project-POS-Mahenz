package repo

import (
	dto "pos_api/domain/supplier/dto"
	model "pos_api/domain/supplier/model"

	"gorm.io/gorm"
)

type (
	SupplierRepo interface {
		GetAll(req *dto.SupplierListRequest) ([]*model.Supplier, int64, error)
		GetOptions() ([]*dto.SupplierOptionResponse, error)
		GetByID(id int) (*model.Supplier, error)
		Create(req *dto.CreateSupplierRequest, code string) (int64, error)
		Update(req *dto.UpdateSupplierRequest) error
		Delete(req *dto.DeleteSupplierRequest) error
		ToggleStatus(req *dto.ToggleStatusSupplierRequest) error

		GetPurchaseHistory(supplierID int) ([]dto.SupplierPurchaseItem, error)
		GetReturnHistory(supplierID int) ([]dto.SupplierReturnHistoryItem, error)
		CheckCodeExists(code string) (bool, error)
		CheckNameExists(name string, excludeID int) (bool, error)
		GetCount() (int, error)
		CountPurchasesBySupplier(supplierID int) (int, error)
		CountActiveDebtBySupplier(supplierID int) (int, error)

		GetDB() *gorm.DB
	}

	supplierRepo struct {
		db *gorm.DB
	}
)

func NewSupplierRepo(db *gorm.DB) *supplierRepo {
	return &supplierRepo{db: db}
}

func (r *supplierRepo) GetDB() *gorm.DB {
	return r.db
}
