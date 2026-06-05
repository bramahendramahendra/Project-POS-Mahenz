package repo

import (
	dto "pos_api/domain/product_category/dto"
	model "pos_api/domain/product_category/model"

	"gorm.io/gorm"
)

type (
	CategoryRepoInterface interface {
		GetAll() ([]*model.Category, error)
		GetByID(id int) (*model.Category, error)
		GetByName(name string) (*model.Category, error)
		CheckNameExists(name string, excludeID int) (bool, error)
		CheckCodeExists(code string) (bool, error)
		CountProductsByCategory(categoryID int) (int, error)
		CountActiveProductsByCategory(categoryID int) (int, error)
		Create(req *dto.CreateCategoryRequest) (int64, error)
		Update(req *dto.UpdateCategoryRequest) error
		Delete(req *dto.DeleteCategoryRequest) error
		ToggleStatus(req *dto.ToggleStatusCategoryRequest) error
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
