package repo

import (
	model "pos_api/domain/auth/model"

	"gorm.io/gorm"
)

type (
	AuthRepoInterface interface {
		GetUserByUsername(username string) (*model.User, error)
		GetUserByID(id int) (*model.User, error)
		CreateSession(session *model.Session) error
		GetSessionByToken(token string) (*model.Session, error)
		GetSessionByRefreshToken(token string) (*model.Session, error)
		DeleteSessionByUserID(userID int) error
		DeleteSessionByToken(token string) error
	}

	authRepo struct {
		db *gorm.DB
	}
)

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{db: db}
}
