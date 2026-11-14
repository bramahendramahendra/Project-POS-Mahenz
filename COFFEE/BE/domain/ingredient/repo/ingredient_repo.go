package repo

import (
	dto "coffee_api/domain/ingredient/dto"
	model "coffee_api/domain/ingredient/model"
	request_helper "coffee_api/helper/request"
	time_helper "coffee_api/helper/time"
)

const (
	ingredientSelect = `SELECT i.id, i.name, i.unit_id, COALESCE(u.name, '') AS unit_name, COALESCE(u.abbreviation, '') AS unit_abbr,
		i.stock, i.total_used, i.total_wasted, i.min_stock, i.is_active, i.created_at, i.updated_at,
		EXISTS(SELECT 1 FROM ingredient_mutations im WHERE im.ingredient_id = i.id) AS has_history
		FROM ingredients i LEFT JOIN units u ON u.id = i.unit_id`
	countIngredientsBase       = `SELECT COUNT(*) FROM ingredients i WHERE 1=1`
	getAllIngredientsQuery     = ingredientSelect + ` WHERE 1=1`
	getIngredientsDefaultOrder = ` ORDER BY i.name ASC`
	getIngredientOptionsQuery  = ingredientSelect + ` WHERE i.is_active = 1 ORDER BY i.name ASC`
	getIngredientByIDQuery     = ingredientSelect + ` WHERE i.id = ? LIMIT 1`
	lockIngredientQuery        = ingredientSelect + ` WHERE i.id = ? LIMIT 1 FOR UPDATE OF i`
	ingredientNameExistsQuery  = `SELECT COUNT(*) FROM ingredients WHERE name = ? AND id <> ?`
	unitExistsQuery            = `SELECT COUNT(*) FROM units WHERE id = ? AND is_active = 1`
	createIngredientQuery      = `INSERT INTO ingredients (name, unit_id, min_stock, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	getLastInsertIDQuery       = `SELECT LAST_INSERT_ID()`
	updateIngredientQuery      = `UPDATE ingredients SET name = ?, unit_id = ?, min_stock = ?, updated_at = ? WHERE id = ?`
	toggleIngredientQuery      = `UPDATE ingredients SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteIngredientQuery      = `DELETE FROM ingredients WHERE id = ?`
	hasMutationsQuery          = `SELECT COUNT(*) FROM ingredient_mutations WHERE ingredient_id = ?`
	usedInRecipesQuery         = `SELECT COUNT(*) FROM menu_recipe_lines WHERE ingredient_id = ?`
	applyStockChangeQuery      = `UPDATE ingredients SET stock = stock + ?, total_used = total_used + ?, total_wasted = total_wasted + ?, updated_at = ? WHERE id = ?`

	createMutationQuery = `INSERT INTO ingredient_mutations (ingredient_id, type, qty_change, stock_before, stock_after, notes, reference_type, reference_id, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)`
	mutationSelect = `SELECT m.id, m.ingredient_id, i.name AS ingredient_name, COALESCE(u.abbreviation, '') AS unit_abbr, m.type, m.qty_change,
		m.stock_before, m.stock_after, COALESCE(m.notes, '') AS notes, COALESCE(m.reference_type, '') AS reference_type, m.reference_id,
		m.cancelled_at, m.cancelled_by, COALESCE(cu.full_name, '') AS cancelled_by_name, COALESCE(m.cancel_reason, '') AS cancel_reason,
		m.created_by, COALESCE(us.full_name, '') AS created_by_name, m.created_at
		FROM ingredient_mutations m
		INNER JOIN ingredients i ON i.id = m.ingredient_id
		LEFT JOIN units u ON u.id = i.unit_id
		LEFT JOIN users us ON us.id = m.created_by
		LEFT JOIN users cu ON cu.id = m.cancelled_by`
	countMutationsBase       = `SELECT COUNT(*) FROM ingredient_mutations m WHERE 1=1`
	getMutationsQuery        = mutationSelect + ` WHERE 1=1`
	getMutationsDefaultOrder = ` ORDER BY m.created_at DESC, m.id DESC`
	getMutationByIDQuery     = mutationSelect + ` WHERE m.id = ? LIMIT 1`
	lockMutationQuery        = `SELECT id, ingredient_id, type, qty_change, cancelled_at FROM ingredient_mutations WHERE id = ? LIMIT 1 FOR UPDATE`
	markCancelledQuery       = `UPDATE ingredient_mutations SET cancelled_at = ?, cancelled_by = ?, cancel_reason = ? WHERE id = ? AND cancelled_at IS NULL`
)

