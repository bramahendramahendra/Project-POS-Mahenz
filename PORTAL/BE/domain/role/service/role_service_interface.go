package service

import (
	"context"

	"portal_api/pkg/authclient"
)

type (
	// RoleServiceInterface: permintaan diteruskan ke BE User atas nama user yang sedang login.
	RoleServiceInterface interface {
		Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error)
	}

	roleService struct {
		client *authclient.Client
	}
)

func NewRoleService(client *authclient.Client) *roleService {
	return &roleService{client: client}
}
