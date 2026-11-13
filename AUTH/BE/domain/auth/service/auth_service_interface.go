package service

import (
	dto "auth_api/domain/auth/dto"
	repo "auth_api/domain/auth/repo"
	"auth_api/helper"
)

type (
	AuthServiceInterface interface {
		Verify(req *dto.VerifyRequest, caller helper.CallerApp) (*dto.UserData, error)
		GetMe(userID int) (*dto.UserData, error)
		ChangeMyPassword(userID int, req *dto.ChangeMyPasswordRequest) error
	}

	authService struct {
		repo repo.AuthRepoInterface
	}
)

func NewAuthService(repo repo.AuthRepoInterface) *authService {
	return &authService{repo: repo}
}
