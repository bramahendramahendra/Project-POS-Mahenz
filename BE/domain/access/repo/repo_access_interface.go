package repo

import (
	"pos_api/domain/access/dto"
	"pos_api/domain/access/model"

	"gorm.io/gorm"
)

type (
	AccessRepoInterface interface {
		GetByRoleID(roleID int) ([]*model.RoleMenuAccessItem, error)
		SetRoleAccess(roleID int, accesses []dto.SetAccessItem) error
	}

	accessRepo struct {
		db *gorm.DB
	}
)

func NewAccessRepo(db *gorm.DB) *accessRepo {
	return &accessRepo{db: db}
}
