package repo

import (
	model "portal_api/domain/auth/model"
	time_helper "portal_api/helper/time"
)

const (
	getProfileByAuthUserIDQuery       = `SELECT id, auth_user_id, username, full_name FROM users WHERE auth_user_id = ? LIMIT 1`
	getUnlinkedProfileByUsernameQuery = `SELECT id, username, full_name FROM users WHERE username = ? AND auth_user_id IS NULL ORDER BY id LIMIT 1`
	linkProfileQuery                  = `UPDATE users SET auth_user_id = ?, username = ?, full_name = ?, updated_at = ? WHERE id = ? AND auth_user_id IS NULL`
	updateProfileQuery                = `UPDATE users SET username = ?, full_name = ?, updated_at = ? WHERE id = ?`
	createProfileQuery                = `INSERT INTO users (auth_user_id, username, full_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	createSessionQuery                = `INSERT INTO sessions (user_id, user_role, must_change_password, token, refresh_token, device_info, ip_address, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE user_role=VALUES(user_role), must_change_password=VALUES(must_change_password), token=VALUES(token), refresh_token=VALUES(refresh_token), device_info=VALUES(device_info), ip_address=VALUES(ip_address), expires_at=VALUES(expires_at), created_at=VALUES(created_at)`
	getSessionByTokenQuery            = `SELECT s.id, s.user_id, u.auth_user_id, s.user_role, s.must_change_password, s.token, s.device_info, s.expires_at FROM sessions s INNER JOIN users u ON u.id = s.user_id WHERE s.token = ? AND u.auth_user_id IS NOT NULL LIMIT 1`
	getSessionByRefreshTokenQuery     = `SELECT s.id, s.user_id, u.auth_user_id, s.device_info, s.refresh_token, s.expires_at FROM sessions s INNER JOIN users u ON u.id = s.user_id WHERE s.refresh_token = ? AND u.auth_user_id IS NOT NULL LIMIT 1`
	clearSessionMustChangeQuery       = `UPDATE sessions SET must_change_password = 0 WHERE user_id = ?`
	deleteSessionByUserIDQuery        = `DELETE FROM sessions WHERE user_id = ?`
	deleteSessionByTokenQuery         = `DELETE FROM sessions WHERE token = ?`
)

func (r *authRepo) GetProfileByAuthUserID(authUserID int) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getProfileByAuthUserIDQuery, authUserID).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

// GetUnlinkedProfileByUsername: profil lama Portal (sebelum BE User) yang belum ditautkan.
func (r *authRepo) GetUnlinkedProfileByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUnlinkedProfileByUsernameQuery, username).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *authRepo) LinkProfile(id, authUserID int, username, fullName string) error {
	return r.db.Exec(linkProfileQuery, authUserID, username, fullName, time_helper.GetTimeNow(), id).Error
}

func (r *authRepo) UpdateProfile(id int, username, fullName string) error {
	return r.db.Exec(updateProfileQuery, username, fullName, time_helper.GetTimeNow(), id).Error
}

func (r *authRepo) CreateProfile(authUserID int, username, fullName string) (int, error) {
	now := time_helper.GetTimeNow()
	result := r.db.Exec(createProfileQuery, authUserID, username, fullName, now, now)
	if result.Error != nil {
		return 0, result.Error
	}
	var id int
	if err := r.db.Raw(`SELECT id FROM users WHERE auth_user_id = ? LIMIT 1`, authUserID).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepo) CreateSession(session *model.Session) error {
	return r.db.Exec(createSessionQuery,
		session.UserID, session.UserRole, session.MustChangePassword, session.Token, session.RefreshToken,
		session.DeviceInfo, session.IPAddress, session.ExpiresAt, session.CreatedAt,
	).Error
}

func (r *authRepo) GetSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	if err := r.db.Raw(getSessionByTokenQuery, token).Scan(&session).Error; err != nil {
		return nil, err
	}
	if session.ID == 0 {
		return nil, nil
	}
	return &session, nil
}

func (r *authRepo) GetSessionByRefreshToken(token string) (*model.Session, error) {
	var session model.Session
	if err := r.db.Raw(getSessionByRefreshTokenQuery, token).Scan(&session).Error; err != nil {
		return nil, err
	}
	if session.ID == 0 {
		return nil, nil
	}
	return &session, nil
}

func (r *authRepo) ClearSessionMustChangePassword(userID int) error {
	return r.db.Exec(clearSessionMustChangeQuery, userID).Error
}

func (r *authRepo) DeleteSessionByUserID(userID int) error {
	return r.db.Exec(deleteSessionByUserIDQuery, userID).Error
}

func (r *authRepo) DeleteSessionByToken(token string) error {
	return r.db.Exec(deleteSessionByTokenQuery, token).Error
}
