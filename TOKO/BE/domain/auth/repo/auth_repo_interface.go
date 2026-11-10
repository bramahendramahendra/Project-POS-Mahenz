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

		// SSO Portal
		GetPortalCredentials() (*model.PortalCredentials, error)
		GetUserByPortalID(portalUserID int) (*model.User, error)
		// UsernameTaken: true bila username sudah dipakai akun lokal (tertaut maupun tidak).
		UsernameTaken(username string) (bool, error)
		GetRoleIDByName(name string) (int, error)
		CreatePortalUser(username, hashedPassword, fullName string, roleID, portalUserID int) error
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
