package repo

import (
	"time"

	model "auth_api/domain/sso/model"
)

const (
	getLaunchAppQuery  = `SELECT a.id, a.code, a.path FROM apps a INNER JOIN user_apps ua ON ua.app_id = a.id WHERE ua.user_id = ? AND a.code = ? AND a.is_active = 1 AND a.show_in_portal = 1 LIMIT 1`
	createCodeQuery    = `INSERT INTO sso_codes (code_hash, user_id, app_id, expires_at, created_at) VALUES (?, ?, ?, ?, ?)`
	getCodeByHashQuery = `SELECT c.id, c.user_id, c.app_id, c.expires_at, c.used_at, a.code AS app_code, a.is_active AS app_active, u.username, u.full_name, u.is_active AS user_active FROM sso_codes c INNER JOIN apps a ON a.id = c.app_id INNER JOIN users u ON u.id = c.user_id WHERE c.code_hash = ? LIMIT 1`
	markUsedQuery      = `UPDATE sso_codes SET used_at = ? WHERE id = ? AND used_at IS NULL`
	hasAccessQuery     = `SELECT COUNT(*) FROM user_apps WHERE user_id = ? AND app_id = ?`
	deleteExpiredQuery = `DELETE FROM sso_codes WHERE expires_at < ?`
)

func (r *ssoRepo) GetLaunchApp(userID int, appCode string) (*model.LaunchApp, error) {
	var app model.LaunchApp
	if err := r.db.Raw(getLaunchAppQuery, userID, appCode).Scan(&app).Error; err != nil {
		return nil, err
	}
	if app.ID == 0 {
		return nil, nil
	}
	return &app, nil
}

func (r *ssoRepo) CreateCode(codeHash string, userID, appID int, expiresAt, now time.Time) error {
	return r.db.Exec(createCodeQuery, codeHash, userID, appID, expiresAt, now).Error
}

func (r *ssoRepo) GetCodeByHash(codeHash string) (*model.Code, error) {
	var code model.Code
	if err := r.db.Raw(getCodeByHashQuery, codeHash).Scan(&code).Error; err != nil {
		return nil, err
	}
	if code.ID == 0 {
		return nil, nil
	}
	return &code, nil
}

func (r *ssoRepo) MarkUsed(id int, now time.Time) (bool, error) {
	res := r.db.Exec(markUsedQuery, now, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (r *ssoRepo) HasAccess(userID, appID int) (bool, error) {
	var count int64
	if err := r.db.Raw(hasAccessQuery, userID, appID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ssoRepo) DeleteExpired(before time.Time) error {
	return r.db.Exec(deleteExpiredQuery, before).Error
}
