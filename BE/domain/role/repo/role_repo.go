package repo

import (
	dto "pos_api/domain/role/dto"
	model "pos_api/domain/role/model"
)

const (
	countRolesQuery      = `SELECT COUNT(*) FROM roles WHERE 1=1`
	getAllRolesQuery     = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE 1=1`
	getActiveRoleOptions = `SELECT id, display_name FROM roles WHERE is_active = 1 ORDER BY is_system DESC, display_name ASC`
	getRoleByIDQuery     = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE id = ? LIMIT 1`
	getRoleByNameQuery   = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE name = ? LIMIT 1`
	createRoleQuery      = `INSERT INTO roles (name, display_name, description) VALUES (?, ?, ?)`
	updateRoleQuery      = `UPDATE roles SET display_name = ?, description = ?, updated_at = NOW() WHERE id = ?`
	deleteRoleQuery      = `DELETE FROM roles WHERE id = ?`
	toggleRoleQuery      = `UPDATE roles SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *roleRepo) GetAll(req *dto.GetAllRequest) ([]*model.Role, int64, error) {
	var args, countArgs []any
	conditions := ""

	if req.Search != "" {
		like := "%" + req.Search + "%"
		conditions += ` AND (name LIKE ? OR display_name LIKE ?)`
		args = append(args, like, like)
		countArgs = append(countArgs, like, like)
	}
	if req.IsActive != nil {
		conditions += ` AND is_active = ?`
		args = append(args, *req.IsActive)
		countArgs = append(countArgs, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countRolesQuery+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := getAllRolesQuery + conditions + ` ORDER BY is_system DESC, id ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var dataDB []*model.Role
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *roleRepo) GetActiveOptions() ([]*dto.RoleOptionResponse, error) {
	var options []*dto.RoleOptionResponse
	if err := r.db.Raw(getActiveRoleOptions).Scan(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *roleRepo) GetByID(id int) (*model.Role, error) {
	var role model.Role
	if err := r.db.Raw(getRoleByIDQuery, id).Scan(&role).Error; err != nil {
		return nil, err
	}
	if role.ID == 0 {
		return nil, nil
	}
	return &role, nil
}

func (r *roleRepo) GetByName(name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Raw(getRoleByNameQuery, name).Scan(&role).Error; err != nil {
		return nil, err
	}
	if role.ID == 0 {
		return nil, nil
	}
	return &role, nil
}

func (r *roleRepo) Create(req *dto.CreateRequest) (int64, error) {
	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}
	if err := r.db.Exec(createRoleQuery, req.Name, req.DisplayName, desc).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *roleRepo) Update(id int, req *dto.UpdateRequest) error {
	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}
	return r.db.Exec(updateRoleQuery, req.DisplayName, desc, id).Error
}

func (r *roleRepo) Delete(id int) error {
	return r.db.Exec(deleteRoleQuery, id).Error
}

func (r *roleRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleRoleQuery, id).Error
}
