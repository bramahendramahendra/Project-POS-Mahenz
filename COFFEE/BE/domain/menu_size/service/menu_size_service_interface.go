package service

import (
	dto "coffee_api/domain/menu_size/dto"
	repo "coffee_api/domain/menu_size/repo"
)

type (
	MenuSizeServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.MenuSizeResponse, error)
		GetOptions() ([]dto.OptionResponse, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	menuSizeService struct {
		repo repo.MenuSizeRepoInterface
	}
)

func NewMenuSizeService(repo repo.MenuSizeRepoInterface) *menuSizeService {
	return &menuSizeService{repo: repo}
}
