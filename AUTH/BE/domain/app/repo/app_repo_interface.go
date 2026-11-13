package repo

import (
	model "auth_api/domain/app/model"

	"gorm.io/gorm"
)

type (
	AppRepoInterface interface {
		GetMyApps(userID int) ([]*model.App, error)

		// Khusus admin
		GetAll() ([]*model.AppAdmin, error)
		GetByID(id int) (*model.AppAdmin, error)
		GetByCode(code string) (*model.AppAdmin, error)
		Create(app *model.App, keyHash string) (int64, error)
		Update(id int, app *model.App) error
		ToggleStatus(id int) error
		SetKeyHash(id int, keyHash string) error
		GetUserAccess(userID int) ([]*model.UserAccess, error)
		// SetUserAccess mengganti seluruh akses user dengan daftar appIDs (transaksional).
		SetUserAccess(userID int, appIDs []int) error
		UserExists(userID int) (bool, error)
	}

	appRepo struct {
		db *gorm.DB
	}
)

func NewAppRepo(db *gorm.DB) *appRepo {
	return &appRepo{db: db}
}
