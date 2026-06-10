package service

import (
	dto "pos_api/domain/access/dto"
	"pos_api/errors"
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

	err = s.repo.SetRoleAccess(req.RoleID, req.Accesses)
	return err
}
