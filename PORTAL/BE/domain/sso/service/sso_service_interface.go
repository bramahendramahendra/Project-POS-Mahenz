package service

import (
	"context"

	"portal_api/pkg/authclient"
)

type (
	// SSOServiceInterface: permintaan diteruskan ke BE User atas nama user yang sedang login.
	SSOServiceInterface interface {
		Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error)
	}

	ssoService struct {
		client *authclient.Client
	}
)

func NewSSOService(client *authclient.Client) *ssoService {
	return &ssoService{client: client}
}
