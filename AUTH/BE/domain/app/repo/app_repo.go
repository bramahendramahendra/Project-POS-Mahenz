package repo

import (
	model "auth_api/domain/app/model"
	time_helper "auth_api/helper/time"

	"gorm.io/gorm"
)

const (
	// Hanya aplikasi yang tampil sebagai kartu (Portal sendiri tidak: show_in_portal = 0).
	getMyAppsQuery = `SELECT a.id, a.code, a.name, COALESCE(a.description, '') AS description, a.path, a.icon, a.sort_order, a.is_active FROM apps a INNER JOIN user_apps ua ON ua.app_id = a.id WHERE ua.user_id = ? AND a.is_active = 1 AND a.show_in_portal = 1 ORDER BY a.sort_order ASC, a.name ASC`

	// Pengenal aplikasi pemanggil dari hash kunci (middleware AppKey).
	getActiveByKeyHashQuery = `SELECT id, code FROM apps WHERE app_key_hash = ? AND is_active = 1 LIMIT 1`

	adminSelect = `SELECT a.id, a.code, a.name, COALESCE(a.description, '') AS description, a.path, a.icon, a.sort_order, a.is_active,
		(a.app_key_hash IS NOT NULL AND a.app_key_hash <> '') AS has_key,
		(SELECT COUNT(*) FROM user_apps ua WHERE ua.app_id = a.id) AS user_count,
		a.created_at FROM apps a`

	getAllAdminQuery    = adminSelect + ` ORDER BY a.sort_order ASC, a.name ASC`
	getByIDAdminQuery   = adminSelect + ` WHERE a.id = ? LIMIT 1`
	getByCodeAdminQuery = adminSelect + ` WHERE a.code = ? LIMIT 1`

	createAppQuery       = `INSERT INTO apps (code, name, description, path, icon, sort_order, app_key_hash) VALUES (?, ?, ?, ?, ?, ?, ?)`
	updateAppQuery       = `UPDATE apps SET name = ?, description = ?, path = ?, icon = ?, sort_order = ?, updated_at = ? WHERE id = ?`
	toggleAppStatusQuery = `UPDATE apps SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	setKeyHashQuery      = `UPDATE apps SET app_key_hash = ?, updated_at = ? WHERE id = ?`

	getUserAccessQuery  = `SELECT a.id AS app_id, a.code, a.name, a.is_active, (ua.user_id IS NOT NULL) AS granted FROM apps a LEFT JOIN user_apps ua ON ua.app_id = a.id AND ua.user_id = ? ORDER BY a.sort_order ASC, a.name ASC`
	deleteUserAppsQuery = `DELETE FROM user_apps WHERE user_id = ?`
	insertUserAppQuery  = `INSERT INTO user_apps (user_id, app_id) SELECT ?, id FROM apps WHERE id = ?`
	userExistsQuery     = `SELECT COUNT(*) FROM users WHERE id = ?`
)

// GetActiveByKeyHash: aplikasi aktif pemilik kunci (hash SHA-256 hex). nil bila tidak ada.
func (r *appRepo) GetActiveByKeyHash(keyHash string) (*model.App, error) {
	var app model.App
	if err := r.db.Raw(getActiveByKeyHashQuery, keyHash).Scan(&app).Error; err != nil {
		return nil, err
	}
	if app.ID == 0 {
		return nil, nil
	}
	return &app, nil
}

func (r *appRepo) GetMyApps(userID int) ([]*model.App, error) {
	var rows []*model.App
	if err := r.db.Raw(getMyAppsQuery, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *appRepo) GetAll() ([]*model.AppAdmin, error) {
	var rows []*model.AppAdmin
	if err := r.db.Raw(getAllAdminQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *appRepo) GetByID(id int) (*model.AppAdmin, error) {
	var row model.AppAdmin
	if err := r.db.Raw(getByIDAdminQuery, id).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &row, nil
}

func (r *appRepo) GetByCode(code string) (*model.AppAdmin, error) {
	var row model.AppAdmin
	if err := r.db.Raw(getByCodeAdminQuery, code).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &row, nil
}

func (r *appRepo) Create(app *model.App, keyHash string) (int64, error) {
	if err := r.db.Exec(createAppQuery, app.Code, app.Name, app.Description, app.Path, app.Icon, app.SortOrder, keyHash).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *appRepo) Update(id int, app *model.App) error {
	return r.db.Exec(updateAppQuery, app.Name, app.Description, app.Path, app.Icon, app.SortOrder, time_helper.GetTimeNow(), id).Error
}

func (r *appRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleAppStatusQuery, time_helper.GetTimeNow(), id).Error
}

func (r *appRepo) SetKeyHash(id int, keyHash string) error {
	return r.db.Exec(setKeyHashQuery, keyHash, time_helper.GetTimeNow(), id).Error
}

func (r *appRepo) GetUserAccess(userID int) ([]*model.UserAccess, error) {
	var rows []*model.UserAccess
	if err := r.db.Raw(getUserAccessQuery, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *appRepo) SetUserAccess(userID int, appIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(deleteUserAppsQuery, userID).Error; err != nil {
			return err
		}
		for _, appID := range appIDs {
			if err := tx.Exec(insertUserAppQuery, userID, appID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *appRepo) UserExists(userID int) (bool, error) {
	var count int64
	if err := r.db.Raw(userExistsQuery, userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
