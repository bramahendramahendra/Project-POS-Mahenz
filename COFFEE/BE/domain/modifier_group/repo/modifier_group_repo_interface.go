package repo

import (
	dto "coffee_api/domain/modifier_group/dto"
	model "coffee_api/domain/modifier_group/model"

	"gorm.io/gorm"
)

type (
	ModifierGroupRepoInterface interface {
		GetAll(req *dto.GetAllRequest, onlyActive bool) ([]*model.ModifierGroup, error)
		GetByID(id int) (*model.ModifierGroup, error)
		GetOptionsByGroupIDs(groupIDs []int, onlyActive bool) ([]*model.ModifierOption, error)
		NameExists(name string, excludeID int) (bool, error)
		Create(req *dto.SaveRequest) (int64, error)
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error

		InsertOption(groupID int, o *dto.OptionInput, sortOrder int) error
		UpdateOption(groupID int, o *dto.OptionInput, sortOrder int) error
		DeleteOptionsNotIn(groupID int, keepIDs []int) error

		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) ModifierGroupRepoInterface
	}

	modifierGroupRepo struct {
		db *gorm.DB
	}
)

func NewModifierGroupRepo(db *gorm.DB) *modifierGroupRepo {
	return &modifierGroupRepo{db: db}
}

func (r *modifierGroupRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *modifierGroupRepo) WithTx(tx *gorm.DB) ModifierGroupRepoInterface {
	return &modifierGroupRepo{db: tx}
}
