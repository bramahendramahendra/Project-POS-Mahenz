package repo

import (
	dto "coffee_api/domain/dining_table/dto"
	model "coffee_api/domain/dining_table/model"

	"gorm.io/gorm"
)

type (
	DiningTableRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.DiningTable, error)
		GetOptions() ([]*model.DiningTable, error)
		GetByID(id int) (*model.DiningTable, error)
		NameExists(name string, excludeID int) (bool, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	diningTableRepo struct {
		db *gorm.DB
	}
)

func NewDiningTableRepo(db *gorm.DB) *diningTableRepo {
	return &diningTableRepo{db: db}
}
