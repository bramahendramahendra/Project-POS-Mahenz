package service

import dto_supplier_return "pos_api/domain/supplier_return/dto"

type SupplierReturnService interface {
	GetAll(req *dto_supplier_return.SupplierReturnListRequest) ([]*dto_supplier_return.SupplierReturnResponse, int, error)
	GetByID(id int) (*dto_supplier_return.SupplierReturnResponse, error)
	Create(req *dto_supplier_return.CreateSupplierReturnRequest) (*dto_supplier_return.SupplierReturnResponse, error)
	UpdateStatus(req *dto_supplier_return.UpdateStatusRequest) error
	Delete(id int) error
}
