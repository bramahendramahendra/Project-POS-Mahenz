package service

import (
	"context"
	"encoding/json"
	"strings"

	dto "coffee_api/domain/auth/dto"
	model "coffee_api/domain/auth/model"
	"coffee_api/errors"
	"coffee_api/pkg/authclient"
)

// Role bawaan untuk profil lokal yang dibuat otomatis saat login pertama (lihat migrasi 011).
const pendingRoleName = "menunggu_akses"

// LoginSSO menukar kode sekali pakai (dibuat BE User saat kartu Coffee dibuka di Portal)
// menjadi sesi Coffee.
//
// Alur: tukar kode ke BE User (server ke server, memakai kunci aplikasi Coffee; BE User hanya
// menerima kode milik aplikasi pemanggil) -> cari/tautkan/buat profil lokal -> buat sesi
// seperti login biasa. Role tidak datang dari BE User, tetapi dari profil lokal.
func (s *authService) LoginSSO(ctx context.Context, req *dto.SSORequest, ip string) (*dto.LoginResponse, error) {
	body, err := json.Marshal(dto.ExchangeRequest{Code: strings.TrimSpace(req.Code)})
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(ctx, "/sso/exchange", 0, body)
	if err != nil {
		if _, rejected := err.(*errors.UpstreamError); rejected {
			return nil, &errors.UnauthenticatedError{Message: "Kode login tidak valid atau sudah kedaluwarsa. Silakan buka lagi dari Portal"}
		}
		return nil, authclient.ToServiceError(err)
	}

	var result dto.ExchangeResult
	if err := resp.DecodeData(&result); err != nil {
		return nil, authclient.ToServiceError(err)
	}
	if !result.Valid || result.User.ID <= 0 || !result.User.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Login lewat Portal ditolak"}
	}

	user, err := s.resolveLocalUser(&result.User)
	if err != nil {
		return nil, err
	}

	device := req.DeviceInfo
	if device == "" {
		device = "web"
	}
	return s.issueSession(user, device, ip)
}

// resolveLocalUser mencari profil lokal untuk akun yang SUDAH diverifikasi BE User:
//  1. lewat auth_user_id (nama disalin dari BE User, keputusan 76);
//  2. bila belum ada: tautkan ke profil lokal dengan username sama yang belum tertaut
//     (keputusan 67; aman karena identitas sudah diverifikasi BE User);
//  3. bila tidak ada: buat profil baru dengan role menunggu_akses.
//
// Profil lokal yang dinonaktifkan admin Coffee tetap ditolak.
func (s *authService) resolveLocalUser(account *dto.AccountData) (*model.User, error) {
	user, err := s.repo.GetUserByAuthID(account.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.linkOrCreateUser(account)
		if err != nil {
			return nil, err
		}
	} else if account.FullName != "" && account.FullName != user.FullName {
		if err := s.repo.UpdateFullName(user.ID, account.FullName); err != nil {
			return nil, err
		}
		user.FullName = account.FullName
	}

	if !user.IsActive {
		return nil, &errors.UnauthenticatedError{Message: "Akun tidak aktif di Coffee. Hubungi admin Coffee"}
	}
	return user, nil
}

func (s *authService) linkOrCreateUser(account *dto.AccountData) (*model.User, error) {
	existing, err := s.repo.GetUnlinkedUserByUsername(account.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fullName := account.FullName
		if fullName == "" {
			fullName = existing.FullName
		}
		if err := s.repo.LinkUser(existing.ID, account.ID, fullName); err != nil {
			return nil, err
		}
		return s.reloadLinkedUser(account.ID)
	}

	// Username sudah dipakai profil lokal yang tertaut ke akun BE User LAIN: jangan diambil alih.
	taken, err := s.repo.UsernameTaken(account.Username)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, &errors.UnauthenticatedError{
			Message: "Username ini sudah dipakai akun lain di Coffee. Hubungi admin Coffee",
		}
	}

	roleID, err := s.repo.GetRoleIDByName(pendingRoleName)
	if err != nil {
		return nil, err
	}
	if roleID == 0 {
		return nil, &errors.InternalServerError{Message: "Role awal untuk akun baru belum tersedia"}
	}

	if err := s.repo.CreateAuthUser(account.Username, account.FullName, roleID, account.ID); err != nil {
		// Dua permintaan pertama bersamaan bisa saling bentrok pada kunci unik; baca ulang.
		if existing, getErr := s.repo.GetUserByAuthID(account.ID); getErr == nil && existing != nil {
			return existing, nil
		}
		return nil, err
	}
	return s.reloadLinkedUser(account.ID)
}

func (s *authService) reloadLinkedUser(authUserID int) (*model.User, error) {
	user, err := s.repo.GetUserByAuthID(authUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &errors.InternalServerError{Message: "Gagal menyiapkan profil lokal"}
	}
	return user, nil
}

// accountServiceAvailable: layanan akun (BE User) sudah diatur (alamat dan kunci aplikasi ada).
func (s *authService) accountServiceAvailable() (bool, error) {
	return s.client.Configured()
}

// directLoginEnabled menentukan apakah form username + password boleh dipakai.
// Hanya nilai "false" di pengaturan yang mematikan (kosong/tidak ada/rusak = aktif).
// Bila layanan akun belum diatur, form tetap tampil (login akan dijawab "belum diatur"),
// sebab tombol Portal juga tidak bisa dipakai. Pemulihan bila terkunci: lewat SQL
// (keputusan 73).
func (s *authService) directLoginEnabled() bool {
	available, err := s.accountServiceAvailable()
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
	available, err := s.accountServiceAvailable()
	if err != nil {
		return nil, err
	}
	return &dto.LoginOptionsResponse{
		DirectLoginEnabled: s.directLoginEnabled(),
		PortalLoginEnabled: available,
	}, nil
}

func (s *authService) SetDirectLogin(enabled bool) error {
	if !enabled {
		// Mematikan login langsung hanya masuk akal bila login lewat Portal siap.
		available, err := s.accountServiceAvailable()
		if err != nil {
			return err
		}
		if !available {
			return &errors.BadRequestError{
				Message: "Layanan akun belum diatur di server ini, sehingga login langsung tidak dapat dimatikan",
			}
		}
	}

	value := "true"
	if !enabled {
		value = "false"
	}
	return s.repo.SetDirectLoginSetting(value)
}
