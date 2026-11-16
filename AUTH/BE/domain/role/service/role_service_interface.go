package service

import (
	dto "auth_api/domain/role/dto"
	repo "auth_api/domain/role/repo"
)

type (
	RoleServiceInterface interface {
		GetActiveOptions() ([]*dto.RoleOptionResponse, error)
	}

	roleService struct {
		repo repo.RoleRepoInterface
	}
)

func NewRoleService(repo repo.RoleRepoInterface) *roleService {
	return &roleService{repo: repo}
}
