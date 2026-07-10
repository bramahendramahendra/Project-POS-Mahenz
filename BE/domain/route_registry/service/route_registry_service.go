package service

import (
	dto "pos_api/domain/route_registry/dto"
)

func (s *routeRegistryService) GetOptions() ([]*dto.RouteOptionResponse, error) {
	data, err := s.repo.GetActiveOptions()
	if err != nil {
		return nil, err
	}
	result := make([]*dto.RouteOptionResponse, 0, len(data))
	for _, d := range data {
		result = append(result, &dto.RouteOptionResponse{Path: d.Path, Label: d.Label})
	}
	return result, nil
}

func (s *routeRegistryService) IsValidPath(path string) (bool, error) {
	return s.repo.ExistsActivePath(path)
}
