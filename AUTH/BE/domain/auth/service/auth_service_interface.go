package service

import (
	dto "auth_api/domain/auth/dto"
	model "auth_api/domain/auth/model"
	repo "auth_api/domain/auth/repo"
)

type (
	AuthServiceInterface interface {
		Login(req *dto.LoginRequest, ip string) (*dto.LoginResponse, error)
		Logout(token string) error
		RefreshToken(refreshToken string) (*dto.RefreshResponse, error)
		GetMe(userID int) (*dto.UserData, error)
		ChangeMyPassword(userID int, req *dto.ChangeMyPasswordRequest) error
		ValidateToken(token string) (*model.Session, error)
		VerifyToken(token string) (*dto.VerifyTokenResponse, error)
	}

	authService struct {
		repo repo.AuthRepoInterface
	}
)

func NewAuthService(repo repo.AuthRepoInterface) *authService {
	return &authService{repo: repo}
}
