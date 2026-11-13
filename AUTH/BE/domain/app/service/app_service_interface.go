package service

import (
	dto "auth_api/domain/app/dto"
	repo "auth_api/domain/app/repo"
)

type (
	AppServiceInterface interface {
		GetMyApps(userID int) ([]*dto.MyAppResponse, error)

		// Khusus admin
		GetAll() ([]*dto.AppAdminResponse, error)
		Create(req *dto.CreateRequest) (*dto.CreateAppResponse, error)
		Update(id int, req *dto.UpdateRequest) error
		ToggleStatus(id int) error
		RotateKey(id int) (*dto.RotateKeyResponse, error)
		GetUserAccess(userID int) ([]*dto.UserAccessItem, error)
		SetUserAccess(userID int, req *dto.SetUserAccessRequest) error
	}

	appService struct {
		repo repo.AppRepoInterface
	}
)

func NewAppService(repo repo.AppRepoInterface) *appService {
	return &appService{repo: repo}
}
