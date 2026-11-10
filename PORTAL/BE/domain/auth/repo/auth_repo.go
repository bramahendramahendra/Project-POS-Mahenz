package repo

import (
	model "portal_api/domain/auth/model"
	time_helper "portal_api/helper/time"
)

const (
	getUserByUsernameQuery        = `SELECT u.id, u.username, u.password, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.must_change_password FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.username = ? LIMIT 1`
	getUserByIDQuery              = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.must_change_password FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.id = ? LIMIT 1`
	getPasswordHashByIDQuery      = `SELECT password FROM users WHERE id = ? LIMIT 1`
	updateOwnPasswordQuery        = `UPDATE users SET password = ?, must_change_password = 0, updated_at = ? WHERE id = ?`
	createSessionQuery            = `INSERT INTO sessions (user_id, user_role, token, refresh_token, device_info, ip_address, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE user_role=VALUES(user_role), token=VALUES(token), refresh_token=VALUES(refresh_token), device_info=VALUES(device_info), ip_address=VALUES(ip_address), expires_at=VALUES(expires_at), created_at=VALUES(created_at)`
	getSessionByTokenQuery        = `SELECT s.id, s.user_id, s.user_role, s.token, s.device_info, s.expires_at, u.must_change_password FROM sessions s INNER JOIN users u ON u.id = s.user_id WHERE s.token = ? LIMIT 1`
	getSessionByRefreshTokenQuery = `SELECT id, user_id, refresh_token, expires_at FROM sessions WHERE refresh_token = ? LIMIT 1`
	deleteSessionByUserIDQuery    = `DELETE FROM sessions WHERE user_id = ?`
	deleteSessionByTokenQuery     = `DELETE FROM sessions WHERE token = ?`
)

func (r *authRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUserByUsernameQuery, username).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *authRepo) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUserByIDQuery, id).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *authRepo) GetPasswordHashByID(id int) (string, error) {
	var hash string
	if err := r.db.Raw(getPasswordHashByIDQuery, id).Scan(&hash).Error; err != nil {
		return "", err
	}
	return hash, nil
}

// UpdateOwnPassword mengganti password milik sendiri sekaligus menghapus
// penanda "wajib ganti password".
func (r *authRepo) UpdateOwnPassword(id int, hashedPassword string) error {
	return r.db.Exec(updateOwnPasswordQuery, hashedPassword, time_helper.GetTimeNow(), id).Error
}

func (r *authRepo) CreateSession(session *model.Session) error {
	return r.db.Exec(createSessionQuery,
		session.UserID, session.UserRole, session.Token, session.RefreshToken,
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

func (r *authRepo) DeleteSessionByUserID(userID int) error {
	return r.db.Exec(deleteSessionByUserIDQuery, userID).Error
}

func (r *authRepo) DeleteSessionByToken(token string) error {
	return r.db.Exec(deleteSessionByTokenQuery, token).Error
}
