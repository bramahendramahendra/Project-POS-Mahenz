package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"time"

	"portal_api/config"
	dto "portal_api/domain/auth/dto"
	model "portal_api/domain/auth/model"
	"portal_api/errors"
	time_helper "portal_api/helper/time"
	"portal_api/pkg/authclient"
	"portal_api/pkg/jwt"
)

// Login: username + password diverifikasi BE User (pesan penolakan diteruskan apa adanya),
// lalu Portal menyiapkan profil lokal dan sesinya sendiri.
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, ip string) (*dto.LoginResponse, error) {
	body, err := json.Marshal(dto.VerifyRequest{Username: req.Username, Password: req.Password, ClientIP: ip})
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(ctx, "/login/verify", 0, body)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}

	var account dto.VerifyUser
	if err := resp.DecodeData(&account); err != nil {
		return nil, authclient.ToServiceError(err)
	}
	if account.ID <= 0 || !account.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	profile, err := s.ensureProfile(&account)
	if err != nil {
		return nil, err
	}

	token, refreshToken, expiresAt, err := s.issueSession(profile, &account, req.DeviceInfo, ip)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         toUserData(&account),
	}, nil
}

func (s *authService) Logout(token string) error {
	return s.repo.DeleteSessionByToken(token)
}

// RefreshToken memperbarui sesi. Akun diperiksa ulang ke BE User, sehingga akun yang sudah
// dinonaktifkan atau dihapus tidak bisa memperpanjang sesi.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
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

	account, err := s.fetchProfile(ctx, session.AuthUserID)
	if err != nil {
		return nil, err
	}
	if !account.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	profile, err := s.ensureProfile(account)
	if err != nil {
		return nil, err
	}

	token, newRefreshToken, expiresAt, err := s.issueSession(profile, account, session.DeviceInfo, "")
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetMe meneruskan profil akun dari BE User (selalu data terbaru: role, wajib ganti password).
func (s *authService) GetMe(ctx context.Context, authUserID int) (*authclient.Response, error) {
	resp, err := s.client.Post(ctx, "/account/profile", authUserID, nil)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}
	return resp, nil
}

// ChangeMyPassword meneruskan ganti password sendiri ke BE User (password lama wajib benar).
// Berhasil: penanda wajib ganti password di sesi Portal dihapus, sesi tetap aktif.
func (s *authService) ChangeMyPassword(ctx context.Context, userID, authUserID int, body []byte) (*authclient.Response, error) {
	resp, err := s.client.Post(ctx, "/account/change-password", authUserID, body)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}
	if err := s.repo.ClearSessionMustChangePassword(userID); err != nil {
		return nil, err
	}
	return resp, nil
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

// fetchProfile membaca akun dari BE User. Akun yang ditolak BE User (nonaktif, dihapus)
// dianggap tidak terautentikasi.
func (s *authService) fetchProfile(ctx context.Context, authUserID int) (*dto.VerifyUser, error) {
	resp, err := s.client.Post(ctx, "/account/profile", authUserID, nil)
	if err != nil {
		var upstream *errors.UpstreamError
		if stderrors.As(err, &upstream) && upstream.HTTPCode == http.StatusUnauthorized {
			return nil, &errors.UnauthenticatedError{Message: upstream.Message}
		}
		return nil, authclient.ToServiceError(err)
	}
	var account dto.VerifyUser
	if err := resp.DecodeData(&account); err != nil {
		return nil, authclient.ToServiceError(err)
	}
	return &account, nil
}

// ensureProfile mencari profil lokal lewat auth_user_id. Bila belum ada: tautkan profil lama
// Portal dengan username sama yang belum tertaut (data sebelum BE User), atau buat baru.
// Username dan nama disalin dari BE User agar tampilan/log tetap terbaca.
func (s *authService) ensureProfile(account *dto.VerifyUser) (*model.User, error) {
	profile, err := s.repo.GetProfileByAuthUserID(account.ID)
	if err != nil {
		return nil, err
	}
	if profile != nil {
		if profile.Username != account.Username || profile.FullName != account.FullName {
			if err := s.repo.UpdateProfile(profile.ID, account.Username, account.FullName); err != nil {
				return nil, err
			}
		}
		return profile, nil
	}

	legacy, err := s.repo.GetUnlinkedProfileByUsername(account.Username)
	if err != nil {
		return nil, err
	}
	if legacy != nil {
		if err := s.repo.LinkProfile(legacy.ID, account.ID, account.Username, account.FullName); err != nil {
			return nil, err
		}
		return &model.User{ID: legacy.ID, AuthUserID: account.ID, Username: account.Username, FullName: account.FullName}, nil
	}

	id, err := s.repo.CreateProfile(account.ID, account.Username, account.FullName)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, AuthUserID: account.ID, Username: account.Username, FullName: account.FullName}, nil
}

// issueSession membuat token + sesi baru (satu sesi aktif per user, sesi lama diganti).
func (s *authService) issueSession(profile *model.User, account *dto.VerifyUser, device, ip string) (string, string, time.Time, error) {
	now := time_helper.GetTimeNow()
	expiresAt := now.Add(time.Second * time.Duration(config.Cfg.RefreshTokenExpire))

	claims := map[string]any{
		"user_id":      profile.ID,
		"auth_user_id": account.ID,
		"username":     account.Username,
		"full_name":    account.FullName,
		"role":         account.RoleName,
		"apps":         device,
	}
	jwt.CreateClaims(claims)
	token, err := jwt.GenerateToken()
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, err
	}

	if err := s.repo.DeleteSessionByUserID(profile.ID); err != nil {
		return "", "", time.Time{}, err
	}

	session := &model.Session{
		UserID:             profile.ID,
		UserRole:           account.RoleName,
		MustChangePassword: account.MustChangePassword,
		Token:              token,
		RefreshToken:       refreshToken,
		DeviceInfo:         device,
		IPAddress:          ip,
		ExpiresAt:          expiresAt,
		CreatedAt:          now,
	}
	if err := s.repo.CreateSession(session); err != nil {
		return "", "", time.Time{}, err
	}
	return token, refreshToken, expiresAt, nil
}

func toUserData(account *dto.VerifyUser) dto.UserData {
	return dto.UserData{
		ID:                 account.ID,
		Username:           account.Username,
		FullName:           account.FullName,
		RoleID:             account.RoleID,
		RoleName:           account.RoleName,
		MustChangePassword: account.MustChangePassword,
	}
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
