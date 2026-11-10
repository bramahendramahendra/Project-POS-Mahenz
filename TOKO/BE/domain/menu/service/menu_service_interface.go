package service

import (
	dto "pos_api/domain/menu/dto"
	repo "pos_api/domain/menu/repo"
	route_registry_service "pos_api/domain/route_registry/service"
)

type (
	MenuServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.MenuResponse, int64, error)
		GetRootOptions() ([]*dto.MenuOptionResponse, error)
		GetByID(id int) (*dto.MenuResponse, error)
		GetMyMenus(roleName string) ([]dto.MyMenuItem, error)
		Create(req *dto.CreateRequest) (*dto.MenuResponse, error)
		Update(id int, req *dto.UpdateRequest) error
		Delete(id int) error
		Reorder(req *dto.ReorderRequest) error
	}

	menuService struct {
		repo          repo.MenuRepoInterface
		routeRegistry route_registry_service.RouteRegistryServiceInterface
	}
)

func NewMenuService(repo repo.MenuRepoInterface, routeRegistry route_registry_service.RouteRegistryServiceInterface) *menuService {
	return &menuService{repo: repo, routeRegistry: routeRegistry}
}
