package service

import (
	"strings"

	dto "pos_api/domain/product_unit/dto"
	"pos_api/errors"
)

func (s *unitService) GetAll(req *dto.GetAllRequest) (data []dto.UnitResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
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

	return data, total, nil
}

func (s *unitService) GetOptions() (data []dto.GetOptionResponse, err error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.GetOptionResponse{
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
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
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

func (s *unitService) Create(req *dto.CreateRequest) (data dto.UnitResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Abbreviation = strings.TrimSpace(req.Abbreviation)

	exists, err := s.repo.CheckNameExists(req.Name, 0)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama satuan sudah digunakan"}
	}

	abbrExists, err := s.repo.CheckAbbreviationExists(req.Abbreviation, 0)
	if err != nil {
		return data, err
	}
	if abbrExists {
		return data, &errors.BadRequestError{Message: "Singkatan satuan sudah digunakan"}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data satuan"}
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

func (s *unitService) Update(req *dto.UpdateRequest) (data dto.UnitResponse, err error) {
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

	abbrExists, err := s.repo.CheckAbbreviationExists(req.Abbreviation, req.ID)
	if err != nil {
		return data, err
	}
	if abbrExists {
		return data, &errors.BadRequestError{Message: "Singkatan satuan sudah digunakan"}
	}

	if err = s.repo.Update(req); err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data satuan"}
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

func (s *unitService) Delete(req *dto.DeleteRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
	}

	count, err := s.repo.CountProductUnitsByUnit(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa penggunaan satuan"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Satuan masih digunakan oleh produk"}
	}

	return s.repo.Delete(req)
}

func (s *unitService) ToggleStatus(req *dto.ToggleStatusRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Satuan tidak ditemukan"}
	}

	if exists.IsActive {
		activeCount, err := s.repo.CountActiveProductsByUnit(req.ID)
		if err != nil {
			return &errors.InternalServerError{Message: "Gagal memeriksa produk aktif satuan"}
		}
		if activeCount > 0 {
			return &errors.BadRequestError{Message: "Satuan tidak bisa dinonaktifkan karena masih digunakan oleh produk aktif"}
		}
	}

	return s.repo.ToggleStatus(req)
}
