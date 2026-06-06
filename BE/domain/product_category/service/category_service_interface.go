package service

import (
	dto "pos_api/domain/product_category/dto"
	repo "pos_api/domain/product_category/repo"
)

type (
	CategoryServiceInterface interface {
		GetAll(req *dto.CategoryListRequest) (data []dto.CategoryResponse, total int64, err error)
		GetOptions() (data []dto.CategoryOptionResponse, err error)
		GetByID(id int) (data dto.CategoryResponse, err error)
		Create(req *dto.CreateCategoryRequest) (data dto.CategoryResponse, err error)
		Update(req *dto.UpdateCategoryRequest) (data dto.CategoryResponse, err error)
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
