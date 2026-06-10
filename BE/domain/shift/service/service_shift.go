package service

import (
	"strings"

	dto "pos_api/domain/shift/dto"
	"pos_api/errors"
)

func (s *shiftService) GetAll(req *dto.GetAllRequest) (data []dto.ShiftResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.ShiftResponse{
			ID:        v.ID,
			Name:      v.Name,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IsActive:  v.IsActive,
		})
	}

	return data, total, nil
}

func (s *shiftService) GetActive() (data []dto.ShiftActiveResponse, err error) {
	dataDB, err := s.repo.GetActive()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.ShiftActiveResponse{
			ID:        v.ID,
			Name:      v.Name,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
	}

	return data, nil
}

func (s *shiftService) GetByID(id int) (data dto.ShiftResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Shift tidak ditemukan"}
	}

	data = dto.ShiftResponse{
		ID:        dataDB.ID,
		Name:      dataDB.Name,
		StartTime: dataDB.StartTime,
		EndTime:   dataDB.EndTime,
		IsActive:  dataDB.IsActive,
	}

	return data, nil
}

func (s *shiftService) Create(req *dto.CreateRequest) (data dto.ShiftResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)

	newID, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data shift"}
	}

	data = dto.ShiftResponse{
		ID:        dataDB.ID,
		Name:      dataDB.Name,
		StartTime: dataDB.StartTime,
		EndTime:   dataDB.EndTime,
		IsActive:  dataDB.IsActive,
	}

	return data, nil
}

func (s *shiftService) Update(req *dto.UpdateRequest) (err error) {
	req.Name = strings.TrimSpace(req.Name)

	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Shift tidak ditemukan"}
	}

	return s.repo.Update(req)
}

func (s *shiftService) Delete(req *dto.DeleteRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Shift tidak ditemukan"}
	}

	count, err := s.repo.CountOpenCashDrawer(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa penggunaan shift"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Shift tidak bisa dihapus karena sedang digunakan"}
	}

	return s.repo.Delete(req)
}

func (s *shiftService) ToggleStatus(req *dto.ToggleStatusRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Shift tidak ditemukan"}
	}

	return s.repo.ToggleStatus(req)
}

func (s *shiftService) GetSummary() (data []dto.ShiftSummaryResponse, err error) {
	dataDB, err := s.repo.GetSummary()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, *v)
	}

	return data, nil
}
