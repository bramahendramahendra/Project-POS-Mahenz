package repo

import (
	dto "coffee_api/domain/menu_size/dto"
	model "coffee_api/domain/menu_size/model"
	time_helper "coffee_api/helper/time"
)

const (
	sizeSelect = `SELECT c.id, c.name, c.sort_order, c.is_active,
		(SELECT COUNT(*) FROM menu_item_variants mv WHERE mv.size_id = c.id) AS usage_count
		FROM menu_sizes c`
	sizeOrder       = ` ORDER BY c.sort_order ASC, c.name ASC`
	getSizeOptions  = `SELECT c.id, c.name FROM menu_sizes c WHERE c.is_active = 1` + sizeOrder
	getSizeByID     = sizeSelect + ` WHERE c.id = ? LIMIT 1`
	sizeNameExists  = `SELECT COUNT(*) FROM menu_sizes WHERE name = ? AND id <> ?`
	createSizeQuery = `INSERT INTO menu_sizes (name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?)`
	updateSizeQuery = `UPDATE menu_sizes SET name = ?, sort_order = ?, updated_at = ? WHERE id = ?`
	toggleSizeQuery = `UPDATE menu_sizes SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteSizeQuery = `DELETE FROM menu_sizes WHERE id = ?`
)

func (r *menuSizeRepo) GetAll(req *dto.GetAllRequest) ([]*model.MenuSize, error) {
	query := sizeSelect + ` WHERE 1=1`
	var args []any
	if req.Search != "" {
		query += ` AND c.name LIKE ?`
		args = append(args, "%"+req.Search+"%")
	}
	if req.IsActive != nil {
		query += ` AND c.is_active = ?`
		args = append(args, *req.IsActive)
	}
	var data []*model.MenuSize
	if err := r.db.Raw(query+sizeOrder, args...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuSizeRepo) GetOptions() ([]*model.MenuSize, error) {
	var data []*model.MenuSize
	if err := r.db.Raw(getSizeOptions).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuSizeRepo) GetByID(id int) (*model.MenuSize, error) {
	var data model.MenuSize
	if err := r.db.Raw(getSizeByID, id).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *menuSizeRepo) NameExists(name string, excludeID int) (bool, error) {
	var count int64
	if err := r.db.Raw(sizeNameExists, name, excludeID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *menuSizeRepo) Create(req *dto.SaveRequest) error {
	now := time_helper.GetTimeNow()
	return r.db.Exec(createSizeQuery, req.Name, req.SortOrder, now, now).Error
}

func (r *menuSizeRepo) Update(req *dto.SaveRequest) error {
	return r.db.Exec(updateSizeQuery, req.Name, req.SortOrder, time_helper.GetTimeNow(), req.ID).Error
}

func (r *menuSizeRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleSizeQuery, time_helper.GetTimeNow(), id).Error
}

func (r *menuSizeRepo) Delete(id int) error {
	return r.db.Exec(deleteSizeQuery, id).Error
}
