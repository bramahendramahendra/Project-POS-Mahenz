package service

import (
	dto "pos_api/domain/product_category/dto"
	repo "pos_api/domain/product_category/repo"
)

type (
	CategoryServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.CategoryResponse, total int64, err error)
		GetOptions() (data []dto.GetOptionResponse, err error)
		GetByID(id int) (data dto.CategoryResponse, err error)
		Create(req *dto.CreateRequest) (data dto.CategoryResponse, err error)
		Update(req *dto.UpdateRequest) (data dto.CategoryResponse, err error)
		Delete(req *dto.DeleteRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusRequest) (err error)
	}

	categoryService struct {
		repo repo.CategoryRepoInterface
	}
)

func NewCategoryService(repo repo.CategoryRepoInterface) *categoryService {
	return &categoryService{repo: repo}
}
