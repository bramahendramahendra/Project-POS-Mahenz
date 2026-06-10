package repo

import (
	dto "pos_api/domain/customer/dto"
	model "pos_api/domain/customer/model"

	"gorm.io/gorm"
)

type (
	CustomerRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Customer, int64, error)
		GetOptions() ([]*model.Customer, error)
		GetByID(id int) (*model.Customer, error)
		GetCount() (int, error)
		CountActiveReceivables(customerID int) (int, error)
		Create(req *dto.CreateRequest, code string) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error

		GetDB() *gorm.DB
	}

	customerRepo struct {
		db *gorm.DB
	}
)

func NewCustomerRepo(db *gorm.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) GetDB() *gorm.DB {
	return r.db
}
