package service

import (
	"fmt"
	"strings"

	dto "pos_api/domain/customer/dto"
	"pos_api/errors"
)

func (s *customerService) GetAll(req *dto.GetAllRequest) (data []dto.CustomerResponse, total int64, err error) {
	data = []dto.CustomerResponse{}
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.CustomerResponse{
			ID:           v.ID,
			CustomerCode: v.CustomerCode,
			Name:         v.Name,
			Phone:        v.Phone,
			Address:      v.Address,
			CreditLimit:  v.CreditLimit,
			IsActive:     v.IsActive,
			CreatedAt:    v.CreatedAt,
		})
	}

	return data, total, nil
}

func (s *customerService) GetOptions() (data []dto.CustomerActiveItem, err error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.CustomerActiveItem{
			ID:           v.ID,
			Name:         v.Name,
			CustomerCode: v.CustomerCode,
			CreditLimit:  v.CreditLimit,
		})
	}

	return data, nil
}

func (s *customerService) GetByID(id int) (data dto.CustomerDetailResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	data = dto.CustomerDetailResponse{
		ID:           dataDB.ID,
		CustomerCode: dataDB.CustomerCode,
		Name:         dataDB.Name,
		Phone:        dataDB.Phone,
		Address:      dataDB.Address,
		CreditLimit:  dataDB.CreditLimit,
		Notes:        dataDB.Notes,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *customerService) Create(req *dto.CreateRequest) (data dto.CustomerResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)

	count, err := s.repo.GetCount()
	if err != nil {
		return data, err
	}
	code := fmt.Sprintf("CUS-%03d", count+1)

	newID, err := s.repo.Create(req, code)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data pelanggan"}
	}

	data = dto.CustomerResponse{
		ID:           dataDB.ID,
		CustomerCode: dataDB.CustomerCode,
		Name:         dataDB.Name,
		Phone:        dataDB.Phone,
		Address:      dataDB.Address,
		CreditLimit:  dataDB.CreditLimit,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *customerService) Update(req *dto.UpdateRequest) (data dto.CustomerResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)

	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if exists == nil {
		return data, &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	if err = s.repo.Update(req); err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data pelanggan"}
	}

	data = dto.CustomerResponse{
		ID:           dataDB.ID,
		CustomerCode: dataDB.CustomerCode,
		Name:         dataDB.Name,
		Phone:        dataDB.Phone,
		Address:      dataDB.Address,
		CreditLimit:  dataDB.CreditLimit,
		IsActive:     dataDB.IsActive,
		CreatedAt:    dataDB.CreatedAt,
	}

	return data, nil
}

func (s *customerService) Delete(req *dto.DeleteRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	count, err := s.repo.CountActiveReceivables(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa piutang pelanggan"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Pelanggan masih memiliki piutang aktif"}
	}

	return s.repo.Delete(req)
}

func (s *customerService) ToggleStatus(req *dto.ToggleStatusRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	return s.repo.ToggleStatus(req)
}
