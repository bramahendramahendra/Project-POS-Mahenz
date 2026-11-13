package repo

import (
	model "pos_api/domain/auth/model"

	"gorm.io/gorm"
)

type (
	AuthRepoInterface interface {
		GetUserByID(id int) (*model.User, error)
		CreateSession(session *model.Session) error
		GetSessionByToken(token string) (*model.Session, error)
		GetSessionByRefreshToken(token string) (*model.Session, error)
		DeleteSessionByUserID(userID int) error
		DeleteSessionByToken(token string) error

		// Profil lokal yang ditautkan ke akun BE User.
		GetUserByAuthID(authUserID int) (*model.User, error)
		// GetUnlinkedUserByUsername: profil lokal dengan username sama yang belum tertaut.
		GetUnlinkedUserByUsername(username string) (*model.User, error)
		// UsernameTaken: true bila username sudah dipakai profil lokal (tertaut maupun tidak).
		UsernameTaken(username string) (bool, error)
		LinkUser(userID, authUserID int, fullName string) error
		GetRoleIDByName(name string) (int, error)
		CreateAuthUser(username, fullName string, roleID, authUserID int) error
		UpdateFullName(userID int, fullName string) error

		// Toggle login langsung (settings.direct_login_enabled). Nilai kosong = belum diatur.
		GetDirectLoginSetting() (string, error)
		SetDirectLoginSetting(value string) error
	}

	authRepo struct {
		db *gorm.DB
	}
)

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{db: db}
}
