package service

import (
	"pos_api/domain/setting/dto"
	repo "pos_api/domain/setting/repo"
)

type (
	SettingServiceInterface interface {
		GetAll() (map[string]string, error)
		GetByKey(key string) (string, error)
		Save(data map[string]string) error
		Reset() error
		GetStoreProfile() (*dto.StoreProfileResponse, error)
		UpdateStoreProfile(req *dto.StoreProfileRequest) error
		GetPrinterSettings() (*dto.PrinterSettingsResponse, error)
		UpdatePrinterSettings(req *dto.PrinterSettingsRequest) error
	}

	settingService struct {
		repo repo.SettingRepoInterface
	}
)

func NewSettingService(r repo.SettingRepoInterface) *settingService {
	return &settingService{repo: r}
}
