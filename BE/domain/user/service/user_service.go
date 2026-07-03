package service

import (
	dto "pos_api/domain/user/dto"
	model "pos_api/domain/user/model"
	"pos_api/errors"
	"pos_api/pkg/bcrypt"
)

func (s *userService) GetAll(req *dto.GetAllRequest) ([]*dto.UserResponse, int64, error) {
	users, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*dto.UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toUserResponse(u))
	}
	return result, total, nil
}

func (s *userService) GetByID(id int) (*dto.UserResponse, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	return toUserResponse(u), nil
}

func (s *userService) Create(req *dto.CreateRequest) (*dto.UserResponse, error) {
	existing, err := s.repo.GetByUsername(req.Username, 0)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Username sudah digunakan"}
	}

	hashed, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashed,
		FullName: req.FullName,
		RoleID:   req.RoleID,
	}

	newID, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(int(newID))
	if err != nil || created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data user baru"}
	}
	return toUserResponse(created), nil
}

func (s *userService) Update(id int, req *dto.UpdateRequest) error {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}

	return s.repo.Update(id, req)
}

func (s *userService) ChangePassword(id int, req *dto.ChangePasswordRequest) error {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}

	hashed, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(id, hashed)
}

func (s *userService) Delete(id, currentUserID int) error {
	if id == currentUserID {
		return &errors.BadRequestError{Message: "Tidak bisa menghapus akun sendiri"}
	}

	u, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}

	if err := s.repo.DeleteSessionByUserID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *userService) ToggleStatus(id int) error {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	return s.repo.ToggleStatus(id)
}

func toUserResponse(u *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		RoleID:    u.RoleID,
		RoleName:  u.RoleName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
