package repo

import (
	dto "auth_api/domain/role/dto"
	model "auth_api/domain/role/model"

	"gorm.io/gorm"
)

type (
	RoleRepoInterface interface {
		GetActiveOptions() ([]*dto.RoleOptionResponse, error)
		GetByID(id int) (*model.Role, error)
	}

	roleRepo struct {
		db *gorm.DB
	}
)

func NewRoleRepo(db *gorm.DB) *roleRepo {
	return &roleRepo{db: db}
}
