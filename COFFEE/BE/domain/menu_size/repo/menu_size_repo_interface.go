package repo

import (
	dto "coffee_api/domain/menu_size/dto"
	model "coffee_api/domain/menu_size/model"

	"gorm.io/gorm"
)

type (
	MenuSizeRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.MenuSize, error)
		GetOptions() ([]*model.MenuSize, error)
		GetByID(id int) (*model.MenuSize, error)
		NameExists(name string, excludeID int) (bool, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	menuSizeRepo struct {
		db *gorm.DB
	}
)

func NewMenuSizeRepo(db *gorm.DB) *menuSizeRepo {
	return &menuSizeRepo{db: db}
}
