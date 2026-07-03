package service

import (
	dto "pos_api/domain/user/dto"
	repo "pos_api/domain/user/repo"
)

type (
	UserServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.UserResponse, int64, error)
		GetByID(id int) (*dto.UserResponse, error)
		Create(req *dto.CreateRequest) (*dto.UserResponse, error)
		Update(id int, req *dto.UpdateRequest) error
		ChangePassword(id int, req *dto.ChangePasswordRequest) error
		Delete(id, currentUserID int) error
		ToggleStatus(id int) error
	}

	userService struct {
		repo repo.UserRepoInterface
	}
)

func NewUserService(repo repo.UserRepoInterface) *userService {
	return &userService{repo: repo}
}
