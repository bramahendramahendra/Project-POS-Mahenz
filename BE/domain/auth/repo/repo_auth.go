package repo

import model "pos_api/domain/auth/model"

const (
	getUserByUsernameQuery        = `SELECT u.id, u.username, u.password, u.full_name, u.role_id, r.name AS role_name, u.is_active FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.username = ? LIMIT 1`
	getUserByIDQuery              = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.id = ? LIMIT 1`
	createSessionQuery            = `INSERT INTO sessions (user_id, user_role, token, refresh_token, device_info, ip_address, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE user_role=VALUES(user_role), token=VALUES(token), refresh_token=VALUES(refresh_token), device_info=VALUES(device_info), ip_address=VALUES(ip_address), expires_at=VALUES(expires_at), created_at=NOW()`
	getSessionByTokenQuery        = `SELECT id, user_id, user_role, token, device_info, expires_at FROM sessions WHERE token = ? LIMIT 1`
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

func (r *authRepo) CreateSession(session *model.Session) error {
	return r.db.Exec(createSessionQuery,
		session.UserID, session.UserRole, session.Token, session.RefreshToken,
		session.DeviceInfo, session.IPAddress, session.ExpiresAt,
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
