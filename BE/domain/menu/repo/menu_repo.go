package repo

import (
	dto "pos_api/domain/menu/dto"
	model "pos_api/domain/menu/model"
	request_helper "pos_api/helper/request"
)

const (
	countMenusQuery         = `SELECT COUNT(*) FROM menus WHERE 1=1`
	getAllMenusQuery        = `SELECT id, parent_id, key_name, label, icon, path, order_index, is_active, created_at, updated_at FROM menus WHERE 1=1`
	getMenuRootOptionsQuery = `SELECT id, label FROM menus WHERE parent_id IS NULL AND is_active = 1 ORDER BY order_index ASC, id ASC`
	getMenuByIDQuery        = `SELECT id, parent_id, key_name, label, icon, path, order_index, is_active, created_at, updated_at FROM menus WHERE id = ? LIMIT 1`
	getMenuByKeyNameQuery   = `SELECT id, parent_id, key_name, label, icon, path, order_index, is_active, created_at, updated_at FROM menus WHERE key_name = ? LIMIT 1`
	createMenuQuery         = `INSERT INTO menus (parent_id, key_name, label, icon, path, order_index) VALUES (?, ?, ?, ?, ?, ?)`
	updateMenuQuery         = `UPDATE menus SET parent_id = ?, label = ?, icon = ?, path = ?, order_index = ?, updated_at = NOW() WHERE id = ?`
	deleteMenuQuery         = `DELETE FROM menus WHERE id = ?`
	updateOrderQuery        = `UPDATE menus SET order_index = ?, updated_at = NOW() WHERE id = ?`
	// LEFT JOIN dengan semua menu aktif (bukan hanya yang punya can_view=1) supaya kategori induk
	// tanpa izin langsung tetap ikut terbawa sebagai node struktural saat ada anak yang punya akses.
	// Filtering item yang benar-benar tidak boleh dilihat dilakukan belakangan oleh pruneTree().
	getMyMenusQuery = `
		SELECT
			m.id,
			m.parent_id,
			m.key_name,
			m.label,
			m.icon,
			m.path,
			m.order_index,
			COALESCE(rma.can_view,   0),
			COALESCE(rma.can_create, 0),
			COALESCE(rma.can_edit,   0),
			COALESCE(rma.can_delete, 0)
		FROM menus m
		LEFT JOIN role_menu_access rma ON rma.menu_id = m.id
			AND rma.role_id = (SELECT id FROM roles WHERE name = ? AND is_active = 1 LIMIT 1)
		WHERE m.is_active = 1
		ORDER BY m.order_index ASC, m.id ASC
		`
)

func (r *menuRepo) GetAll(req *dto.GetAllRequest) ([]*model.Menu, int64, error) {
	var args, countArgs []any
	conditions := ""

	if req.Search != "" {
		like := "%" + req.Search + "%"
		conditions += ` AND (key_name LIKE ? OR label LIKE ?)`
		args = append(args, like, like)
		countArgs = append(countArgs, like, like)
	}
	if req.IsActive != nil {
		conditions += ` AND is_active = ?`
		args = append(args, *req.IsActive)
		countArgs = append(countArgs, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countMenusQuery+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	query := getAllMenusQuery + conditions + ` ORDER BY order_index ASC, id ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var dataDB []*model.Menu
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *menuRepo) GetRootOptions() ([]*dto.MenuOptionResponse, error) {
	var options []*dto.MenuOptionResponse
	if err := r.db.Raw(getMenuRootOptionsQuery).Scan(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *menuRepo) GetByID(id int) (*model.Menu, error) {
	var m model.Menu
	if err := r.db.Raw(getMenuByIDQuery, id).Scan(&m).Error; err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *menuRepo) GetByKeyName(keyName string) (*model.Menu, error) {
	var m model.Menu
	if err := r.db.Raw(getMenuByKeyNameQuery, keyName).Scan(&m).Error; err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *menuRepo) GetMyMenus(roleName string) ([]*dto.MyMenuItem, error) {
	type rawRow struct {
		ID         int
		ParentID   *int
		KeyName    string
		Label      string
		Icon       *string
		Path       *string
		OrderIndex int
		CanView    bool
		CanCreate  bool
		CanEdit    bool
		CanDelete  bool
	}

	rows, err := r.db.Raw(getMyMenusQuery, roleName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*dto.MyMenuItem
	for rows.Next() {
		var row rawRow
		if err := rows.Scan(&row.ID, &row.ParentID, &row.KeyName, &row.Label, &row.Icon,
			&row.Path, &row.OrderIndex, &row.CanView, &row.CanCreate, &row.CanEdit, &row.CanDelete); err != nil {
			return nil, err
		}
		items = append(items, &dto.MyMenuItem{
			KeyName:    row.KeyName,
			Label:      row.Label,
			Icon:       row.Icon,
			Path:       row.Path,
			OrderIndex: row.OrderIndex,
			Permission: dto.MenuPermission{
				CanView:   row.CanView,
				CanCreate: row.CanCreate,
				CanEdit:   row.CanEdit,
				CanDelete: row.CanDelete,
			},
			Children: []dto.MyMenuItem{},
		})
	}
	return items, nil
}

func (r *menuRepo) Create(req *dto.CreateRequest) (int64, error) {
	var icon, path *string
	if req.Icon != "" {
		icon = &req.Icon
	}
	if req.Path != "" {
		path = &req.Path
	}
	if err := r.db.Exec(createMenuQuery, req.ParentID, req.KeyName, req.Label, icon, path, req.OrderIndex).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(`SELECT LAST_INSERT_ID()`).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *menuRepo) Update(id int, req *dto.UpdateRequest) error {
	var icon, path *string
	if req.Icon != "" {
		icon = &req.Icon
	}
	if req.Path != "" {
		path = &req.Path
	}
	return r.db.Exec(updateMenuQuery, req.ParentID, req.Label, icon, path, req.OrderIndex, id).Error
}

func (r *menuRepo) Delete(id int) error {
	return r.db.Exec(deleteMenuQuery, id).Error
}

func (r *menuRepo) Reorder(items []dto.ReorderItem) error {
	for _, item := range items {
		if err := r.db.Exec(updateOrderQuery, item.OrderIndex, item.ID).Error; err != nil {
			return err
		}
	}
	return nil
}
