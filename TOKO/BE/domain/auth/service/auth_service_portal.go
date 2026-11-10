package service

import (
	"context"
	stderrors "errors"
	"strings"

	"pos_api/config"
	dto "pos_api/domain/auth/dto"
	model "pos_api/domain/auth/model"
	"pos_api/errors"
	"pos_api/pkg/bcrypt"
	"pos_api/pkg/portal"
)

// Role bawaan untuk profil lokal yang dibuat otomatis saat SSO pertama (lihat migrasi 011).
const pendingRoleName = "menunggu_akses"

// LoginSSO menukar kode sekali pakai dari Portal menjadi sesi Toko.
//
// Alur: tukar kode ke Portal (server ke server, memakai kunci aplikasi) -> pastikan jawaban
// untuk aplikasi ini -> cari profil lokal lewat portal_user_id (buat bila belum ada) -> buat
// sesi seperti login biasa. Role tidak datang dari Portal, tetapi dari profil lokal.
func (s *authService) LoginSSO(ctx context.Context, req *dto.SSORequest, ip string) (*dto.LoginResponse, error) {
	creds, err := s.repo.GetPortalCredentials()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, &errors.InternalServerError{Message: "Login lewat Portal belum diatur di server ini"}
	}

	result, err := portal.NewClient(config.Cfg.PortalUrl).Exchange(ctx, creds.AppID, creds.AppKey, strings.TrimSpace(req.Code))
	if err != nil {
		switch {
		case stderrors.Is(err, portal.ErrRejected):
			return nil, &errors.UnauthenticatedError{Message: "Kode login tidak valid atau sudah kedaluwarsa. Silakan buka lagi dari Portal"}
		default:
			return nil, &errors.InternalServerError{Message: "Portal tidak dapat dihubungi. Coba lagi beberapa saat"}
		}
	}

	// app_id yang dipercaya adalah yang datang dari jawaban Portal, bukan dari URL.
	if !result.Valid || result.AppID != creds.AppID || result.User.PortalUserID <= 0 || !result.User.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Login lewat Portal ditolak"}
	}

	user, err := s.repo.GetUserByPortalID(result.User.PortalUserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.provisionPortalUser(result.User.PortalUserID, result.User.Username, result.User.FullName)
		if err != nil {
			return nil, err
		}
	} else if result.User.FullName != "" && result.User.FullName != user.FullName {
		// Nama dikelola di Portal; role dan PIN tetap milik Toko.
		if err := s.repo.UpdateFullName(user.ID, result.User.FullName); err != nil {
			return nil, err
		}
		user.FullName = result.User.FullName
	}

	if !user.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif"}
	}

	device := req.DeviceInfo
	if device == "" {
		device = "web"
	}
	return s.issueSession(user, device, ip)
}

// provisionPortalUser membuat profil lokal untuk user Portal yang baru pertama kali masuk.
// Akun lokal dengan username sama yang belum ditautkan TIDAK otomatis dipakai: menautkan
// berdasarkan nama saja berisiko mengambil alih akun orang lain.
func (s *authService) provisionPortalUser(portalUserID int, username, fullName string) (*model.User, error) {
	taken, err := s.repo.UsernameTaken(username)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, &errors.UnauthenticatedError{
			Message: "Sudah ada akun Toko dengan username yang sama yang belum ditautkan ke Portal. Hubungi admin",
		}
	}

	roleID, err := s.repo.GetRoleIDByName(pendingRoleName)
	if err != nil {
		return nil, err
	}
	if roleID == 0 {
		return nil, &errors.InternalServerError{Message: "Role awal untuk akun Portal belum tersedia"}
	}

	// Password acak yang tidak diketahui siapa pun: akun ini hanya bisa masuk lewat Portal.
	unusable, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}
	hashed, err := bcrypt.HashPassword(unusable)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreatePortalUser(username, hashed, fullName, roleID, portalUserID); err != nil {
		// Dua permintaan pertama bersamaan bisa saling bentrok pada kunci unik; baca ulang.
		existing, getErr := s.repo.GetUserByPortalID(portalUserID)
		if getErr == nil && existing != nil {
			return existing, nil
		}
		return nil, err
	}

	created, err := s.repo.GetUserByPortalID(portalUserID)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal membuat profil lokal"}
	}
	return created, nil
}

// portalLoginAvailable: login lewat Portal siap dipakai (alamat Portal dan kredensial aplikasi ada).
func (s *authService) portalLoginAvailable() (bool, error) {
	if config.Cfg.PortalUrl == "" {
		return false, nil
	}
	creds, err := s.repo.GetPortalCredentials()
	if err != nil {
		return false, err
	}
	return creds != nil, nil
}

// directLoginEnabled menentukan apakah login langsung (username + password) boleh dipakai.
//
// Aturan pengaman agar tidak ada yang terkunci di luar:
//   - flag darurat di config server (DirectLoginForceEnabled) selalu menang;
//   - hanya nilai "false" di pengaturan yang mematikan (kosong/tidak ada/rusak = aktif);
//   - bila login lewat Portal belum siap, login langsung tetap aktif;
//   - bila pengaturan gagal dibaca, anggap aktif.
func (s *authService) directLoginEnabled() bool {
	if config.Cfg.DirectLoginForceEnabled {
		return true
	}

	available, err := s.portalLoginAvailable()
	if err != nil || !available {
		return true
	}

	value, err := s.repo.GetDirectLoginSetting()
	if err != nil {
		return true
	}
	return strings.TrimSpace(strings.ToLower(value)) != "false"
}

func (s *authService) GetLoginOptions() (*dto.LoginOptionsResponse, error) {
	available, err := s.portalLoginAvailable()
	if err != nil {
		return nil, err
	}
	return &dto.LoginOptionsResponse{
		DirectLoginEnabled: s.directLoginEnabled(),
		PortalLoginEnabled: available,
		DirectLoginForced:  config.Cfg.DirectLoginForceEnabled,
	}, nil
}

func (s *authService) SetDirectLogin(enabled bool) error {
	if !enabled {
		// Mematikan login langsung hanya masuk akal bila Portal siap; kalau tidak, semua orang terkunci.
		available, err := s.portalLoginAvailable()
		if err != nil {
			return err
		}
		if !available {
			return &errors.BadRequestError{
				Message: "Login lewat Portal belum diatur di server ini, sehingga login langsung tidak dapat dimatikan",
			}
		}
	}

	value := "true"
	if !enabled {
		value = "false"
	}
	return s.repo.SetDirectLoginSetting(value)
}
