package service

import (
	dto "coffee_api/domain/menu_item/dto"
	repo "coffee_api/domain/menu_item/repo"
)

type (
	MenuItemServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.MenuItemResponse, int64, error)
		GetByID(id int) (*dto.MenuItemDetailResponse, error)
		Create(req *dto.SaveRequest, userID int) (*dto.MenuItemDetailResponse, error)
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		SetSoldOut(req *dto.SoldOutRequest) error
		Delete(id int) error
		Catalog() (*dto.CatalogResponse, error)
	}

	menuItemService struct {
		repo repo.MenuItemRepoInterface
	}
)

func NewMenuItemService(repo repo.MenuItemRepoInterface) *menuItemService {
	return &menuItemService{repo: repo}
}
