package service

import (
	role_repo "pos_api/domain/role/repo"
	dto "pos_api/domain/user/dto"
	repo "pos_api/domain/user/repo"
)

type (
	UserServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.UserResponse, int64, error)
		GetByID(id int) (*dto.UserResponse, error)
		Create(req *dto.CreateRequest) (*dto.UserResponse, error)
		Update(id, currentUserID int, req *dto.UpdateRequest) error
		ChangePassword(id int, req *dto.ChangePasswordRequest) error
		Delete(id, currentUserID int) error
		ToggleStatus(id, currentUserID int) error
	}

	userService struct {
		repo     repo.UserRepoInterface
		roleRepo role_repo.RoleRepoInterface
	}
)

func NewUserService(repo repo.UserRepoInterface, roleRepo role_repo.RoleRepoInterface) *userService {
	return &userService{repo: repo, roleRepo: roleRepo}
}
