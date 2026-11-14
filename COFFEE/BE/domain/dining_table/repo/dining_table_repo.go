package repo

import (
	dto "coffee_api/domain/dining_table/dto"
	model "coffee_api/domain/dining_table/model"
	time_helper "coffee_api/helper/time"
)

// usage_count = jumlah pesanan yang memakai meja (sejak 4.7a); meja yang sudah dipakai tidak bisa dihapus.
const (
	tableSelect = `SELECT c.id, c.name, c.sort_order, c.is_active,
		(SELECT COUNT(*) FROM transactions t WHERE t.table_id = c.id) AS usage_count
		FROM dining_tables c`
	tableOrder       = ` ORDER BY c.sort_order ASC, c.name ASC`
	getTableOptions  = `SELECT c.id, c.name FROM dining_tables c WHERE c.is_active = 1` + tableOrder
	getTableByID     = tableSelect + ` WHERE c.id = ? LIMIT 1`
	tableNameExists  = `SELECT COUNT(*) FROM dining_tables WHERE name = ? AND id <> ?`
	createTableQuery = `INSERT INTO dining_tables (name, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?)`
	updateTableQuery = `UPDATE dining_tables SET name = ?, sort_order = ?, updated_at = ? WHERE id = ?`
	toggleTableQuery = `UPDATE dining_tables SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	deleteTableQuery = `DELETE FROM dining_tables WHERE id = ?`
)

func (r *diningTableRepo) GetAll(req *dto.GetAllRequest) ([]*model.DiningTable, error) {
	query := tableSelect + ` WHERE 1=1`
	var args []any
	if req.Search != "" {
		query += ` AND c.name LIKE ?`
		args = append(args, "%"+req.Search+"%")
	}
	if req.IsActive != nil {
		query += ` AND c.is_active = ?`
		args = append(args, *req.IsActive)
	}
	var data []*model.DiningTable
	if err := r.db.Raw(query+tableOrder, args...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *diningTableRepo) GetOptions() ([]*model.DiningTable, error) {
	var data []*model.DiningTable
	if err := r.db.Raw(getTableOptions).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *diningTableRepo) GetByID(id int) (*model.DiningTable, error) {
	var data model.DiningTable
	if err := r.db.Raw(getTableByID, id).Scan(&data).Error; err != nil {
		return nil, err
	}
	if data.ID == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *diningTableRepo) NameExists(name string, excludeID int) (bool, error) {
	var count int64
	if err := r.db.Raw(tableNameExists, name, excludeID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *diningTableRepo) Create(req *dto.SaveRequest) error {
	now := time_helper.GetTimeNow()
	return r.db.Exec(createTableQuery, req.Name, req.SortOrder, now, now).Error
}

func (r *diningTableRepo) Update(req *dto.SaveRequest) error {
	return r.db.Exec(updateTableQuery, req.Name, req.SortOrder, time_helper.GetTimeNow(), req.ID).Error
}

func (r *diningTableRepo) ToggleStatus(id int) error {
	return r.db.Exec(toggleTableQuery, time_helper.GetTimeNow(), id).Error
}

func (r *diningTableRepo) Delete(id int) error {
	return r.db.Exec(deleteTableQuery, id).Error
}
