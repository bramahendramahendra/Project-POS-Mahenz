package service

import (
	dto "coffee_api/domain/dining_table/dto"
	repo "coffee_api/domain/dining_table/repo"
)

type (
	DiningTableServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.DiningTableResponse, error)
		GetOptions() ([]dto.OptionResponse, error)
		Create(req *dto.SaveRequest) error
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	diningTableService struct {
		repo repo.DiningTableRepoInterface
	}
)

func NewDiningTableService(repo repo.DiningTableRepoInterface) *diningTableService {
	return &diningTableService{repo: repo}
}
