package service

import repo "pos_api/domain/setting/repo"

type (
	SettingServiceInterface interface {
		GetAll() (map[string]string, error)
		GetByKey(key string) (string, error)
		Save(data map[string]string) error
		Reset() error
	}

	settingService struct {
		repo repo.SettingRepoInterface
	}
)

func NewSettingService(r repo.SettingRepoInterface) *settingService {
	return &settingService{repo: r}
}
