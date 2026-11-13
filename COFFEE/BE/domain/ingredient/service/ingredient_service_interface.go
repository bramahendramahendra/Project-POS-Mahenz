package service

import (
	dto "coffee_api/domain/ingredient/dto"
	repo "coffee_api/domain/ingredient/repo"
)

type (
	IngredientServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]dto.IngredientResponse, int64, error)
		GetOptions() ([]dto.IngredientOptionResponse, error)
		GetByID(id int) (*dto.IngredientResponse, error)
		Create(req *dto.CreateRequest, userID int) (*dto.IngredientResponse, error)
		Update(req *dto.UpdateRequest) error
		ToggleStatus(id int) error
		Delete(id int) error

		RecordUsage(req *dto.RecordUsageRequest, userID int) (*dto.RecordResponse, error)
		RecordWaste(req *dto.RecordWasteRequest, userID int) (*dto.RecordResponse, error)
		GetMutations(req *dto.GetMutationsRequest) ([]dto.MutationResponse, int64, error)
		CancelMutation(req *dto.CancelMutationRequest, userID int) (*dto.MutationResponse, error)
	}

	ingredientService struct {
		repo repo.IngredientRepoInterface
	}
)

func NewIngredientService(repo repo.IngredientRepoInterface) *ingredientService {
	return &ingredientService{repo: repo}
}
