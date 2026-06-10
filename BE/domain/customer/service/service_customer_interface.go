package service

import (
	dto "pos_api/domain/customer/dto"
	repo "pos_api/domain/customer/repo"
)

type (
	CustomerServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.CustomerResponse, total int64, err error)
		GetOptions() (data []dto.CustomerActiveItem, err error)
		GetByID(id int) (data dto.CustomerDetailResponse, err error)
		Create(req *dto.CreateRequest) (data dto.CustomerResponse, err error)
		Update(req *dto.UpdateRequest) (data dto.CustomerResponse, err error)
		Delete(req *dto.DeleteRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusRequest) (err error)
	}

	customerService struct {
		repo repo.CustomerRepoInterface
	}
)

func NewCustomerService(repo repo.CustomerRepoInterface) *customerService {
	return &customerService{repo: repo}
}