func (r *ingredientRepo) GetAll(req *dto.GetAllRequest) ([]*model.Ingredient, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		conditions += " AND i.name LIKE ?"
		args = append(args, "%"+req.Search+"%")
	}
	if req.IsActive != nil {
		conditions += " AND i.is_active = ?"
		args = append(args, *req.IsActive)
	}
	if req.LowStock {
		conditions += " AND i.min_stock IS NOT NULL AND i.stock <= i.min_stock"
	}

	var total int64
	if err := r.db.Raw(countIngredientsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
	allowedSortFields := map[string]string{
		"name":         "i.name",
		"stock":        "i.stock",
		"total_used":   "i.total_used",
		"total_wasted": "i.total_wasted",
		"is_active":    "i.is_active",
	}
	query := getAllIngredientsQuery + conditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, getIngredientsDefaultOrder) +
		" LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Ingredient
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *ingredientRepo) GetOptions() ([]*model.Ingredient, error) {
	var dataDB []*model.Ingredient
	if err := r.db.Raw(getIngredientOptionsQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *ingredientRepo) getOne(query string, id int) (*model.Ingredient, error) {
	var dataDB model.Ingredient
	if err := r.db.Raw(query, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *ingredientRepo) GetByID(id int) (*model.Ingredient, error) {
	return r.getOne(getIngredientByIDQuery, id)
}

func (r *ingredientRepo) GetByIDForUpdate(id int) (*model.Ingredient, error) {
	return r.getOne(lockIngredientQuery, id)
}

func (r *ingredientRepo) NameExists(name string, excludeID int) (bool, error) {
	var count int64
	if err := r.db.Raw(ingredientNameExistsQuery, name, excludeID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ingredientRepo) UnitExists(unitID int) (bool, error) {
	var count int64
	if err := r.db.Raw(unitExistsQuery, unitID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ingredientRepo) Create(req *dto.CreateRequest, userID int) (int64, error) {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(createIngredientQuery, req.Name, req.UnitID, req.MinStock, userID, now, now).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(getLastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ingredientRepo) Update(req *dto.UpdateRequest) error {
	return r.db.Exec(updateIngredientQuery, req.Name, req.UnitID, req.MinStock, time_helper.GetTimeNow(), req.ID).Error
}

func (r *ingredientRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleIngredientQuery, time_helper.GetTimeNow(), id).Error
}

func (r *ingredientRepo) Delete(id int) error {
	return r.db.Exec(deleteIngredientQuery, id).Error
}

func (r *ingredientRepo) HasMutations(id int) (bool, error) {
	var count int64
	if err := r.db.Raw(hasMutationsQuery, id).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ingredientRepo) UsedInRecipes(id int) (bool, error) {
	var count int64
	if err := r.db.Raw(usedInRecipesQuery, id).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ingredientRepo) ApplyStockChange(id int, stockDelta, usedDelta, wastedDelta float64) error {
	return r.db.Exec(applyStockChangeQuery, stockDelta, usedDelta, wastedDelta, time_helper.GetTimeNow(), id).Error
}

func (r *ingredientRepo) CreateMutation(m *model.IngredientMutation) (int64, error) {
	if err := r.db.Exec(createMutationQuery,
		m.IngredientID, m.Type, m.QtyChange, m.StockBefore, m.StockAfter, m.Notes,
		m.ReferenceType, m.ReferenceID, m.CreatedBy, m.CreatedAt,
	).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(getLastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ingredientRepo) GetMutations(req *dto.GetMutationsRequest) ([]*model.IngredientMutation, int64, error) {
	var args []any
	conditions := ""

	if req.IngredientID != nil {
		conditions += " AND m.ingredient_id = ?"
		args = append(args, *req.IngredientID)
	}
	if req.Type != "" {
		conditions += " AND m.type = ?"
		args = append(args, req.Type)
	}
	if req.StartDate != "" {
		conditions += " AND DATE(m.created_at) >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		conditions += " AND DATE(m.created_at) <= ?"
		args = append(args, req.EndDate)
	}

	var total int64
	if err := r.db.Raw(countMutationsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
	query := getMutationsQuery + conditions + getMutationsDefaultOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.IngredientMutation
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *ingredientRepo) GetMutationByID(id int) (*model.IngredientMutation, error) {
	var dataDB model.IngredientMutation
	if err := r.db.Raw(getMutationByIDQuery, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *ingredientRepo) GetMutationByIDForUpdate(id int) (*model.IngredientMutation, error) {
	var dataDB model.IngredientMutation
	if err := r.db.Raw(lockMutationQuery, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *ingredientRepo) MarkMutationCancelled(id, userID int, reason string) error {
	return r.db.Exec(markCancelledQuery, time_helper.GetTimeNow(), userID, reason, id).Error
}
