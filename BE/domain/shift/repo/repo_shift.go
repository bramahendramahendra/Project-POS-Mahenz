package repo

import (
	dto "pos_api/domain/shift/dto"
	model "pos_api/domain/shift/model"
)

const (
	countShiftsQuery       = `SELECT COUNT(*) FROM shifts WHERE 1=1`
	getAllShiftsQuery      = `SELECT id, name, start_time, end_time, is_active FROM shifts WHERE 1=1`
	getAllShiftsOrder      = ` ORDER BY start_time`
	getActiveShiftsQuery   = `SELECT id, name, start_time, end_time FROM shifts WHERE is_active = 1 ORDER BY start_time`
	getShiftByIDQuery      = `SELECT id, name, start_time, end_time, is_active, created_at FROM shifts WHERE id = ? LIMIT 1`
	checkShiftUsedQuery    = `SELECT COUNT(*) FROM cash_drawer WHERE shift_id = ? AND status = 'open'`
	createShiftQuery       = `INSERT INTO shifts (name, start_time, end_time) VALUES (?, ?, ?)`
	getLastInsertIDQuery   = `SELECT LAST_INSERT_ID()`
	updateShiftQuery       = `UPDATE shifts SET name=?, start_time=?, end_time=?, updated_at=NOW() WHERE id=?`
	deleteShiftQuery       = `DELETE FROM shifts WHERE id = ?`
	toggleShiftStatusQuery = `UPDATE shifts SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
	getShiftSummaryQuery   = `SELECT cd.shift_id, s.name as shift_name, COUNT(t.id) as total_transactions, SUM(t.total_amount) as total_sales, SUM(CASE WHEN t.payment_method='cash' THEN t.total_amount ELSE 0 END) as total_cash, SUM(CASE WHEN t.payment_method!='cash' THEN t.total_amount ELSE 0 END) as total_non_cash FROM cash_drawer cd LEFT JOIN shifts s ON cd.shift_id = s.id LEFT JOIN transactions t ON DATE(t.transaction_date) = DATE(cd.open_time) WHERE 1=1 GROUP BY cd.shift_id, s.name`
)

func (r *shiftRepo) GetAll(req *dto.GetAllRequest) ([]*model.Shift, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND name LIKE ?`
		args = append(args, search)
	}

	var total int64
	if err := r.db.Raw(countShiftsQuery+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := getAllShiftsQuery + conditions + getAllShiftsOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Shift
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *shiftRepo) GetOptions() ([]*model.Shift, error) {
	var dataDB []*model.Shift
	if err := r.db.Raw(getActiveShiftsQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *shiftRepo) GetByID(id int) (*model.Shift, error) {
	var dataDB model.Shift
	if err := r.db.Raw(getShiftByIDQuery, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *shiftRepo) CountOpenCashDrawer(shiftID int) (int, error) {
	var count int
	if err := r.db.Raw(checkShiftUsedQuery, shiftID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *shiftRepo) Create(req *dto.CreateRequest) (int64, error) {
	if err := r.db.Exec(createShiftQuery, req.Name, req.StartTime, req.EndTime).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(getLastInsertIDQuery).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *shiftRepo) Update(req *dto.UpdateRequest) error {
	return r.db.Exec(updateShiftQuery, req.Name, req.StartTime, req.EndTime, req.ID).Error
}

func (r *shiftRepo) Delete(req *dto.DeleteRequest) error {
	return r.db.Exec(deleteShiftQuery, req.ID).Error
}

func (r *shiftRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	return r.db.Exec(toggleShiftStatusQuery, req.ID).Error
}

func (r *shiftRepo) GetSummary() ([]*dto.ShiftSummaryResponse, error) {
	var dataDB []*dto.ShiftSummaryResponse
	if err := r.db.Raw(getShiftSummaryQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}
