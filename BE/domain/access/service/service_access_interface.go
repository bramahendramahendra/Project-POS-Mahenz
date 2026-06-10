package service

import (
	dto "pos_api/domain/access/dto"
	repo "pos_api/domain/access/repo"
	role_repo "pos_api/domain/role/repo"
)

type (
	AccessServiceInterface interface {
		GetByRoleID(roleID int) (data []*dto.RoleMenuAccessItem, err error)
		SetRoleAccess(req *dto.SetRoleAccessRequest) (err error)
	}

	accessService struct {
		repo     repo.AccessRepoInterface
		roleRepo role_repo.RoleRepoInterface
	}
)

func NewAccessService(repo repo.AccessRepoInterface, roleRepo role_repo.RoleRepoInterface) *accessService {
	return &accessService{repo: repo, roleRepo: roleRepo}
}
