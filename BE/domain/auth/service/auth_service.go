package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"pos_api/config"
	dto "pos_api/domain/auth/dto"
	model "pos_api/domain/auth/model"
	"pos_api/errors"
	time_helper "pos_api/helper/time"
	"pos_api/pkg/bcrypt"
	"pos_api/pkg/jwt"
)

func (s *authService) Login(req *dto.LoginRequest, ip string) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &errors.UnauthenticatedError{Message: "Username atau password salah"}
	}
	if !user.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}
	if !bcrypt.VerifyPassword(req.Password, user.Password) {
		return nil, &errors.UnauthenticatedError{Message: "Username atau password salah"}
	}

	expiresAt := time_helper.GetTimeNow().Add(time.Second * time.Duration(config.Cfg.TokenExpire))

	claims := map[string]any{
		"user_id":   user.ID,
		"username":  user.Username,
		"full_name": user.FullName,
		"role":      user.RoleName,
		"apps":      req.DeviceInfo,
	}
	jwt.CreateClaims(claims)
	token, err := jwt.GenerateToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteSessionByUserID(user.ID); err != nil {
		return nil, err
	}

	session := &model.Session{
		UserID:       user.ID,
		UserRole:     user.RoleName,
		Token:        token,
		RefreshToken: refreshToken,
		DeviceInfo:   req.DeviceInfo,
		IPAddress:    ip,
		ExpiresAt:    expiresAt,
	}
	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: dto.UserData{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			RoleID:   user.RoleID,
			RoleName: user.RoleName,
		},
	}, nil
}

func (s *authService) Logout(token string) error {
	return s.repo.DeleteSessionByToken(token)
}

func (s *authService) RefreshToken(refreshToken string) (*dto.RefreshResponse, error) {
	session, err := s.repo.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, &errors.UnauthenticatedError{Message: "Refresh token tidak valid"}
	}
	if time_helper.GetTimeNow().After(session.ExpiresAt) {
		return nil, &errors.UnauthenticatedError{Message: "Refresh token sudah expired"}
	}

	user, err := s.repo.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	expiresAt := time_helper.GetTimeNow().Add(time.Second * time.Duration(config.Cfg.TokenExpire))

	claims := map[string]any{
		"user_id":   user.ID,
		"username":  user.Username,
		"full_name": user.FullName,
		"role":      user.RoleName,
		"apps":      session.DeviceInfo,
	}
	jwt.CreateClaims(claims)
	newToken, err := jwt.GenerateToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteSessionByUserID(user.ID); err != nil {
		return nil, err
	}

	newSession := &model.Session{
		UserID:       user.ID,
		UserRole:     user.RoleName,
		Token:        newToken,
		RefreshToken: newRefreshToken,
		DeviceInfo:   session.DeviceInfo,
		ExpiresAt:    expiresAt,
	}
	if err := s.repo.CreateSession(newSession); err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) GetMe(userID int) (*dto.UserData, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &errors.NotFoundError{Message: "User tidak ditemukan"}
	}
	return &dto.UserData{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		RoleID:   user.RoleID,
		RoleName: user.RoleName,
	}, nil
}

func (s *authService) VerifyToken(token string) (*dto.VerifyTokenResponse, error) {
	claims, err := jwt.VerifyToken(token)
	if err != nil {
		return &dto.VerifyTokenResponse{Valid: false, Error: err.Error()}, nil
	}

	claimsMap := make(map[string]any)
	for k, v := range *claims {
		claimsMap[k] = v
	}

	expReadable := ""
	if exp, ok := (*claims)["exp"].(float64); ok {
		expReadable = time.Unix(int64(exp), 0).Format(time.RFC3339)
	}

	return &dto.VerifyTokenResponse{
		Valid:       true,
		Claims:      claimsMap,
		ExpReadable: expReadable,
	}, nil
}

func (s *authService) ValidateToken(token string) (*model.Session, error) {
	if _, err := jwt.VerifyToken(token); err != nil {
		return nil, err
	}

	session, err := s.repo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, &errors.UnauthenticatedError{Message: "Token tidak valid atau sudah logout"}
	}
	if time_helper.GetTimeNow().After(session.ExpiresAt) {
		return nil, &errors.UnauthenticatedError{Message: "Token expired"}
	}

	return session, nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
