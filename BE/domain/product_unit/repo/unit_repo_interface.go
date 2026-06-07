package repo

import (
	dto "pos_api/domain/product_unit/dto"
	model "pos_api/domain/product_unit/model"

	"gorm.io/gorm"
)

type (
	UnitRepoInterface interface {
		GetAll(req *dto.UnitListRequest) ([]*model.Unit, int64, error)
		GetOptions() ([]*model.UnitOption, error)
		GetByID(id int) (*model.Unit, error)
		Create(req *dto.CreateUnitRequest) (int64, error)
		Update(req *dto.UpdateUnitRequest) error
		Delete(req *dto.DeleteUnitRequest) error
		ToggleStatus(req *dto.ToggleStatusUnitRequest) error

		CheckNameExists(name string, excludeID int) (bool, error)
		CountProductUnitsByUnit(unitID int) (int, error)
		CountActiveProductsByUnit(unitID int) (int, error)

		GetDB() *gorm.DB
	}

	unitRepo struct {
		db *gorm.DB
	}
)

func NewUnitRepo(db *gorm.DB) *unitRepo {
	return &unitRepo{db: db}
}

func (r *unitRepo) GetDB() *gorm.DB {
	return r.db
}
