package service

import (
	"strings"

	dto "coffee_api/domain/menu_size/dto"
	model "coffee_api/domain/menu_size/model"
	"coffee_api/errors"
)

func (s *menuSizeService) GetAll(req *dto.GetAllRequest) ([]dto.MenuSizeResponse, error) {
	rows, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}
	data := make([]dto.MenuSizeResponse, 0, len(rows))
	for _, v := range rows {
		data = append(data, toResponse(v))
	}
	return data, nil
}

func (s *menuSizeService) GetOptions() ([]dto.OptionResponse, error) {
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

func (s *menuSizeService) Create(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if err := s.ensureUniqueName(req.Name, 0); err != nil {
		return err
	}
	return s.repo.Create(req)
}

func (s *menuSizeService) Update(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if _, err := s.mustGet(req.ID); err != nil {
		return err
	}
	if err := s.ensureUniqueName(req.Name, req.ID); err != nil {
		return err
	}
	return s.repo.Update(req)
}

func (s *menuSizeService) ToggleStatus(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// Delete: hanya ukuran yang belum dipakai varian menu. Yang sudah dipakai cukup dinonaktifkan.
func (s *menuSizeService) Delete(id int) error {
	current, err := s.mustGet(id)
	if err != nil {
		return err
	}
	if current.UsageCount > 0 {
		return &errors.BadRequestError{Message: "Ukuran masih dipakai varian menu dan tidak bisa dihapus. Nonaktifkan saja"}
	}
	return s.repo.Delete(id)
}

func (s *menuSizeService) mustGet(id int) (*model.MenuSize, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, &errors.NotFoundError{Message: "Ukuran tidak ditemukan"}
	}
	return current, nil
}

func (s *menuSizeService) ensureUniqueName(name string, excludeID int) error {
	exists, err := s.repo.NameExists(name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama ukuran sudah digunakan"}
	}
	return nil
}

func toResponse(v *model.MenuSize) dto.MenuSizeResponse {
	return dto.MenuSizeResponse{ID: v.ID, Name: v.Name, SortOrder: v.SortOrder, IsActive: v.IsActive, UsageCount: v.UsageCount}
}
