package repo

import (
	model "portal_api/domain/auth/model"

	"gorm.io/gorm"
)

type (
	AuthRepoInterface interface {
		GetProfileByAuthUserID(authUserID int) (*model.User, error)
		GetUnlinkedProfileByUsername(username string) (*model.User, error)
		LinkProfile(id, authUserID int, username, fullName string) error
		UpdateProfile(id int, username, fullName string) error
		CreateProfile(authUserID int, username, fullName string) (int, error)
		CreateSession(session *model.Session) error
		GetSessionByToken(token string) (*model.Session, error)
		GetSessionByRefreshToken(token string) (*model.Session, error)
		ClearSessionMustChangePassword(userID int) error
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
