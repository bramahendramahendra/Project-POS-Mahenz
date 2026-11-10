package service

import (
	"context"
	dto "pos_api/domain/auth/dto"
	model "pos_api/domain/auth/model"
	repo "pos_api/domain/auth/repo"
)

type (
	AuthServiceInterface interface {
		Login(req *dto.LoginRequest, ip string) (*dto.LoginResponse, error)
		Logout(token string) error
		RefreshToken(refreshToken string) (*dto.RefreshResponse, error)
		GetMe(userID int) (*dto.UserData, error)
		ValidateToken(token string) (*model.Session, error)
		VerifyToken(token string) (*dto.VerifyTokenResponse, error)
		// LoginSSO menukar kode sekali pakai dari Portal menjadi sesi Toko.
		LoginSSO(ctx context.Context, req *dto.SSORequest, ip string) (*dto.LoginResponse, error)
		// GetLoginOptions: cara login yang tersedia (publik).
		GetLoginOptions() (*dto.LoginOptionsResponse, error)
		// SetDirectLogin menyalakan/mematikan login langsung (hanya admin; dicek di route).
		SetDirectLogin(enabled bool) error
	}

	authService struct {
		repo repo.AuthRepoInterface
	}
)

func NewAuthService(repo repo.AuthRepoInterface) *authService {
	return &authService{repo: repo}
}
