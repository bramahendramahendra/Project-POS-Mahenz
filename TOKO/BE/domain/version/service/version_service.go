package service

import (
	"pos_api/domain/version/dto"
	"pos_api/errors"
)

func (s *versionService) GetAll() ([]dto.AppVersionListItem, error) {
	versions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	items := make([]dto.AppVersionListItem, len(versions))
	for i, v := range versions {
		items[i] = dto.AppVersionListItem{
			ID:           v.ID,
			Platform:     v.Platform,
			Version:      v.Version,
			DownloadURL:  v.DownloadURL,
			ReleaseNotes: v.ReleaseNotes,
			IsMandatory:  v.IsMandatory,
			IsLatest:     v.IsLatest,
			CreatedAt:    v.CreatedAt,
		}
	}
	return items, nil
}

func (s *versionService) CheckAndroid(currentVersion string) (*dto.VersionCheckResponse, error) {
	latest, err := s.repo.GetLatestAndroid()
	if err != nil {
		return nil, &errors.NotFoundError{Message: "Data versi tidak ditemukan"}
	}

	hasUpdate := latest.Version != currentVersion

	resp := &dto.VersionCheckResponse{
		LatestVersion:  latest.Version,
		CurrentVersion: currentVersion,
		HasUpdate:      hasUpdate,
	}

	if hasUpdate {
		resp.DownloadURL = latest.DownloadURL
		resp.ReleaseNotes = latest.ReleaseNotes
		resp.IsMandatory = latest.IsMandatory
	}

	return resp, nil
}

func (s *versionService) UpdateAndroidVersion(req *dto.UpdateVersionRequest) error {
	if err := s.repo.SetAllNotLatest(); err != nil {
		return &errors.InternalServerError{Message: "Gagal mereset versi sebelumnya"}
	}

	if err := s.repo.CreateVersion(req.Version, req.DownloadURL, req.ReleaseNotes, req.IsMandatory); err != nil {
		return &errors.InternalServerError{Message: "Gagal menyimpan versi baru"}
	}

	return nil
}
