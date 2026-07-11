package repo

import (
	"encoding/json"
	"fmt"

	"pos_api/domain/sync/dto"
	"pos_api/domain/sync/model"
	request_helper "pos_api/helper/request"
)

const (
	getConflictsQuery      = `SELECT id, entity_type, entity_id, local_id, device_id, desktop_data, online_data, desktop_time, online_time, status, created_at FROM sync_conflicts WHERE 1=1`
	countConflictsBase     = `SELECT COUNT(*) FROM sync_conflicts WHERE 1=1`
	getConflictByIDQuery   = `SELECT id, entity_type, entity_id, desktop_data, online_data, desktop_time, online_time, status, resolved_by, resolution, resolved_at FROM sync_conflicts WHERE id = ?`
	resolveConflictQuery   = `UPDATE sync_conflicts SET status='resolved', resolved_by=?, resolved_action=?, resolved_at=NOW() WHERE id=?`
	createConflictQuery    = `INSERT INTO sync_conflicts (entity_type, entity_id, local_id, device_id, desktop_data, online_data, desktop_time, online_time) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	getQueueQuery          = `SELECT id, device_id, entity_type, entity_id, action, status, retry_count, created_at FROM sync_queue WHERE 1=1`
	countQueueBase         = `SELECT COUNT(*) FROM sync_queue WHERE 1=1`
	createQueueItemQuery   = `INSERT INTO sync_queue (device_id, entity_type, entity_id, action, payload, status) VALUES (?, ?, ?, ?, ?, 'pending')`
	updateQueueStatusQuery = `UPDATE sync_queue SET status=?, synced_at=CASE WHEN ? = 'synced' THEN NOW() ELSE NULL END, error_message=? WHERE id=?`
	insertHistoryQuery     = `INSERT INTO sync_history (device_id, device_type, total_items, synced_items, conflict_items, failed_items, pending_items, duration_ms, status, started_at, finished_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getHistoryQuery        = `SELECT id, device_id, device_type, total_items, synced_items, conflict_items, failed_items, pending_items, duration_ms, status, DATE_FORMAT(started_at,'%Y-%m-%dT%H:%i:%sZ'), DATE_FORMAT(finished_at,'%Y-%m-%dT%H:%i:%sZ') FROM sync_history WHERE 1=1`
	countHistoryBase       = `SELECT COUNT(*) FROM sync_history WHERE 1=1`
)

func (r *syncRepo) GetConflicts(filter *dto.ConflictFilter) ([]dto.ConflictResponse, int, error) {
	var args, countArgs []interface{}
	conditions := ""

	if filter.Status != "" {
		conditions += " AND status = ?"
		args = append(args, filter.Status)
		countArgs = append(countArgs, filter.Status)
	}

	var total int
	if err := r.db.Raw(countConflictsBase+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(filter.Page, filter.Limit, 20, 0)

	query := getConflictsQuery + conditions + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.ConflictResponse
	for rows.Next() {
		var item dto.ConflictResponse
		if err := rows.Scan(
			&item.ID, &item.EntityType, &item.EntityID,
			&item.LocalID, &item.DeviceID,
			&item.DesktopData, &item.OnlineData,
			&item.DesktopTime, &item.OnlineTime,
			&item.Status, &item.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.ConflictResponse{}
	}
	return items, total, nil
}

func (r *syncRepo) GetConflictByID(id int) (*model.SyncConflict, error) {
	row := r.db.Raw(getConflictByIDQuery, id).Row()
	var c model.SyncConflict
	if err := row.Scan(&c.ID, &c.EntityType, &c.EntityID, &c.DesktopData, &c.OnlineData, &c.DesktopTime, &c.OnlineTime, &c.Status, &c.ResolvedBy, &c.Resolution, &c.ResolvedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *syncRepo) CountPendingConflicts() (int, error) {
	var count int
	err := r.db.Raw(`SELECT COUNT(*) FROM sync_conflicts WHERE status = 'pending'`).Scan(&count).Error
	return count, err
}

func (r *syncRepo) ResolveConflict(id, userID int, action string) error {
	return r.db.Exec(resolveConflictQuery, userID, action, id).Error
}

func (r *syncRepo) MarkResolved(id, resolvedBy int, action string) error {
	return r.db.Exec(resolveConflictQuery, resolvedBy, action, id).Error
}

func (r *syncRepo) CreateConflict(deviceID string, item *dto.SyncItem, onlineData string) (int, error) {
	result := r.db.Exec(createConflictQuery,
		item.EntityType, item.ServerID, item.LocalID, deviceID,
		item.Payload, onlineData,
	)
	if result.Error != nil {
		return 0, result.Error
	}
	var id int
	r.db.Raw("SELECT LAST_INSERT_ID()").Scan(&id)
	return id, nil
}

func entityTable(entityType string) string {
	switch entityType {
	case "product":
		return "products"
	case "transaction":
		return "transactions"
	case "expense":
		return "expenses"
	}
	return ""
}

func (r *syncRepo) GetRawByEntityAndID(entityType string, serverID int) (string, error) {
	table := entityTable(entityType)
	if table == "" {
		return "{}", nil
	}

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ? LIMIT 1`, table)
	rows, err := r.db.Raw(query, serverID).Rows()
	if err != nil {
		return "{}", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "{}", nil
	}

	cols, err := rows.Columns()
	if err != nil {
		return "{}", err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}
	if err := rows.Scan(valuePtrs...); err != nil {
		return "{}", err
	}

	rowMap := make(map[string]interface{}, len(cols))
	for i, col := range cols {
		val := values[i]
		if b, ok := val.([]byte); ok {
			rowMap[col] = string(b)
		} else {
			rowMap[col] = val
		}
	}

	b, err := json.Marshal(rowMap)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func (r *syncRepo) GetEntitySnapshot(entityType string, serverID int) (*model.EntitySnapshot, error) {
	table := entityTable(entityType)
	if table == "" {
		return nil, nil
	}

	query := fmt.Sprintf(
		`SELECT DATE_FORMAT(updated_at, '%%Y-%%m-%%dT%%H:%%i:%%sZ') FROM %s WHERE id = ? LIMIT 1`,
		table,
	)

	var updatedAt string
	row := r.db.Raw(query, serverID).Row()
	if err := row.Scan(&updatedAt); err != nil {
		return nil, nil
	}

	rawJSON, _ := r.GetRawByEntityAndID(entityType, serverID)
	return &model.EntitySnapshot{
		UpdatedAt: updatedAt,
		Data:      rawJSON,
	}, nil
}

func (r *syncRepo) GetQueue(filter *dto.QueueFilter) ([]dto.QueueResponse, int, error) {
	var args, countArgs []interface{}
	conditions := ""

	if filter.DeviceID != "" {
		conditions += " AND device_id = ?"
		args = append(args, filter.DeviceID)
		countArgs = append(countArgs, filter.DeviceID)
	}
	if filter.Status != "" {
		conditions += " AND status = ?"
		args = append(args, filter.Status)
		countArgs = append(countArgs, filter.Status)
	}
	if filter.EntityType != "" {
		conditions += " AND entity_type = ?"
		args = append(args, filter.EntityType)
		countArgs = append(countArgs, filter.EntityType)
	}

	var total int
	if err := r.db.Raw(countQueueBase+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(filter.Page, filter.Limit, 20, 0)

	query := getQueueQuery + conditions + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.QueueResponse
	for rows.Next() {
		var item dto.QueueResponse
		if err := rows.Scan(&item.ID, &item.DeviceID, &item.EntityType, &item.EntityID, &item.Action, &item.Status, &item.RetryCount, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.QueueResponse{}
	}
	return items, total, nil
}

func (r *syncRepo) CreateQueueItem(deviceID string, item *dto.SyncItem) (int, error) {
	result := r.db.Exec(createQueueItemQuery, deviceID, item.EntityType, item.EntityID, item.Action, item.Payload)
	if result.Error != nil {
		return 0, result.Error
	}
	var id int
	r.db.Raw("SELECT LAST_INSERT_ID()").Scan(&id)
	return id, nil
}

func (r *syncRepo) UpdateQueueStatus(id int, status, errMsg string) error {
	return r.db.Exec(updateQueueStatusQuery, status, status, errMsg, id).Error
}

func (r *syncRepo) InsertHistory(h model.SyncHistory) error {
	return r.db.Exec(insertHistoryQuery,
		h.DeviceID, h.DeviceType, h.TotalItems, h.SyncedItems,
		h.ConflictItems, h.FailedItems, h.PendingItems, h.DurationMs, h.Status,
		h.StartedAt, h.FinishedAt,
	).Error
}

func (r *syncRepo) GetHistory(filter *dto.HistoryFilter) ([]dto.SyncHistoryResponse, int, error) {
	var args, countArgs []interface{}
	conditions := ""

	if filter.DeviceID != "" {
		conditions += " AND device_id = ?"
		args = append(args, filter.DeviceID)
		countArgs = append(countArgs, filter.DeviceID)
	}
	if filter.StartDate != "" {
		conditions += " AND DATE(started_at) >= ?"
		args = append(args, filter.StartDate)
		countArgs = append(countArgs, filter.StartDate)
	}
	if filter.EndDate != "" {
		conditions += " AND DATE(started_at) <= ?"
		args = append(args, filter.EndDate)
		countArgs = append(countArgs, filter.EndDate)
	}

	var total int
	if err := r.db.Raw(countHistoryBase+conditions, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(filter.Page, filter.Limit, 20, 0)

	query := getHistoryQuery + conditions + ` ORDER BY started_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []dto.SyncHistoryResponse
	for rows.Next() {
		var item dto.SyncHistoryResponse
		if err := rows.Scan(
			&item.ID, &item.DeviceID, &item.DeviceType,
			&item.TotalItems, &item.SyncedItems, &item.ConflictItems,
			&item.FailedItems, &item.PendingItems, &item.DurationMs, &item.Status,
			&item.StartedAt, &item.FinishedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []dto.SyncHistoryResponse{}
	}
	return items, total, nil
}
