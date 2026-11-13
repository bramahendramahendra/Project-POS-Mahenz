package service

import (
	"context"

	dto "portal_api/domain/auth/dto"
	model "portal_api/domain/auth/model"
	repo "portal_api/domain/auth/repo"
	"portal_api/pkg/authclient"
)

type (
	AuthServiceInterface interface {
		Login(ctx context.Context, req *dto.LoginRequest, ip string) (*dto.LoginResponse, error)
		Logout(token string) error
		RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
		GetMe(ctx context.Context, authUserID int) (*authclient.Response, error)
		ChangeMyPassword(ctx context.Context, userID, authUserID int, body []byte) (*authclient.Response, error)
		ValidateToken(token string) (*model.Session, error)
		VerifyToken(token string) (*dto.VerifyTokenResponse, error)
	}

	authService struct {
		repo   repo.AuthRepoInterface
		client *authclient.Client
	}
)

func NewAuthService(repo repo.AuthRepoInterface, client *authclient.Client) *authService {
	return &authService{repo: repo, client: client}
}
