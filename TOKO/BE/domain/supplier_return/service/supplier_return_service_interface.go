package service

import (
	dto "pos_api/domain/supplier_return/dto"
	repo "pos_api/domain/supplier_return/repo"
)

type (
	SupplierReturnServiceInterface interface {
		GetAll(req *dto.SupplierReturnListRequest) (data []dto.SupplierReturnResponse, total int64, err error)
		GetByID(id int) (data dto.SupplierReturnResponse, err error)
		Create(req *dto.CreateSupplierReturnRequest) (data dto.SupplierReturnResponse, err error)
		UpdateStatus(req *dto.UpdateStatusRequest) error
		Delete(req *dto.GetSupplierReturnByIDRequest) error
	}

	supplierReturnService struct {
		repo repo.SupplierReturnRepoInterface
	}
)

func NewSupplierReturnService(repo repo.SupplierReturnRepoInterface) *supplierReturnService {
	return &supplierReturnService{repo: repo}
}
