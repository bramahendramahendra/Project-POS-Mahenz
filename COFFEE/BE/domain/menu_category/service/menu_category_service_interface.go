package service

import (
	dto "coffee_api/domain/menu_category/dto"
	repo "coffee_api/domain/menu_category/repo"
)

type (
	MenuCategoryServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.MenuCategoryResponse, error)
		GetOptions() ([]dto.OptionResponse, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	menuCategoryService struct {
		repo repo.MenuCategoryRepoInterface
	}
)

func NewMenuCategoryService(repo repo.MenuCategoryRepoInterface) *menuCategoryService {
	return &menuCategoryService{repo: repo}
}
