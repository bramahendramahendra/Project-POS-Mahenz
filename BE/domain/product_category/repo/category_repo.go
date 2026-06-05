package repo

import (
	"fmt"
	"strings"
	"unicode"

	dto "pos_api/domain/product_category/dto"
	model "pos_api/domain/product_category/model"
)

const (
	getAllCategoriesQuery            = `SELECT c.id, c.name, COALESCE(c.code, '') as code, c.description, COALESCE(c.is_active, 1) as is_active, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id) AS product_count, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id AND p.is_active = 1) AS active_product_count, c.created_at FROM categories c ORDER BY c.name`
	getCategoryByIDQuery             = `SELECT c.id, c.name, COALESCE(c.code, '') as code, c.description, COALESCE(c.is_active, 1) as is_active, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id) AS product_count, (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id AND p.is_active = 1) AS active_product_count, c.created_at FROM categories c WHERE c.id = ? LIMIT 1`
	getCategoryByNameQuery           = `SELECT id, name, COALESCE(code, '') as code, description, COALESCE(is_active, 1) as is_active, created_at FROM categories WHERE name = ? LIMIT 1`
	checkCategoryNameQuery           = `SELECT id FROM categories WHERE name = ? AND id != ? LIMIT 1`
	checkCategoryCodeQuery           = `SELECT id FROM categories WHERE code = ? LIMIT 1`
	checkCategoryUsedQuery           = `SELECT COUNT(*) FROM products WHERE category_id = ?`
	checkCategoryActiveProductsQuery = `SELECT COUNT(*) FROM products WHERE category_id = ? AND is_active = 1`
	createCategoryQuery              = `INSERT INTO categories (name, code, description) VALUES (?, ?, ?)`
	updateCategoryQuery              = `UPDATE categories SET name = ?, description = ?, updated_at = NOW() WHERE id = ?`
	deleteCategoryQuery              = `DELETE FROM categories WHERE id = ?`
	toggleCategoryStatusQuery        = `UPDATE categories SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *categoryRepo) GetAll() ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Raw(getAllCategoriesQuery).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) GetByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Raw(getCategoryByNameQuery, name).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) GetByID(id int) (*model.Category, error) {
	var category model.Category
	err := r.db.Raw(getCategoryByIDQuery, id).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) CheckNameExists(name string, excludeID int) (bool, error) {
	var id int
	result := r.db.Raw(checkCategoryNameQuery, name, excludeID).Scan(&id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *categoryRepo) CountProductsByCategory(categoryID int) (int, error) {
	var count int
	if err := r.db.Raw(checkCategoryUsedQuery, categoryID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *categoryRepo) CountActiveProductsByCategory(categoryID int) (int, error) {
	var count int
	if err := r.db.Raw(checkCategoryActiveProductsQuery, categoryID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *categoryRepo) GenerateUniqueCode(name string) (string, error) {
	letters := strings.Map(func(ru rune) rune {
		if unicode.IsLetter(ru) {
			return unicode.ToUpper(ru)
		}
		return -1
	}, name)

	base := letters
	if len(base) > 3 {
		base = base[:3]
	}
	for len(base) < 3 {
		base += "X"
	}

	candidate := base
	for i := 2; i <= 99; i++ {
		exists, err := r.CheckCodeExists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s%d", base, i)
	}
	return "", fmt.Errorf("tidak bisa generate kode kategori yang unik")
}

func (r *categoryRepo) CheckCodeExists(code string) (bool, error) {
	var id int
	result := r.db.Raw(checkCategoryCodeQuery, code).Scan(&id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *categoryRepo) Create(req *dto.CreateCategoryRequest) (int64, error) {
	err := r.db.Exec(createCategoryQuery, req.Name, req.Code, req.Description).Error
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *categoryRepo) Update(id int, name, description string) error {
	return r.db.Exec(updateCategoryQuery, name, description, id).Error
}

func (r *categoryRepo) Delete(id int) error {
	return r.db.Exec(deleteCategoryQuery, id).Error
}

func (r *categoryRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleCategoryStatusQuery, id).Error
}
