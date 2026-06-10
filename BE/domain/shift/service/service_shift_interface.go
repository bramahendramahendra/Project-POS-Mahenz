package service

import (
	dto "pos_api/domain/shift/dto"
	repo "pos_api/domain/shift/repo"
)

type (
	ShiftServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.ShiftResponse, total int64, err error)
		GetActive() (data []dto.ShiftActiveResponse, err error)
		GetByID(id int) (data dto.ShiftResponse, err error)
		Create(req *dto.CreateRequest) (data dto.ShiftResponse, err error)
		Update(req *dto.UpdateRequest) (err error)
		Delete(req *dto.DeleteRequest) (err error)
		ToggleStatus(req *dto.ToggleStatusRequest) (err error)
		GetSummary() (data []dto.ShiftSummaryResponse, err error)
	}

	shiftService struct {
		repo repo.ShiftRepoInterface
	}
)

func NewShiftService(repo repo.ShiftRepoInterface) *shiftService {
	return &shiftService{repo: repo}
}
