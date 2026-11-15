package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"coffee_api/config"
	dto "coffee_api/domain/auth/dto"
	model "coffee_api/domain/auth/model"
	"coffee_api/errors"
	time_helper "coffee_api/helper/time"
	"coffee_api/pkg/authclient"
	"coffee_api/pkg/jwt"
)

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, ip string) (*dto.LoginResponse, error) {
	if !s.directLoginEnabled() {
		return nil, &errors.UnauthorizededError{Message: "Login langsung dinonaktifkan. Silakan masuk lewat Portal"}
	}

	// Username + password diperiksa BE User (keputusan 66): salah, nonaktif, tanpa akses ke
	// Coffee, dan wajib ganti password ditolak di sana, pesannya diteruskan apa adanya.
	body, err := json.Marshal(dto.VerifyRequest{Username: req.Username, Password: req.Password, ClientIP: ip})
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(ctx, "/login/verify", 0, body)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}
	var account dto.AccountData
	if err := resp.DecodeData(&account); err != nil {
		return nil, authclient.ToServiceError(err)
	}
	if account.ID <= 0 || !account.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	user, err := s.resolveLocalUser(&account)
	if err != nil {
		return nil, err
	}
	return s.issueSession(user, req.DeviceInfo, ip)
}

// issueSession membuat token + sesi untuk user yang identitasnya sudah terverifikasi
// (dipakai login langsung dan login lewat Portal). Satu user satu sesi aktif.
func (s *authService) issueSession(user *model.User, deviceInfo, ip string) (*dto.LoginResponse, error) {
	now := time_helper.GetTimeNow()
	expiresAt := now.Add(time.Second * time.Duration(config.Cfg.RefreshTokenExpire))

	claims := map[string]any{
		"user_id":   user.ID,
		"username":  user.Username,
		"full_name": user.FullName,
		"role":      user.RoleName,
		"apps":      deviceInfo,
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
		DeviceInfo:   deviceInfo,
		IPAddress:    ip,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
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

	now := time_helper.GetTimeNow()
	expiresAt := now.Add(time.Second * time.Duration(config.Cfg.RefreshTokenExpire))

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
		CreatedAt:    now,
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
