package repo

import (
	model "auth_api/domain/auth/model"

	"gorm.io/gorm"
)

type (
	AuthRepoInterface interface {
		GetUserByUsername(username string) (*model.User, error)
		GetUserByID(id int) (*model.User, error)
		GetPasswordHashByID(id int) (string, error)
		UpdateOwnPassword(id int, hashedPassword string) error
		HasAccess(userID, appID int) (bool, error)
	}

	authRepo struct {
		db *gorm.DB
	}
)

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{db: db}
}
