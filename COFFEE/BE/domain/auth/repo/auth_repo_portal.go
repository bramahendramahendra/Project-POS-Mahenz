package repo

import (
	model "coffee_api/domain/auth/model"
	time_helper "coffee_api/helper/time"
)

// Query untuk integrasi login terpusat (Portal SSO).
const (
	getPortalCredentialsQuery = `SELECT app_id, app_key FROM portal_credentials WHERE id = 1 LIMIT 1`
	getUserByPortalIDQuery    = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.portal_user_id = ? LIMIT 1`
	usernameTakenQuery        = `SELECT COUNT(*) FROM users WHERE username = ?`
	getRoleIDByNameQuery      = `SELECT id FROM roles WHERE name = ? LIMIT 1`
	createPortalUserQuery     = `INSERT INTO users (username, password, full_name, role_id, portal_user_id) VALUES (?, ?, ?, ?, ?)`
	updateFullNameQuery       = `UPDATE users SET full_name = ?, updated_at = ? WHERE id = ?`

	getDirectLoginSettingQuery = `SELECT COALESCE(setting_value, '') FROM settings WHERE setting_key = 'direct_login_enabled' LIMIT 1`
	setDirectLoginSettingQuery = `INSERT INTO settings (setting_key, setting_value) VALUES ('direct_login_enabled', ?) ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_at = ?`
)

func (r *authRepo) GetPortalCredentials() (*model.PortalCredentials, error) {
	var creds model.PortalCredentials
	if err := r.db.Raw(getPortalCredentialsQuery).Scan(&creds).Error; err != nil {
		return nil, err
	}
	if creds.AppID == "" || creds.AppKey == "" {
		return nil, nil
	}
	return &creds, nil
}

func (r *authRepo) GetUserByPortalID(portalUserID int) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUserByPortalIDQuery, portalUserID).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *authRepo) UsernameTaken(username string) (bool, error) {
	var count int64
	if err := r.db.Raw(usernameTakenQuery, username).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *authRepo) GetRoleIDByName(name string) (int, error) {
	var id int
	if err := r.db.Raw(getRoleIDByNameQuery, name).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepo) CreatePortalUser(username, hashedPassword, fullName string, roleID, portalUserID int) error {
	return r.db.Exec(createPortalUserQuery, username, hashedPassword, fullName, roleID, portalUserID).Error
}

func (r *authRepo) GetDirectLoginSetting() (string, error) {
	var value string
	if err := r.db.Raw(getDirectLoginSettingQuery).Scan(&value).Error; err != nil {
		return "", err
	}
	return value, nil
}

func (r *authRepo) SetDirectLoginSetting(value string) error {
	return r.db.Exec(setDirectLoginSettingQuery, value, time_helper.GetTimeNow()).Error
}

func (r *authRepo) UpdateFullName(userID int, fullName string) error {
	return r.db.Exec(updateFullNameQuery, fullName, time_helper.GetTimeNow(), userID).Error
}
