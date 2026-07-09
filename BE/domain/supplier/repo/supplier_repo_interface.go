package repo

import (
	dto "pos_api/domain/supplier/dto"
	model "pos_api/domain/supplier/model"

	"gorm.io/gorm"
)

type (
	SupplierRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Supplier, int64, error)
		GetOptions(search string) ([]*model.SupplierOption, error)
		GetByID(id int) (*model.Supplier, error)
		Create(req *dto.CreateRequest, code string) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error

		GetPurchaseHistory(supplierID int) ([]*model.SupplierPurchase, error)
		GetReturnHistory(supplierID int) ([]*model.SupplierReturn, error)
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
