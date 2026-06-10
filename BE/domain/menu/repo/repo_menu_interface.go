package repo

import (
	dto "pos_api/domain/menu/dto"
	model "pos_api/domain/menu/model"

	"gorm.io/gorm"
)

type (
	MenuRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Menu, error)
		GetByID(id int) (*model.Menu, error)
		GetByKeyName(keyName string) (*model.Menu, error)
		GetMyMenus(roleName string) ([]*dto.MyMenuItem, error)
		Create(req *dto.CreateRequest) (int64, error)
		Update(id int, req *dto.UpdateRequest) error
		Delete(id int) error
		Reorder(items []dto.ReorderItem) error
	}

	menuRepo struct {
		db *gorm.DB
	}
)

func NewMenuRepo(db *gorm.DB) *menuRepo {
	return &menuRepo{db: db}
}
