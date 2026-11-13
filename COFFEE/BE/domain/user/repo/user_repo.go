package repo

import (
	dto "coffee_api/domain/user/dto"
	model "coffee_api/domain/user/model"
	request_helper "coffee_api/helper/request"
	time_helper "coffee_api/helper/time"
)

const (
	countUsersQuery        = `SELECT COUNT(*) FROM users u WHERE 1=1`
	getAllUsersQuery       = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.created_at, u.updated_at FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE 1=1`
	getUserByIDQuery       = `SELECT u.id, u.username, u.full_name, u.role_id, r.name AS role_name, u.is_active, u.created_at, u.updated_at FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE u.id = ? LIMIT 1`
	updateUserQuery        = `UPDATE users SET role_id = ?, updated_at = ? WHERE id = ?`
	deleteUserQuery        = `DELETE FROM users WHERE id = ?`
	toggleUserStatusQuery  = `UPDATE users SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteSessionQuery     = `DELETE FROM sessions WHERE user_id = ?`
	countActiveAdminsQuery = `SELECT COUNT(*) FROM users u INNER JOIN roles r ON r.id = u.role_id WHERE r.name IN ('owner', 'admin') AND u.is_active = 1 AND u.id != ?`
)

func (r *userRepo) GetAll(req *dto.GetAllRequest) ([]*model.User, int64, error) {
	conditions := ""
	var args []any

	if req.Search != "" {
		like := "%" + req.Search + "%"
		conditions += ` AND (u.username LIKE ? OR u.full_name LIKE ?)`
		args = append(args, like, like)
	}
	if req.RoleID != nil {
		conditions += ` AND u.role_id = ?`
		args = append(args, *req.RoleID)
	}
	if req.IsActive != nil {
		conditions += ` AND u.is_active = ?`
		args = append(args, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countUsersQuery+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortColumns := map[string]string{
		"username":  "u.username",
		"full_name": "u.full_name",
		"role_name": "r.name",
		"is_active": "u.is_active",
	}
	sortCol := "u.id"
	if col, ok := allowedSortColumns[req.SortBy]; ok {
		sortCol = col
	}
	sortDir := "ASC"
	if req.SortOrder == "desc" {
		sortDir = "DESC"
	}

	query := getAllUsersQuery + conditions + " ORDER BY " + sortCol + " " + sortDir + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.User
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	var u model.User
	if err := r.db.Raw(getUserByIDQuery, id).Scan(&u).Error; err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, nil
	}
	return &u, nil
}

func (r *userRepo) Update(id int, req *dto.UpdateRequest) error {
	return r.db.Exec(updateUserQuery, req.RoleID, time_helper.GetTimeNow(), id).Error
}

func (r *userRepo) Delete(id int) error {
	return r.db.Exec(deleteUserQuery, id).Error
}

func (r *userRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleUserStatusQuery, time_helper.GetTimeNow(), id).Error
}

func (r *userRepo) DeleteSessionByUserID(userID int) error {
	return r.db.Exec(deleteSessionQuery, userID).Error
}

func (r *userRepo) CountActiveAdmins(excludeID int) (int64, error) {
	var count int64
	if err := r.db.Raw(countActiveAdminsQuery, excludeID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
