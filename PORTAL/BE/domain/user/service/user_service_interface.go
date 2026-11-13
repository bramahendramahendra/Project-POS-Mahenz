package service

import (
	"context"

	repo "portal_api/domain/user/repo"
	"portal_api/pkg/authclient"
)

type (
	// UserServiceInterface: kelola user diteruskan ke BE User (aturan dan validasi ada di sana).
	// Setelah berhasil, data lokal Portal disesuaikan (sesi/profil user yang diubah).
	UserServiceInterface interface {
		Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error)
		ResetPassword(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error)
		ToggleStatus(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error)
		Delete(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error)
	}

	userService struct {
		repo   repo.UserRepoInterface
		client *authclient.Client
	}
)

func NewUserService(repo repo.UserRepoInterface, client *authclient.Client) *userService {
	return &userService{repo: repo, client: client}
}
