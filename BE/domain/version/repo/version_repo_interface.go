package repo

import (
	"pos_api/domain/version/model"

	"gorm.io/gorm"
)

type (
	VersionRepoInterface interface {
		GetAll() ([]model.AppVersion, error)
		GetLatestAndroid() (*model.AppVersion, error)
		SetAllNotLatest() error
		CreateVersion(version, downloadURL, releaseNotes string, isMandatory bool) error
	}

	versionRepo struct {
		db *gorm.DB
	}
)

func NewVersionRepo(db *gorm.DB) *versionRepo {
	return &versionRepo{db: db}
}
