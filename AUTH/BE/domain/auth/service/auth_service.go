package service

import (
	dto "auth_api/domain/auth/dto"
	model "auth_api/domain/auth/model"
	"auth_api/errors"
	"auth_api/helper"
	"auth_api/pkg/bcrypt"
)

// Pesan seragam untuk username tidak ada / password salah agar tidak membocorkan keberadaan akun.
const wrongCredentials = "Username atau password salah"

// Verify memeriksa username + password untuk aplikasi pemanggil (keputusan 63).
// Urutan: akun ada dan password benar -> aktif -> wajib ganti password -> akses aplikasi.
// Portal boleh menerima akun yang wajib ganti password (Portal memaksa ganti password);
// aplikasi lain menolaknya (keputusan 70). Akses aplikasi tidak diperiksa untuk Portal.
func (s *authService) Verify(req *dto.VerifyRequest, caller helper.CallerApp) (*dto.UserData, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !bcrypt.VerifyPassword(req.Password, user.Password) {
		return nil, &errors.UnauthenticatedError{Message: wrongCredentials}
	}
	if !user.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	if caller.Code != helper.PortalAppCode {
		if user.MustChangePassword {
			return nil, &errors.UnauthorizededError{Message: "Ganti password dulu di Portal"}
		}
		ok, err := s.repo.HasAccess(user.ID, caller.ID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke aplikasi ini"}
		}
	}

	return toUserData(user), nil
}

func (s *authService) GetMe(userID int) (*dto.UserData, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	return toUserData(user), nil
}

// ChangeMyPassword mengganti password milik sendiri. Password lama wajib benar
// dan password baru harus berbeda.
func (s *authService) ChangeMyPassword(userID int, req *dto.ChangeMyPasswordRequest) error {
	currentHash, err := s.repo.GetPasswordHashByID(userID)
	if err != nil {
		return err
	}
	if currentHash == "" {
		return &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	if !bcrypt.VerifyPassword(req.OldPassword, currentHash) {
		return &errors.BadRequestError{Message: "Password lama salah"}
	}
	if req.OldPassword == req.NewPassword {
		return &errors.BadRequestError{Message: "Password baru harus berbeda dari password lama"}
	}

	hashed, err := bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdateOwnPassword(userID, hashed)
}

func toUserData(u *model.User) *dto.UserData {
	return &dto.UserData{
		ID:                 u.ID,
		Username:           u.Username,
		FullName:           u.FullName,
		RoleID:             u.RoleID,
		RoleName:           u.RoleName,
		IsActive:           u.IsActive,
		MustChangePassword: u.MustChangePassword,
	}
}
