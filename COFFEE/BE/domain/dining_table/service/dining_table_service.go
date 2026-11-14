package service

import (
	"strings"

	dto "coffee_api/domain/dining_table/dto"
	model "coffee_api/domain/dining_table/model"
	"coffee_api/errors"
)

func (s *diningTableService) GetAll(req *dto.GetAllRequest) ([]dto.DiningTableResponse, error) {
	rows, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}
	data := make([]dto.DiningTableResponse, 0, len(rows))
	for _, v := range rows {
		data = append(data, toResponse(v))
	}
	return data, nil
}

func (s *diningTableService) GetOptions() ([]dto.OptionResponse, error) {
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

func (s *diningTableService) Create(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if err := s.ensureUniqueName(req.Name, 0); err != nil {
		return err
	}
	return s.repo.Create(req)
}

func (s *diningTableService) Update(req *dto.SaveRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if _, err := s.mustGet(req.ID); err != nil {
		return err
	}
	if err := s.ensureUniqueName(req.Name, req.ID); err != nil {
		return err
	}
	return s.repo.Update(req)
}

func (s *diningTableService) ToggleStatus(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// Delete: hanya meja yang belum pernah dipakai pesanan (mulai 4.7). Yang sudah dipakai cukup dinonaktifkan.
func (s *diningTableService) Delete(id int) error {
	current, err := s.mustGet(id)
	if err != nil {
		return err
	}
	if current.UsageCount > 0 {
		return &errors.BadRequestError{Message: "Meja sudah pernah dipakai pesanan dan tidak bisa dihapus. Nonaktifkan saja"}
	}
	return s.repo.Delete(id)
}

func (s *diningTableService) mustGet(id int) (*model.DiningTable, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, &errors.NotFoundError{Message: "Meja tidak ditemukan"}
	}
	return current, nil
}

func (s *diningTableService) ensureUniqueName(name string, excludeID int) error {
	exists, err := s.repo.NameExists(name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama meja sudah digunakan"}
	}
	return nil
}

func toResponse(v *model.DiningTable) dto.DiningTableResponse {
	return dto.DiningTableResponse{ID: v.ID, Name: v.Name, SortOrder: v.SortOrder, IsActive: v.IsActive, UsageCount: v.UsageCount}
}
