package repo

import (
	model "coffee_api/domain/auth/model"
	time_helper "coffee_api/helper/time"
)

// Query profil lokal yang ditautkan ke akun BE User, dan toggle login langsung.
const (
	getUserByAuthIDQuery           = `SELECT u.id, u.auth_user_id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.auth_user_id = ? LIMIT 1`
	getUnlinkedUserByUsernameQuery = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.username = ? AND u.auth_user_id IS NULL LIMIT 1`
	usernameTakenQuery             = `SELECT COUNT(*) FROM users WHERE username = ?`
	linkUserQuery                  = `UPDATE users SET auth_user_id = ?, full_name = ?, updated_at = ? WHERE id = ? AND auth_user_id IS NULL`
	getRoleIDByNameQuery           = `SELECT id FROM roles WHERE name = ? LIMIT 1`
	createAuthUserQuery            = `INSERT INTO users (username, password, full_name, role_id, auth_user_id) VALUES (?, NULL, ?, ?, ?)`
	updateFullNameQuery            = `UPDATE users SET full_name = ?, updated_at = ? WHERE id = ?`

	getDirectLoginSettingQuery = `SELECT COALESCE(setting_value, '') FROM settings WHERE setting_key = 'direct_login_enabled' LIMIT 1`
	setDirectLoginSettingQuery = `INSERT INTO settings (setting_key, setting_value) VALUES ('direct_login_enabled', ?) ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_at = ?`
)

func (r *authRepo) GetUserByAuthID(authUserID int) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUserByAuthIDQuery, authUserID).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *authRepo) GetUnlinkedUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Raw(getUnlinkedUserByUsernameQuery, username).Scan(&user).Error; err != nil {
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

// LinkUser menautkan profil lokal yang belum tertaut (syarat auth_user_id IS NULL mencegah
// menimpa tautan yang dibuat permintaan lain di saat bersamaan).
func (r *authRepo) LinkUser(userID, authUserID int, fullName string) error {
	return r.db.Exec(linkUserQuery, authUserID, fullName, time_helper.GetTimeNow(), userID).Error
}

func (r *authRepo) GetRoleIDByName(name string) (int, error) {
	var id int
	if err := r.db.Raw(getRoleIDByNameQuery, name).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

// CreateAuthUser membuat profil lokal baru tanpa password lokal (password ada di BE User).
func (r *authRepo) CreateAuthUser(username, fullName string, roleID, authUserID int) error {
	return r.db.Exec(createAuthUserQuery, username, fullName, roleID, authUserID).Error
}

func (r *authRepo) UpdateFullName(userID int, fullName string) error {
	return r.db.Exec(updateFullNameQuery, fullName, time_helper.GetTimeNow(), userID).Error
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
