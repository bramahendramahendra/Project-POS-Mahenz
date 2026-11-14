package repo

import (
	dto "coffee_api/domain/menu_category/dto"
	model "coffee_api/domain/menu_category/model"

	"gorm.io/gorm"
)

type (
	MenuCategoryRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.MenuCategory, error)
		GetOptions() ([]*model.MenuCategory, error)
		GetByID(id int) (*model.MenuCategory, error)
		NameExists(name string, excludeID int) (bool, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	menuCategoryRepo struct {
		db *gorm.DB
	}
)

func NewMenuCategoryRepo(db *gorm.DB) *menuCategoryRepo {
	return &menuCategoryRepo{db: db}
}
