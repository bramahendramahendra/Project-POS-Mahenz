package service

import (
	dto "pos_api/domain/role/dto"
	model "pos_api/domain/role/model"
	"pos_api/errors"
)

func (s *roleService) GetAll(req *dto.GetAllRequest) ([]*dto.RoleResponse, error) {
	roles, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		result = append(result, toRoleResponse(r))
	}
	return result, nil
}

func (s *roleService) GetByID(id int) (*dto.RoleResponse, error) {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}
	return toRoleResponse(r), nil
}

func (s *roleService) Create(req *dto.CreateRequest) (*dto.RoleResponse, error) {
	existing, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Nama role sudah digunakan"}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(int(newID))
	if err != nil || created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data role baru"}
	}
	return toRoleResponse(created), nil
}

func (s *roleService) Update(req *dto.UpdateRequest) error {
	r, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if r == nil {
		return &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}
	if r.IsSystem {
		return &errors.BadRequestError{Message: "Role sistem tidak dapat diubah"}
	}
	return s.repo.Update(req.ID, req)
}

func (s *roleService) Delete(id int) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if r == nil {
		return &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}
	if r.IsSystem {
		return &errors.BadRequestError{Message: "Role sistem tidak dapat dihapus"}
	}
	return s.repo.Delete(id)
}

func (s *roleService) ToggleStatus(id int) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if r == nil {
		return &errors.NotFoundError{Message: "Role tidak ditemukan"}
	}
	if r.IsSystem {
		return &errors.BadRequestError{Message: "Status role sistem tidak dapat diubah"}
	}
	return s.repo.ToggleStatus(id)
}

func toRoleResponse(r *model.Role) *dto.RoleResponse {
	return &dto.RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
	}
}
