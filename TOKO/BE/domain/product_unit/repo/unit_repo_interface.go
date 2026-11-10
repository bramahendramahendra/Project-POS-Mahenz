package repo

import (
	dto "pos_api/domain/product_unit/dto"
	model "pos_api/domain/product_unit/model"

	"gorm.io/gorm"
)

type (
	UnitRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Unit, int64, error)
		GetOptions() ([]*model.UnitOption, error)
		GetByID(id int) (*model.Unit, error)
		Create(req *dto.CreateRequest) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error

		CheckNameExists(name string, excludeID int) (bool, error)
		CheckAbbreviationExists(abbreviation string, excludeID int) (bool, error)
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
