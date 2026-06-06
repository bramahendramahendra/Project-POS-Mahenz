package service

import (
	dto "pos_api/domain/supplier/dto"
	repo "pos_api/domain/supplier/repo"
)

type (
	SupplierServiceInterface interface {
		GetAll(req *dto.SupplierListRequest) (data []dto.SupplierResponse, total int64, err error)
		GetOptions() (data []dto.SupplierOptionResponse, err error)
		GetDetail(id int) (data dto.SupplierDetailResponse, err error)
		Create(req *dto.CreateSupplierRequest) (data dto.SupplierResponse, err error)
		Update(req *dto.UpdateSupplierRequest) (data dto.SupplierResponse, err error)
		Delete(req *dto.DeleteSupplierRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusSupplierRequest) (err error)
	}

	supplierService struct {
		repo repo.SupplierRepo
	}
)

func NewSupplierService(repo repo.SupplierRepo) *supplierService {
	return &supplierService{repo: repo}
}
