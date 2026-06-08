package repo

import (
	dto "pos_api/domain/product_category/dto"
	model "pos_api/domain/product_category/model"

	"gorm.io/gorm"
)

type (
	CategoryRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Category, int64, error)
		GetOptions() ([]*model.CategoryOption, error)
		GetByID(id int) (*model.Category, error)
		Create(req *dto.CreateRequest) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error

		GetByName(name string) (*model.Category, error)
		CheckCodeExists(code string) (bool, error)
		CheckNameExists(name string, excludeID int) (bool, error)
		CountProductsByCategory(categoryID int) (int, error)
		CountActiveProductsByCategory(categoryID int) (int, error)

		GetDB() *gorm.DB
	}

	categoryRepo struct {
		db *gorm.DB
	}
)

func NewCategoryRepo(db *gorm.DB) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetDB() *gorm.DB {
	return r.db
}
