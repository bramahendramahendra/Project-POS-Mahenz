package service

import (
	"strings"

	dto "coffee_api/domain/menu_category/dto"
	model "coffee_api/domain/menu_category/model"
	"coffee_api/errors"
)

func (s *menuCategoryService) GetAll(req *dto.GetAllRequest) ([]dto.MenuCategoryResponse, error) {
	rows, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}
	data := make([]dto.MenuCategoryResponse, 0, len(rows))
	for _, v := range rows {
		data = append(data, toResponse(v))
	}
	return data, nil
}

func (s *menuCategoryService) GetOptions() ([]dto.OptionResponse, error) {
	rows, err := s.repo.GetOptions()
	if err != nil {
		return nil, err
	}
	data := make([]dto.OptionResponse, 0, len(rows))
	for _, v := range rows {
		data = append(data, dto.OptionResponse{ID: v.ID, Name: v.Name})
	}
	return data, nil
}

func (s *menuCategoryService) Create(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if err := s.ensureUniqueName(req.Name, 0); err != nil {
		return err
	}
	return s.repo.Create(req)
}

func (s *menuCategoryService) Update(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if _, err := s.mustGet(req.ID); err != nil {
		return err
	}
	if err := s.ensureUniqueName(req.Name, req.ID); err != nil {
		return err
	}
	return s.repo.Update(req)
}

func (s *menuCategoryService) ToggleStatus(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// Delete: hanya kategori tanpa menu. Kategori yang masih dipakai cukup dinonaktifkan.
func (s *menuCategoryService) Delete(id int) error {
	current, err := s.mustGet(id)
	if err != nil {
		return err
	}
	if current.ItemCount > 0 {
		return &errors.BadRequestError{Message: "Kategori masih dipakai menu dan tidak bisa dihapus. Nonaktifkan saja"}
	}
	return s.repo.Delete(id)
}

func (s *menuCategoryService) mustGet(id int) (*model.MenuCategory, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, &errors.NotFoundError{Message: "Kategori menu tidak ditemukan"}
	}
	return current, nil
}

func (s *menuCategoryService) ensureUniqueName(name string, excludeID int) error {
	exists, err := s.repo.NameExists(name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama kategori menu sudah digunakan"}
	}
	return nil
}

func toResponse(v *model.MenuCategory) dto.MenuCategoryResponse {
	return dto.MenuCategoryResponse{ID: v.ID, Name: v.Name, SortOrder: v.SortOrder, IsActive: v.IsActive, ItemCount: v.ItemCount}
}
