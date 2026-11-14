package repo

import (
	dto "coffee_api/domain/menu_category/dto"
	model "coffee_api/domain/menu_category/model"
	time_helper "coffee_api/helper/time"
)

const (
	categorySelect = `SELECT c.id, c.name, c.sort_order, c.is_active,
		(SELECT COUNT(*) FROM menu_items mi WHERE mi.category_id = c.id) AS item_count
		FROM menu_categories c`
	categoryOrder       = ` ORDER BY c.sort_order ASC, c.name ASC`
	getCategoryOptions  = `SELECT c.id, c.name FROM menu_categories c WHERE c.is_active = 1` + categoryOrder
	getCategoryByID     = categorySelect + ` WHERE c.id = ? LIMIT 1`
	categoryNameExists  = `SELECT COUNT(*) FROM menu_categories WHERE name = ? AND id <> ?`
	createCategoryQuery = `INSERT INTO menu_categories (name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?)`
	updateCategoryQuery = `UPDATE menu_categories SET name = ?, sort_order = ?, updated_at = ? WHERE id = ?`
	toggleCategoryQuery = `UPDATE menu_categories SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteCategoryQuery = `DELETE FROM menu_categories WHERE id = ?`
)

func (r *menuCategoryRepo) GetAll(req *dto.GetAllRequest) ([]*model.MenuCategory, error) {
	query := categorySelect + ` WHERE 1=1`
	var args []any
	if req.Search != "" {
		query += ` AND c.name LIKE ?`
		args = append(args, "%"+req.Search+"%")
	}
	if req.IsActive != nil {
		query += ` AND c.is_active = ?`
		args = append(args, *req.IsActive)
	}
	var data []*model.MenuCategory
	if err := r.db.Raw(query+categoryOrder, args...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuCategoryRepo) GetOptions() ([]*model.MenuCategory, error) {
	var data []*model.MenuCategory
	if err := r.db.Raw(getCategoryOptions).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuCategoryRepo) GetByID(id int) (*model.MenuCategory, error) {
	var data model.MenuCategory
	if err := r.db.Raw(getCategoryByID, id).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *menuCategoryRepo) NameExists(name string, excludeID int) (bool, error) {
	var count int64
	if err := r.db.Raw(categoryNameExists, name, excludeID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *menuCategoryRepo) Create(req *dto.SaveRequest) error {
	now := time_helper.GetTimeNow()
	return r.db.Exec(createCategoryQuery, req.Name, req.SortOrder, now, now).Error
}

func (r *menuCategoryRepo) Update(req *dto.SaveRequest) error {
	return r.db.Exec(updateCategoryQuery, req.Name, req.SortOrder, time_helper.GetTimeNow(), req.ID).Error
}

func (r *menuCategoryRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleCategoryQuery, time_helper.GetTimeNow(), id).Error
}

func (r *menuCategoryRepo) Delete(id int) error {
	return r.db.Exec(deleteCategoryQuery, id).Error
}
