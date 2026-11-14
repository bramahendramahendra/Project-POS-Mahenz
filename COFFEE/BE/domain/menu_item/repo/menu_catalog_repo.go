package repo

import (
	model "coffee_api/domain/menu_item/model"
	product_repo "coffee_api/domain/product/repo"
)

const (
	catalogCategoriesQuery = `SELECT c.id, c.name FROM menu_categories c WHERE c.is_active = 1 ORDER BY c.sort_order, c.name`
	catalogItemsQuery      = `SELECT mi.id, mi.name, mi.category_id, mi.stock_mode, mi.product_id, mi.has_variants, mi.is_sold_out
		FROM menu_items mi JOIN menu_categories c ON c.id = mi.category_id
		WHERE mi.is_active = 1 AND c.is_active = 1
		ORDER BY c.sort_order, c.name, mi.name`
	// Urutan kombinasi sama dengan form menu: ukuran (urutan daftar ukuran), lalu panas sebelum dingin.
	catalogVariantsQuery = `SELECT v.id, v.menu_item_id, v.size_id, COALESCE(s.name, '') AS size_name, v.temperature, v.price
		FROM menu_item_variants v LEFT JOIN menu_sizes s ON s.id = v.size_id
		WHERE v.is_active = 1
		ORDER BY v.menu_item_id, COALESCE(s.sort_order, 0), s.name, FIELD(COALESCE(v.temperature, ''), '', 'hot', 'iced'), v.id`
	catalogGroupsQuery = `SELECT mg.menu_item_id, g.id, g.name, g.selection_type, g.is_required
		FROM menu_item_modifier_groups mg JOIN modifier_groups g ON g.id = mg.group_id
		WHERE g.is_active = 1
		ORDER BY mg.menu_item_id, mg.sort_order, g.name`
	catalogOptionsQuery = `SELECT o.id, o.group_id, o.name, o.price, o.is_default
		FROM modifier_options o WHERE o.is_active = 1 ORDER BY o.group_id, o.sort_order, o.id`
)

// CatalogData memuat seluruh bahan katalog kasir dalam lima query (tanpa query per menu).
func (r *menuItemRepo) CatalogData() (cats []model.CatalogCategory, items []model.CatalogItem, variants []model.CatalogVariant, groups []model.CatalogGroup, options []model.CatalogOption, err error) {
	if err = r.db.Raw(catalogCategoriesQuery).Scan(&cats).Error; err != nil {
		return
	}
	if err = r.db.Raw(catalogItemsQuery).Scan(&items).Error; err != nil {
		return
	}
	if err = r.db.Raw(catalogVariantsQuery).Scan(&variants).Error; err != nil {
		return
	}
	if err = r.db.Raw(catalogGroupsQuery).Scan(&groups).Error; err != nil {
		return
	}
	err = r.db.Raw(catalogOptionsQuery).Scan(&options).Error
	return
}

// AvailableStock: sisa stok yang bisa dijual per produk (satuan dasar, dikurangi yang ditahan),
// dihitung dengan cara yang sama seperti daftar produk.
func (r *menuItemRepo) AvailableStock(productIDs []int) (map[int]float64, error) {
	result := make(map[int]float64, len(productIDs))
	if len(productIDs) == 0 {
		return result, nil
	}
	minStock := make(map[int]float64, len(productIDs))
	for _, id := range productIDs {
		minStock[id] = 0
	}
	summaries, err := product_repo.BuildStockSummaries(r.db, minStock)
	if err != nil {
		return nil, err
	}
	for _, id := range productIDs {
		if s, ok := summaries[id]; ok {
			result[id] = s.AnchorStock - s.AnchorReserved
		}
	}
	return result, nil
}
