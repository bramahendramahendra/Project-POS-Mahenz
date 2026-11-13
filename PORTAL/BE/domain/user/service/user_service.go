package service

import (
	"context"
	"fmt"

	"portal_api/pkg/authclient"
)

func (s *userService) Forward(ctx context.Context, actingUserID int, path string, body []byte) (*authclient.Response, error) {
	resp, err := s.client.Post(ctx, path, actingUserID, body)
	if err != nil {
		return nil, authclient.ToServiceError(err)
	}
	return resp, nil
}

// ResetPassword: admin menetapkan password sementara. User wajib menggantinya saat login
// berikutnya, dan sesi Portal-nya dihapus.
func (s *userService) ResetPassword(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error) {
	resp, err := s.Forward(ctx, actingUserID, fmt.Sprintf("/users/reset-password/%d", id), body)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteSessionByAuthUserID(id); err != nil {
		return nil, err
	}
	return resp, nil
}

// ToggleStatus: user yang dinonaktifkan langsung dikeluarkan dari Portal. (User yang
// diaktifkan kembali tidak punya sesi, jadi penghapusan sesi tidak berdampak.)
func (s *userService) ToggleStatus(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error) {
	resp, err := s.Forward(ctx, actingUserID, fmt.Sprintf("/users/toggle-status/%d", id), body)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteSessionByAuthUserID(id); err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete: akun dihapus di BE User, profil dan sesi lokal Portal ikut dihapus.
func (s *userService) Delete(ctx context.Context, actingUserID, id int, body []byte) (*authclient.Response, error) {
	resp, err := s.Forward(ctx, actingUserID, fmt.Sprintf("/users/delete/%d", id), body)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteProfileByAuthUserID(id); err != nil {
		return nil, err
	}
	return resp, nil
}
