package repo

import (
	dto "coffee_api/domain/menu_item/dto"
	model "coffee_api/domain/menu_item/model"
	request_helper "coffee_api/helper/request"
	time_helper "coffee_api/helper/time"
)

const (
	// Harga min/maks dan jumlah varian dihitung dari kombinasi yang aktif.
	itemSelect = `SELECT mi.id, mi.category_id, COALESCE(c.name, '') AS category_name, mi.name, mi.stock_mode,
		mi.product_id, COALESCE(p.name, '') AS product_name,
		mi.has_variants, mi.is_active, mi.is_sold_out,
		COALESCE((SELECT MIN(v.price) FROM menu_item_variants v WHERE v.menu_item_id = mi.id AND v.is_active = 1), 0) AS price_min,
		COALESCE((SELECT MAX(v.price) FROM menu_item_variants v WHERE v.menu_item_id = mi.id AND v.is_active = 1), 0) AS price_max,
		(SELECT COUNT(*) FROM menu_item_variants v WHERE v.menu_item_id = mi.id AND v.is_active = 1) AS variant_count
		FROM menu_items mi LEFT JOIN menu_categories c ON c.id = mi.category_id LEFT JOIN products p ON p.id = mi.product_id`
	countItemsBase    = `SELECT COUNT(*) FROM menu_items mi WHERE 1=1`
	itemsDefaultOrder = ` ORDER BY c.sort_order ASC, c.name ASC, mi.name ASC`
	getItemByIDQuery  = itemSelect + ` WHERE mi.id = ? LIMIT 1`

	// Urutan varian: ukuran (urutan daftar ukuran), lalu panas sebelum dingin.
	getVariantsQuery = `SELECT v.id, v.menu_item_id, v.size_id, COALESCE(s.name, '') AS size_name, COALESCE(s.sort_order, 0) AS size_order,
		v.temperature, v.price, v.is_active
		FROM menu_item_variants v LEFT JOIN menu_sizes s ON s.id = v.size_id
		WHERE v.menu_item_id = ?
		ORDER BY COALESCE(s.sort_order, 0), s.name, FIELD(COALESCE(v.temperature, ''), '', 'hot', 'iced'), v.id`
	getRecipeLinesQuery = `SELECT l.id, l.variant_id, l.ingredient_id, i.name AS ingredient_name, l.amount, l.sort_order
		FROM menu_recipe_lines l INNER JOIN ingredients i ON i.id = l.ingredient_id
		WHERE l.variant_id IN ? ORDER BY l.variant_id, l.sort_order, l.id`
	getItemGroupsQuery = `SELECT mg.group_id, g.name AS group_name, mg.sort_order
		FROM menu_item_modifier_groups mg INNER JOIN modifier_groups g ON g.id = mg.group_id
		WHERE mg.menu_item_id = ? ORDER BY mg.sort_order, g.name`

	itemNameExistsQuery    = `SELECT COUNT(*) FROM menu_items WHERE name = ? AND id <> ?`
	categoryActiveQuery    = `SELECT COUNT(*) FROM menu_categories WHERE id = ? AND is_active = 1`
	countActiveSizesQuery  = `SELECT COUNT(*) FROM menu_sizes WHERE id IN ? AND is_active = 1`
	countIngredientsQuery  = `SELECT COUNT(*) FROM ingredients WHERE id IN ?`
	productActiveQuery     = `SELECT COUNT(*) FROM products WHERE id = ? AND is_active = 1`
	productLinkedQuery     = `SELECT COUNT(*) FROM menu_items WHERE product_id = ? AND id <> ?`
	countActiveGroupsQuery = `SELECT COUNT(*) FROM modifier_groups WHERE id IN ? AND is_active = 1`

	createItemQuery      = `INSERT INTO menu_items (category_id, name, stock_mode, product_id, has_variants, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	lastInsertIDQuery    = `SELECT LAST_INSERT_ID()`
	updateItemQuery      = `UPDATE menu_items SET category_id = ?, name = ?, stock_mode = ?, product_id = ?, has_variants = ?, updated_at = ? WHERE id = ?`
	insertVariantQuery   = `INSERT INTO menu_item_variants (menu_item_id, size_id, temperature, price, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	updateVariantQuery   = `UPDATE menu_item_variants SET price = ?, is_active = ?, updated_at = ? WHERE id = ?`
	deleteVariantsQuery  = `DELETE FROM menu_item_variants WHERE menu_item_id = ?`
	deleteVariantsNotIn  = `DELETE FROM menu_item_variants WHERE menu_item_id = ? AND id NOT IN ?`
	deleteRecipeQuery    = `DELETE FROM menu_recipe_lines WHERE variant_id = ?`
	insertRecipeQuery    = `INSERT INTO menu_recipe_lines (variant_id, ingredient_id, amount, sort_order) VALUES (?, ?, ?, ?)`
	deleteItemGroups     = `DELETE FROM menu_item_modifier_groups WHERE menu_item_id = ?`
	insertItemGroupQuery = `INSERT INTO menu_item_modifier_groups (menu_item_id, group_id, sort_order) VALUES (?, ?, ?)`
	toggleItemQuery      = `UPDATE menu_items SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	setSoldOutQuery      = `UPDATE menu_items SET is_sold_out = ?, updated_at = ? WHERE id = ?`
	deleteItemQuery      = `DELETE FROM menu_items WHERE id = ?`
)

func (r *menuItemRepo) GetAll(req *dto.GetAllRequest) ([]*model.MenuItem, int64, error) {
	var args []any
	conditions := ""
	if req.Search != "" {
		conditions += " AND mi.name LIKE ?"
		args = append(args, "%"+req.Search+"%")
	}
	if req.CategoryID != nil {
		conditions += " AND mi.category_id = ?"
		args = append(args, *req.CategoryID)
	}
	if req.StockMode != "" {
		conditions += " AND mi.stock_mode = ?"
		args = append(args, req.StockMode)
	}
	if req.IsActive != nil {
		conditions += " AND mi.is_active = ?"
		args = append(args, *req.IsActive)
	}
	if req.IsSoldOut != nil {
		conditions += " AND mi.is_sold_out = ?"
		args = append(args, *req.IsSoldOut)
	}

	var total int64
	if err := r.db.Raw(countItemsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
	allowedSortFields := map[string]string{
		"name":          "mi.name",
		"category_name": "c.name",
		"price_min":     "price_min",
		"is_active":     "mi.is_active",
	}
	query := itemSelect + " WHERE 1=1" + conditions +
		request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, itemsDefaultOrder) +
		" LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var data []*model.MenuItem
	if err := r.db.Raw(query, args...).Scan(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r *menuItemRepo) GetByID(id int) (*model.MenuItem, error) {
	var data model.MenuItem
	if err := r.db.Raw(getItemByIDQuery, id).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *menuItemRepo) GetVariants(itemID int) ([]*model.Variant, error) {
	var data []*model.Variant
	if err := r.db.Raw(getVariantsQuery, itemID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuItemRepo) GetRecipeLines(variantIDs []int) ([]*model.RecipeLine, error) {
	if len(variantIDs) == 0 {
		return nil, nil
	}
	var data []*model.RecipeLine
	if err := r.db.Raw(getRecipeLinesQuery, variantIDs).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuItemRepo) GetModifierGroups(itemID int) ([]*model.ItemModifierGroup, error) {
	var data []*model.ItemModifierGroup
	if err := r.db.Raw(getItemGroupsQuery, itemID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *menuItemRepo) count(query string, args ...any) (int64, error) {
	var n int64
	if err := r.db.Raw(query, args...).Scan(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

func (r *menuItemRepo) NameExists(name string, excludeID int) (bool, error) {
	n, err := r.count(itemNameExistsQuery, name, excludeID)
	return n > 0, err
}

func (r *menuItemRepo) CategoryIsActive(id int) (bool, error) {
	n, err := r.count(categoryActiveQuery, id)
	return n > 0, err
}

// ProductIsActive: produk barang jadi ada dan aktif.
func (r *menuItemRepo) ProductIsActive(id int) (bool, error) {
	n, err := r.count(productActiveQuery, id)
	return n > 0, err
}

// ProductLinked: produk sudah ditautkan ke menu lain (satu produk untuk satu menu, keputusan 91).
func (r *menuItemRepo) ProductLinked(productID int, excludeID int) (bool, error) {
	n, err := r.count(productLinkedQuery, productID, excludeID)
	return n > 0, err
}

func (r *menuItemRepo) CountActiveSizes(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return r.count(countActiveSizesQuery, ids)
}

func (r *menuItemRepo) CountIngredients(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return r.count(countIngredientsQuery, ids)
}

func (r *menuItemRepo) CountActiveModifierGroups(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return r.count(countActiveGroupsQuery, ids)
}

func (r *menuItemRepo) Create(req *dto.SaveRequest, userID int) (int64, error) {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(createItemQuery, req.CategoryID, req.Name, req.StockMode, req.ProductID, req.HasVariants, userID, now, now).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(lastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *menuItemRepo) Update(req *dto.SaveRequest) error {
	return r.db.Exec(updateItemQuery, req.CategoryID, req.Name, req.StockMode, req.ProductID, req.HasVariants, time_helper.GetTimeNow(), req.ID).Error
}

func (r *menuItemRepo) InsertVariant(itemID int, v *dto.VariantInput) (int64, error) {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(insertVariantQuery, itemID, v.SizeID, v.Temperature, v.Price, v.IsActive, now, now).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(lastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *menuItemRepo) UpdateVariant(variantID int, price float64, isActive bool) error {
	return r.db.Exec(updateVariantQuery, price, isActive, time_helper.GetTimeNow(), variantID).Error
}

func (r *menuItemRepo) DeleteVariantsNotIn(itemID int, keepIDs []int) error {
	if len(keepIDs) == 0 {
		return r.db.Exec(deleteVariantsQuery, itemID).Error
	}
	return r.db.Exec(deleteVariantsNotIn, itemID, keepIDs).Error
}

func (r *menuItemRepo) DeleteRecipeLines(variantID int) error {
	return r.db.Exec(deleteRecipeQuery, variantID).Error
}

func (r *menuItemRepo) InsertRecipeLine(variantID int, line *dto.RecipeInput, sortOrder int) error {
	return r.db.Exec(insertRecipeQuery, variantID, line.IngredientID, line.Amount, sortOrder).Error
}

func (r *menuItemRepo) ReplaceModifierGroups(itemID int, groupIDs []int) error {
	if err := r.db.Exec(deleteItemGroups, itemID).Error; err != nil {
		return err
	}
	for i, gid := range groupIDs {
		if err := r.db.Exec(insertItemGroupQuery, itemID, gid, i).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *menuItemRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleItemQuery, time_helper.GetTimeNow(), id).Error
}

func (r *menuItemRepo) SetSoldOut(id int, soldOut bool) error {
	return r.db.Exec(setSoldOutQuery, soldOut, time_helper.GetTimeNow(), id).Error
}

func (r *menuItemRepo) Delete(id int) error {
	return r.db.Exec(deleteItemQuery, id).Error
}
