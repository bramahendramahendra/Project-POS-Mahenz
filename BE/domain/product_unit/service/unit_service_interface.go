package service

import (
	dto "pos_api/domain/product_unit/dto"
	repo "pos_api/domain/product_unit/repo"
)

type (
	UnitServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.UnitResponse, total int64, err error)
		GetOptions() (data []dto.GetOptionResponse, err error)
		GetByID(id int) (data dto.UnitResponse, err error)
		Create(req *dto.CreateRequest) (data dto.UnitResponse, err error)
		Update(req *dto.UpdateRequest) (data dto.UnitResponse, err error)
		Delete(req *dto.DeleteRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusRequest) (err error)
	}

	unitService struct {
		repo repo.UnitRepoInterface
	}
)

func NewUnitService(repo repo.UnitRepoInterface) *unitService {
	return &unitService{repo: repo}
}
