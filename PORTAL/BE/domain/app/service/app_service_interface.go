package service

import (
	"context"

	"portal_api/pkg/authclient"
)

type (
	// AppServiceInterface: permintaan diteruskan ke BE User atas nama user yang sedang login.
	AppServiceInterface interface {
		Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error)
	}

	appService struct {
		client *authclient.Client
	}
)

func NewAppService(client *authclient.Client) *appService {
	return &appService{client: client}
}
