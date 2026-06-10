package repo

import (
	"pos_api/domain/setting/model"

	"gorm.io/gorm"
)

type (
	SettingRepoInterface interface {
		GetAll() ([]model.Setting, error)
		GetByKey(key string) (*model.Setting, error)
		Upsert(key, value string) error
		ResetAll() error
	}

	settingRepo struct {
		db *gorm.DB
	}
)

func NewSettingRepo(db *gorm.DB) *settingRepo {
	return &settingRepo{db: db}
}
