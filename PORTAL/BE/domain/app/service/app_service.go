package service

import (
	"context"

	"portal_api/pkg/authclient"
)

func (s *appService) Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error) {
	resp, err := s.client.Post(ctx, path, actingUserID, body)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}
	return resp, nil
}
