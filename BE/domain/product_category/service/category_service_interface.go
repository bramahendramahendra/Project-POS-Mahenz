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
		Update(req *dto.UpdateCategoryRequest) (err error)
		Delete(req *dto.DeleteCategoryRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusCategoryRequest) (err error)
	}

	categoryService struct {
		repo repo.CategoryRepoInterface
	}
)

func NewCategoryService(repo repo.CategoryRepoInterface) *categoryService {
	return &categoryService{repo: repo}
}
