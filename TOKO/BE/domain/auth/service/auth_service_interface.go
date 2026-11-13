package service

import (
	"context"

	dto "pos_api/domain/auth/dto"
	model "pos_api/domain/auth/model"
	repo "pos_api/domain/auth/repo"
	"pos_api/pkg/authclient"
)

type (
	AuthServiceInterface interface {
		// Login: form username + password, diverifikasi BE User.
		Login(ctx context.Context, req *dto.LoginRequest, ip string) (*dto.LoginResponse, error)
		Logout(token string) error
		RefreshToken(refreshToken string) (*dto.RefreshResponse, error)
		GetMe(userID int) (*dto.UserData, error)
		ValidateToken(token string) (*model.Session, error)
		VerifyToken(token string) (*dto.VerifyTokenResponse, error)
		// LoginSSO menukar kode sekali pakai dari Portal (dibuat BE User) menjadi sesi Toko.
		LoginSSO(ctx context.Context, req *dto.SSORequest, ip string) (*dto.LoginResponse, error)
		// GetLoginOptions: cara login yang tersedia (publik).
		GetLoginOptions() (*dto.LoginOptionsResponse, error)
		// SetDirectLogin menyalakan/mematikan login langsung (hanya admin; dicek di route).
		SetDirectLogin(enabled bool) error
	}

	authService struct {
		repo   repo.AuthRepoInterface
		client *authclient.Client
	}
)

func NewAuthService(repo repo.AuthRepoInterface, client *authclient.Client) *authService {
	return &authService{repo: repo, client: client}
}
