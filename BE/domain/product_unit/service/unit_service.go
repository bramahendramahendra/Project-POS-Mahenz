package service

import (
	"strings"

	dto "pos_api/domain/product_unit/dto"
	"pos_api/errors"
)

func (s *unitService) GetAll() (data []dto.UnitResponse, err error) {
	dataDB, err := s.repo.GetAll()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.UnitResponse{
			ID:           v.ID,
			Name:         v.Name,
			Abbreviation: v.Abbreviation,
			IsActive:     v.IsActive,
			CreatedAt:    v.CreatedAt,
		})
	}

	return data, nil
}

func (s *unitService) GetActive() (data []dto.UnitActiveResponse, err error) {
	dataDB, err := s.repo.GetActive()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.UnitActiveResponse{
			ID:           v.ID,
			Name:         v.Name,
			Abbreviation: v.Abbreviation,
		})
	}

	return data, nil
}

func (s *unitService) GetByID(id int) (data dto.UnitResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}

	data = dto.UnitResponse{
		ID:           dataDB.ID,
		Name:         dataDB.Name,
		Abbreviation: dataDB.Abbreviation,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *unitService) Create(req *dto.CreateUnitRequest) (data dto.UnitResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Abbreviation = strings.TrimSpace(req.Abbreviation)

	exists, err := s.repo.CheckNameExists(req.Name, 0)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama satuan sudah digunakan"}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}

	data = dto.UnitResponse{
		ID:           dataDB.ID,
		Name:         dataDB.Name,
		Abbreviation: dataDB.Abbreviation,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *unitService) Update(req *dto.UpdateUnitRequest) (data dto.UnitResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Abbreviation = strings.TrimSpace(req.Abbreviation)

	existsUpdate, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if existsUpdate == nil {
		return data, &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
	}

	exists, err := s.repo.CheckNameExists(req.Name, req.ID)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama satuan sudah digunakan"}
	}

	err = s.repo.Update(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}

	data = dto.UnitResponse{
		ID:           dataDB.ID,
		Name:         dataDB.Name,
		Abbreviation: dataDB.Abbreviation,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *unitService) Delete(req *dto.DeleteUnitRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
	}

	count, err := s.repo.CountProductUnitsByUnit(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Satuan masih digunakan oleh produk"}
	}

	return s.repo.Delete(req)
}

func (s *unitService) ToggleStatus(req *dto.ToggleStatusUnitRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
	}

	return s.repo.ToggleStatus(req)
}
