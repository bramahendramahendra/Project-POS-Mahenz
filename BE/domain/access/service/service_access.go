package service

import (
	"fmt"
	dto "pos_api/domain/access/dto"
	"pos_api/errors"
	log_helper "pos_api/helper/log"
	"pos_api/pkg/permcache"
)

func (s *accessService) GetByRoleID(roleID int) (data []*dto.RoleMenuAccessItem, err error) {
	exists, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return data, err
	}
	if exists == nil {
		return data, &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}

	dataDB, err := s.repo.GetByRoleID(roleID)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.RoleMenuAccessItem{
			MenuID:    v.MenuID,
			KeyName:   v.KeyName,
			Label:     v.Label,
			ParentID:  v.ParentID,
			CanView:   v.CanView,
			CanCreate: v.CanCreate,
			CanEdit:   v.CanEdit,
			CanDelete: v.CanDelete,
		})
	}

	return data, nil
}

func (s *accessService) SetRoleAccess(req *dto.SetRoleAccessRequest) (err error) {
	exists, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}

	if err = s.repo.SetRoleAccess(req.RoleID, req.Accesses); err != nil {
		return err
	}

	// Invalidate cache — load ulang semua menu key role ini lalu hapus dari cache.
	rows, _ := s.repo.GetByRoleName(exists.Name)
	menuKeys := make([]string, 0, len(rows))
	for _, row := range rows {
		menuKeys = append(menuKeys, row.KeyName)
	}
	permcache.InvalidateRole(exists.Name, menuKeys)
	s.logInvalidate(exists.Name, menuKeys)
	return nil
}

func (s *accessService) GetPermission(roleName, menuKey string) (permcache.Permission, error) {
	endpoint := fmt.Sprintf("/permission/%s/%s", roleName, menuKey)

	// Cek cache dulu.
	if perm, ok := permcache.Get(roleName, menuKey); ok {
		entry := log_helper.FromBackground("GetPermission", "permcache", "[permcache] HIT")
		entry.Method = "CACHE"
		entry.Endpoint = endpoint
		entry.Data = map[string]any{"role": roleName, "menu": menuKey, "perm": perm}
		log_helper.LogDebug(entry)
		return perm, nil
	}

	// Cache miss — load semua permission role dari DB lalu populate cache.
	entry := log_helper.FromBackground("GetPermission", "permcache", "[permcache] MISS — query DB")
	entry.Method = "CACHE"
	entry.Endpoint = endpoint
	entry.Data = map[string]any{"role": roleName, "menu": menuKey}
	log_helper.LogDebug(entry)

	rows, err := s.repo.GetByRoleName(roleName)
	if err != nil {
		return permcache.Permission{}, err
	}

	var result permcache.Permission
	for _, row := range rows {
		p := permcache.Permission{
			CanView:   row.CanView,
			CanCreate: row.CanCreate,
			CanEdit:   row.CanEdit,
			CanDelete: row.CanDelete,
		}
		permcache.Set(roleName, row.KeyName, p)
		if row.KeyName == menuKey {
			result = p
		}
	}
	return result, nil
}

func (s *accessService) logInvalidate(roleName string, menuKeys []string) {
	entry := log_helper.FromBackground("InvalidateRole", "permcache", "[permcache] INVALIDATE")
	entry.Method = "CACHE"
	entry.Endpoint = fmt.Sprintf("/permission/%s/*", roleName)
	entry.Data = map[string]any{"role": roleName, "invalidated_keys": menuKeys}
	log_helper.LogDebug(entry)
}
