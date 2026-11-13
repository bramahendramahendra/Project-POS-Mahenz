package service

import (
	dto "coffee_api/domain/access/dto"
	repo "coffee_api/domain/access/repo"
	role_repo "coffee_api/domain/role/repo"
	"coffee_api/pkg/permcache"
)

type (
	AccessServiceInterface interface {
		GetByRoleID(roleID int) (data []*dto.RoleMenuAccessItem, err error)
		SetRoleAccess(req *dto.SetRoleAccessRequest) (err error)
		// GetPermission mengambil permission satu menu untuk satu role.
		// Hasil di-cache; cache di-invalidate saat SetRoleAccess dipanggil.
		GetPermission(roleName, menuKey string) (permcache.Permission, error)
	}

	accessService struct {
		repo     repo.AccessRepoInterface
		roleRepo role_repo.RoleRepoInterface
	}
)

func NewAccessService(repo repo.AccessRepoInterface, roleRepo role_repo.RoleRepoInterface) *accessService {
	return &accessService{repo: repo, roleRepo: roleRepo}
}
