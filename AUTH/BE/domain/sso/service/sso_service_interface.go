package service

import (
	dto "auth_api/domain/sso/dto"
	repo "auth_api/domain/sso/repo"
)

type (
	SSOServiceInterface interface {
		Launch(userID int, req *dto.LaunchRequest) (*dto.LaunchResponse, error)
		Exchange(appKey string, req *dto.ExchangeRequest) (*dto.ExchangeResponse, error)
	}

	ssoService struct {
		repo repo.SSORepoInterface
	}
)

func NewSSOService(repo repo.SSORepoInterface) *ssoService {
	return &ssoService{repo: repo}
}
