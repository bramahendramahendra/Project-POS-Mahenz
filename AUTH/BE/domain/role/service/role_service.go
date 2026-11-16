package service

import (
	dto "auth_api/domain/role/dto"
)

func (s *roleService) GetActiveOptions() ([]*dto.RoleOptionResponse, error) {
	return s.repo.GetActiveOptions()
}
