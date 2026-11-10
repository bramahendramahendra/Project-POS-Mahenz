package repo

import (
	model "pos_api/domain/route_registry/model"
)

func (r *routeRegistryRepo) GetActiveOptions() ([]*model.RouteRegistry, error) {
	var data []*model.RouteRegistry
	err := r.db.Where("is_active = 1").Order("path ASC").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *routeRegistryRepo) ExistsActivePath(path string) (bool, error) {
	var count int64
	err := r.db.Model(&model.RouteRegistry{}).Where("path = ? AND is_active = 1", path).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
