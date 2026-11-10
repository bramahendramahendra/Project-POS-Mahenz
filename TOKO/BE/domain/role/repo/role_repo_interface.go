package repo

import (
	dto "pos_api/domain/role/dto"
	model "pos_api/domain/role/model"

	"gorm.io/gorm"
)

type (
	RoleRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Role, int64, error)
		GetActiveOptions() ([]*dto.RoleOptionResponse, error)
		GetByID(id int) (*model.Role, error)
		GetByName(name string) (*model.Role, error)
		Create(req *dto.CreateRequest) (int64, error)
		Update(id int, req *dto.UpdateRequest) error
		Delete(id int) error
		ToggleStatus(id int) error
	}

	roleRepo struct {
		db *gorm.DB
	}
)

func NewRoleRepo(db *gorm.DB) *roleRepo {
	return &roleRepo{db: db}
}
