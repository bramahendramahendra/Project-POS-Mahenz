package repo

import (
	model "auth_api/domain/auth/model"
	time_helper "auth_api/helper/time"
)

const (
	getUserByUsernameQuery   = `SELECT u.id, u.username, u.password, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.must_change_password FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.username = ? LIMIT 1`
	getUserByIDQuery         = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.must_change_password FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.id = ? LIMIT 1`
	getPasswordHashByIDQuery = `SELECT password FROM users WHERE id = ? LIMIT 1`
	updateOwnPasswordQuery   = `UPDATE users SET password = ?, must_change_password = 0, updated_at = ? WHERE id = ?`
	hasAccessQuery           = `SELECT COUNT(*) FROM user_apps WHERE user_id = ? AND app_id = ?`
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

func (r *authRepo) HasAccess(userID, appID int) (bool, error) {
	var count int64
	if err := r.db.Raw(hasAccessQuery, userID, appID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
