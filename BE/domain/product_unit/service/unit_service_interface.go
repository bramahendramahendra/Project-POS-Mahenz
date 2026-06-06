package service

import (
	dto "pos_api/domain/product_unit/dto"
	repo "pos_api/domain/product_unit/repo"
)

type (
	UnitServiceInterface interface {
		GetAll(req *dto.UnitListRequest) (data []dto.UnitResponse, total int64, err error)
		GetOptions() (data []dto.UnitOptionResponse, err error)
		GetByID(id int) (data dto.UnitResponse, err error)
		Create(req *dto.CreateUnitRequest) (data dto.UnitResponse, err error)
		Update(req *dto.UpdateUnitRequest) (data dto.UnitResponse, err error)
		Delete(req *dto.DeleteUnitRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusUnitRequest) (err error)
	}

	unitService struct {
		repo repo.UnitRepoInterface
	}
)

func NewUnitService(repo repo.UnitRepoInterface) *unitService {
	return &unitService{repo: repo}
}
