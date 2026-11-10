package service

import (
	"pos_api/domain/version/dto"
	repo "pos_api/domain/version/repo"
)

type (
	VersionServiceInterface interface {
		GetAll() ([]dto.AppVersionListItem, error)
		CheckAndroid(currentVersion string) (*dto.VersionCheckResponse, error)
		UpdateAndroidVersion(req *dto.UpdateVersionRequest) error
	}

	versionService struct {
		repo repo.VersionRepoInterface
	}
)

func NewVersionService(r repo.VersionRepoInterface) *versionService {
	return &versionService{repo: r}
}
