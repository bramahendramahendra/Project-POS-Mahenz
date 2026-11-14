package service

import (
	dto "coffee_api/domain/modifier_group/dto"
	repo "coffee_api/domain/modifier_group/repo"
)

type (
	ModifierGroupServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.GroupResponse, error)
		// GetOptions: kelompok aktif beserta pilihan aktif (untuk form menu dan kasir).
		GetOptions() ([]dto.GroupResponse, error)
		GetByID(id int) (*dto.GroupResponse, error)
		Create(req *dto.SaveRequest) (*dto.GroupResponse, error)
		Update(req *dto.SaveRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
	}

	modifierGroupService struct {
		repo repo.ModifierGroupRepoInterface
	}
)

func NewModifierGroupService(repo repo.ModifierGroupRepoInterface) *modifierGroupService {
	return &modifierGroupService{repo: repo}
}
