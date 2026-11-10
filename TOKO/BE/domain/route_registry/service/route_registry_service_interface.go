package service

import (
	dto "pos_api/domain/route_registry/dto"
	repo "pos_api/domain/route_registry/repo"
)

type (
	RouteRegistryServiceInterface interface {
		GetOptions() ([]*dto.RouteOptionResponse, error)
		IsValidPath(path string) (bool, error)
	}

	routeRegistryService struct {
		repo repo.RouteRegistryRepoInterface
	}
)

func NewRouteRegistryService(repo repo.RouteRegistryRepoInterface) *routeRegistryService {
	return &routeRegistryService{repo: repo}
}
