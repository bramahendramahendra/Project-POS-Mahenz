package service

import (
	dto "pos_api/domain/menu/dto"
	model "pos_api/domain/menu/model"
	"pos_api/errors"
)

func (s *menuService) GetAll(req *dto.GetAllRequest) ([]*dto.MenuResponse, error) {
	menus, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.MenuResponse, 0, len(menus))
	for _, m := range menus {
		result = append(result, toMenuResponse(m))
	}
	return result, nil
}

func (s *menuService) GetByID(id int) (*dto.MenuResponse, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}
	return toMenuResponse(m), nil
}

func (s *menuService) GetMyMenus(roleName string) ([]dto.MyMenuItem, error) {
	flat, err := s.repo.GetMyMenus(roleName)
	if err != nil {
		return nil, err
	}
	return buildTree(flat), nil
}

func (s *menuService) Create(req *dto.CreateRequest) (*dto.MenuResponse, error) {
	existing, err := s.repo.GetByKeyName(req.KeyName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Key menu sudah digunakan"}
	}

	if req.ParentID != nil {
		parent, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, &errors.BadRequestError{Message: "Parent menu tidak ditemukan"}
		}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(int(newID))
	if err != nil || created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data menu baru"}
	}
	return toMenuResponse(created), nil
}

func (s *menuService) Update(id int, req *dto.UpdateRequest) error {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}

	if req.ParentID != nil && *req.ParentID == id {
		return &errors.BadRequestError{Message: "Menu tidak bisa menjadi parent dirinya sendiri"}
	}

	if req.ParentID != nil {
		parent, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return &errors.BadRequestError{Message: "Parent menu tidak ditemukan"}
		}
	}

	return s.repo.Update(id, req)
}

func (s *menuService) Delete(id int) error {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}
	return s.repo.Delete(id)
}

func (s *menuService) Reorder(req *dto.ReorderRequest) error {
	return s.repo.Reorder(req.Items)
}

func buildTree(flat []*dto.MyMenuItem) []dto.MyMenuItem {
	var roots []dto.MyMenuItem
	childMap := make(map[string][]dto.MyMenuItem)

	for _, item := range flat {
		dotIdx := lastDot(item.KeyName)
		if dotIdx == -1 {
			roots = append(roots, *item)
		} else {
			parentKey := item.KeyName[:dotIdx]
			childMap[parentKey] = append(childMap[parentKey], *item)
		}
	}

	var attachChildren func(items []dto.MyMenuItem) []dto.MyMenuItem
	attachChildren = func(items []dto.MyMenuItem) []dto.MyMenuItem {
		result := make([]dto.MyMenuItem, 0, len(items))
		for _, item := range items {
			if children, ok := childMap[item.KeyName]; ok {
				item.Children = attachChildren(children)
			}
			result = append(result, item)
		}
		return result
	}

	return attachChildren(roots)
}

func lastDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}

func toMenuResponse(m *model.Menu) *dto.MenuResponse {
	return &dto.MenuResponse{
		ID:         m.ID,
		ParentID:   m.ParentID,
		KeyName:    m.KeyName,
		Label:      m.Label,
		Icon:       m.Icon,
		Path:       m.Path,
		OrderIndex: m.OrderIndex,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
	}
}
