package repo

import "gorm.io/gorm"

type (
	// UserRepoInterface: data lokal Portal yang ikut berubah saat akun dikelola di BE User.
	UserRepoInterface interface {
		DeleteSessionByAuthUserID(authUserID int) error
		DeleteProfileByAuthUserID(authUserID int) error
	}

	userRepo struct {
		db *gorm.DB
	}
)

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}
