package service

import (
	dto "pos_api/domain/role/dto"
	repo "pos_api/domain/role/repo"
)

type (
	RoleServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.RoleResponse, int64, error)
		GetActiveOptions() ([]*dto.RoleOptionResponse, error)
		GetByID(id int) (*dto.RoleResponse, error)
		Create(req *dto.CreateRequest) (*dto.RoleResponse, error)
		Update(req *dto.UpdateRequest) error
		Delete(id int) error
		ToggleStatus(id int) error
	}

	roleService struct {
		repo repo.RoleRepoInterface
	}
)

func NewRoleService(repo repo.RoleRepoInterface) *roleService {
	return &roleService{repo: repo}
}
