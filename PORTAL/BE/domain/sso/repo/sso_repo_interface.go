package repo

import (
	model "portal_api/domain/sso/model"
	"time"

	"gorm.io/gorm"
)

type (
	SSORepoInterface interface {
		// GetLaunchApp mengembalikan aplikasi aktif yang boleh dibuka user (nil bila tidak ada).
		GetLaunchApp(userID int, appCode string) (*model.LaunchApp, error)
		CreateCode(codeHash string, userID, appID int, expiresAt, now time.Time) error
		GetCodeByHash(codeHash string) (*model.Code, error)
		// MarkUsed menandai kode terpakai. Mengembalikan false bila sudah terpakai lebih dulu.
		MarkUsed(id int, now time.Time) (bool, error)
		HasAccess(userID, appID int) (bool, error)
		DeleteExpired(before time.Time) error
	}

	ssoRepo struct {
		db *gorm.DB
	}
)

func NewSSORepo(db *gorm.DB) *ssoRepo {
	return &ssoRepo{db: db}
}
