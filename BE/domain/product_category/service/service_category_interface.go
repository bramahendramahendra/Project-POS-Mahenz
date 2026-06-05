package service

import (
	dto "pos_api/domain/product_category/dto"
	repo "pos_api/domain/product_category/repo"
)

type (
	CategoryServiceInterface interface {
		GetAll() (data []dto.CategoryResponse, err error)
		GetByID(id int) (data dto.CategoryResponse, err error)
		Create(req *dto.CreateCategoryRequest) (data dto.CategoryResponse, err error)
		Update(id int, req *dto.UpdateCategoryRequest) error
		Delete(id int) error
		ToggleStatus(id int) error
	}

	categoryService struct {
		repo repo.CategoryRepoInterface
	}
)

func NewCategoryService(repo repo.CategoryRepoInterface) CategoryServiceInterface {
	return &categoryService{repo: repo}
}
