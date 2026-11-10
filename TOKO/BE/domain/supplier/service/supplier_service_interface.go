package service

import (
	dto "pos_api/domain/supplier/dto"
	repo "pos_api/domain/supplier/repo"
)

type (
	SupplierServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.SupplierResponse, total int64, err error)
		GetOptions(search string) (data []dto.GetOptionResponse, err error)
		GetDetail(id int) (data dto.GetDetailResponse, err error)
		Create(req *dto.CreateRequest) (data dto.SupplierResponse, err error)
		Update(req *dto.UpdateRequest) (data dto.SupplierResponse, err error)
		Delete(req *dto.DeleteRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusRequest) (err error)
	}

	supplierService struct {
		repo repo.SupplierRepoInterface
	}
)

func NewSupplierService(repo repo.SupplierRepoInterface) *supplierService {
	return &supplierService{repo: repo}
}
