package service

import (
	role_repo "portal_api/domain/role/repo"
	dto "portal_api/domain/user/dto"
	repo "portal_api/domain/user/repo"
)

type (
	UserServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.UserResponse, int64, error)
		GetByID(id int) (*dto.UserResponse, error)
		Create(req *dto.CreateRequest) (*dto.UserResponse, error)
		Update(id, currentUserID int, req *dto.UpdateRequest) error
		ResetPassword(id, currentUserID int, req *dto.ResetPasswordRequest) error
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
