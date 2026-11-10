package segment

import (
	access_repo "pos_api/domain/access/repo"
	access_service "pos_api/domain/access/service"
	role_repo "pos_api/domain/role/repo"
	pkgdatabase "pos_api/pkg/database"
)

// newAccessService membuat instance AccessService yang dipakai bersama
// oleh semua segment untuk PermissionMiddleware.
func newAccessService() access_service.AccessServiceInterface {
	return access_service.NewAccessService(
		access_repo.NewAccessRepo(pkgdatabase.DB),
		role_repo.NewRoleRepo(pkgdatabase.DB),
	)
}
