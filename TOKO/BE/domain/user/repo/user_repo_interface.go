package repo

import (
	dto "pos_api/domain/user/dto"
	model "pos_api/domain/user/model"

	"gorm.io/gorm"
)

type (
	UserRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.User, int64, error)
		GetByID(id int) (*model.User, error)
		Update(id int, req *dto.UpdateRequest) error
		Delete(id int) error
		ToggleStatus(id int) error
		DeleteSessionByUserID(userID int) error
		CountActiveAdmins(excludeID int) (int64, error)
	}

	userRepo struct {
		db *gorm.DB
	}
)

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}
