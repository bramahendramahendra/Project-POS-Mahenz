package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	dto "auth_api/domain/sso/dto"
	"auth_api/errors"
	time_helper "auth_api/helper/time"
)

// Masa berlaku kode SSO. Kode hanya perlu hidup sampai browser sampai di aplikasi tujuan.
const codeTTL = 60 * time.Second

// Pesan sengaja seragam agar penyerang tidak bisa membedakan penyebab penolakan.
const exchangeRejected = "Permintaan SSO tidak valid"

func (s *ssoService) Launch(userID int, req *dto.LaunchRequest) (*dto.LaunchResponse, error) {
	app, err := s.repo.GetLaunchApp(userID, req.AppCode)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke aplikasi ini"}
	}

	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	now := time_helper.GetTimeNow()
	if err := s.repo.CreateCode(hashHex(code), userID, app.ID, now.Add(codeTTL), now); err != nil {
		return nil, err
	}

	// Pembersihan kode lama secara oportunistik (gagal tidak menggagalkan login).
	_ = s.repo.DeleteExpired(now.Add(-24 * time.Hour))

	return &dto.LaunchResponse{
		RedirectURL:      buildRedirectURL(app.Path, code),
		ExpiresInSeconds: int(codeTTL.Seconds()),
	}, nil
}

// Exchange: kunci aplikasi sudah diperiksa middleware AppKey; di sini kode harus milik
// aplikasi pemanggil itu sendiri (kode untuk Toko tidak bisa ditukar oleh Coffee).
func (s *ssoService) Exchange(callerAppID int, req *dto.ExchangeRequest) (*dto.ExchangeResponse, error) {
	reject := &errors.UnauthenticatedError{Message: exchangeRejected}

	row, err := s.repo.GetCodeByHash(hashHex(strings.TrimSpace(req.Code)))
	if err != nil {
		return nil, err
	}
	if row == nil || row.AppID != callerAppID {
		return nil, reject
	}

	now := time_helper.GetTimeNow()
	if row.UsedAt != nil || now.After(row.ExpiresAt) {
		return nil, reject
	}
	if !row.AppActive || !row.UserActive {
		return nil, reject
	}

	hasAccess, err := s.repo.HasAccess(row.UserID, row.AppID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, reject
	}

	// Penandaan terpakai bersifat atomik: bila dua permintaan datang bersamaan, hanya satu yang lolos.
	ok, err := s.repo.MarkUsed(row.ID, now)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, reject
	}

	return &dto.ExchangeResponse{
		Valid: true,
		AppID: row.AppCode,
		User: dto.ExchangeUser{
			ID:       row.UserID,
			Username: row.Username,
			FullName: row.FullName,
			IsActive: row.UserActive,
		},
	}, nil
}

func generateCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashHex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// buildRedirectURL membentuk alamat tujuan: <path aplikasi>sso?code=<kode>.
// path boleh berupa path di domain yang sama ("/toko/") atau URL penuh.
func buildRedirectURL(appPath, code string) string {
	base := appPath
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base + "sso?code=" + code
}
