package repo

import (
	dto "auth_api/domain/role/dto"
	model "auth_api/domain/role/model"
)

const (
	countRolesQuery      = `SELECT COUNT(*) FROM roles WHERE 1=1`
	getAllRolesQuery     = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE 1=1`
	getActiveRoleOptions = `SELECT id, display_name FROM roles WHERE is_active = 1 ORDER BY is_system DESC, display_name ASC`
	getRoleByIDQuery     = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE id = ? LIMIT 1`
	getRoleByNameQuery   = `SELECT id, name, display_name, description, is_system, is_active, created_at, updated_at FROM roles WHERE name = ? LIMIT 1`
	createRoleQuery      = `INSERT INTO roles (name, display_name, description) VALUES (?, ?, ?)`
	updateRoleQuery      = `UPDATE roles SET display_name = ?, description = ?, updated_at = ? WHERE id = ?`
	deleteRoleQuery      = `DELETE FROM roles WHERE id = ?`
	toggleRoleQuery      = `UPDATE roles SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
)

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
