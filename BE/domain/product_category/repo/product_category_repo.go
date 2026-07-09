package repo

import (
	dto "pos_api/domain/product_category/dto"
	model "pos_api/domain/product_category/model"
	request_helper "pos_api/helper/request"
)

const (
	countCategoriesQuery       = `SELECT COUNT(*) FROM categories c WHERE 1=1`
	getAllCategoriesQuery       = `SELECT c.id, c.name, COALESCE(c.code, '') as code, c.description, COALESCE(c.is_active, 1) as is_active, COUNT(p.id) AS product_count, COUNT(CASE WHEN p.is_active = 1 THEN 1 END) AS active_product_count, c.created_at FROM categories c LEFT JOIN products p ON p.category_id = c.id WHERE 1=1`
	getAllCategoriesGroupBy     = ` GROUP BY c.id, c.name, c.code, c.description, c.is_active, c.created_at`
	getAllCategoriesDefaultSort = ` ORDER BY c.name`
	getAllCategoryOptionsQuery       = `SELECT id, name FROM categories WHERE is_active = 1 ORDER BY name`
	getCategoryByIDQuery             = `SELECT c.id, c.name, COALESCE(c.code, '') as code, c.description, COALESCE(c.is_active, 1) as is_active, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id) AS product_count, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id AND p.is_active = 1) AS active_product_count, c.created_at FROM categories c WHERE c.id = ? LIMIT 1`
	getCategoryByNameQuery           = `SELECT id, name, COALESCE(code, '') as code, description, COALESCE(is_active, 1) as is_active, created_at FROM categories WHERE name = ? LIMIT 1`
	checkCategoryNameQuery           = `SELECT id FROM categories WHERE name = ? AND id != ? LIMIT 1`
	checkCategoryCodeQuery           = `SELECT id FROM categories WHERE code = ? LIMIT 1`
	checkCategoryUsedQuery           = `SELECT COUNT(*) FROM products WHERE category_id = ?`
	checkCategoryActiveProductsQuery = `SELECT COUNT(*) FROM products WHERE category_id = ? AND is_active = 1`
	createCategoryQuery              = `INSERT INTO categories (name, code, description) VALUES (?, ?, ?)`
	getLastInsertIDQuery             = `SELECT LAST_INSERT_ID()`
	updateCategoryQuery              = `UPDATE categories SET name = ?, description = ?, updated_at = NOW() WHERE id = ?`
	deleteCategoryQuery              = `DELETE FROM categories WHERE id = ?`
	toggleCategoryStatusQuery        = `UPDATE categories SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *categoryRepo) GetAll(req *dto.GetAllRequest) ([]*model.Category, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND c.name LIKE ?`
		args = append(args, search)
	}

	if req.IsActive != nil {
		conditions += ` AND c.is_active = ?`
		args = append(args, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countCategoriesQuery+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortFields := map[string]string{
		"name":          "c.name",
		"product_count": "product_count",
		"is_active":     "c.is_active",
		"created_at":    "c.created_at",
	}
	query := getAllCategoriesQuery + conditions + getAllCategoriesGroupBy
	query += request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, getAllCategoriesDefaultSort)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Category
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *categoryRepo) GetOptions() ([]*model.CategoryOption, error) {
	var dataDB []*model.CategoryOption
	err := r.db.Raw(getAllCategoryOptionsQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *categoryRepo) GetByID(id int) (*model.Category, error) {
	var dataDB model.Category
	err := r.db.Raw(getCategoryByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *categoryRepo) Create(req *dto.CreateRequest) (int64, error) {
	err := r.db.Exec(createCategoryQuery, req.Name, req.Code, req.Description).Error
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.Raw(getLastInsertIDQuery).Scan(&id).Error
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *categoryRepo) Update(req *dto.UpdateRequest) error {
	err := r.db.Exec(updateCategoryQuery, req.Name, req.Description, req.ID).Error
	return err
}

func (r *categoryRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteCategoryQuery, req.ID).Error
	return err
}

func (r *categoryRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleCategoryStatusQuery, req.ID).Error
	return err
}

func (r *categoryRepo) GetByName(name string) (*model.Category, error) {
	var dataDB model.Category
	err := r.db.Raw(getCategoryByNameQuery, name).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *categoryRepo) CheckCodeExists(code string) (bool, error) {
	var id int
	err := r.db.Raw(checkCategoryCodeQuery, code).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *categoryRepo) CheckNameExists(name string, excludeID int) (bool, error) {
	var id int
	err := r.db.Raw(checkCategoryNameQuery, name, excludeID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *categoryRepo) CountProductsByCategory(categoryID int) (int, error) {
	var count int
	err := r.db.Raw(checkCategoryUsedQuery, categoryID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *categoryRepo) CountActiveProductsByCategory(categoryID int) (int, error) {
	var count int
	err := r.db.Raw(checkCategoryActiveProductsQuery, categoryID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
