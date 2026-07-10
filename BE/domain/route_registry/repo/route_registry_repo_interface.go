package repo

import (
	model "pos_api/domain/route_registry/model"

	"gorm.io/gorm"
)

type (
	RouteRegistryRepoInterface interface {
		GetActiveOptions() ([]*model.RouteRegistry, error)
		ExistsActivePath(path string) (bool, error)
	}

	routeRegistryRepo struct {
		db *gorm.DB
	}
)

func NewRouteRegistryRepo(db *gorm.DB) *routeRegistryRepo {
	return &routeRegistryRepo{db: db}
}
