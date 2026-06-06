package repo

import (
	"fmt"

	dto "pos_api/domain/product_unit/dto"
	model "pos_api/domain/product_unit/model"
)

const (
	countUnitsQuery        = `SELECT COUNT(*) FROM units WHERE 1=1`
	countUnitsSearchQuery  = `SELECT COUNT(*) FROM units WHERE name LIKE ?`
	getAllUnitsQuery       = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE 1=1`
	getAllUnitsOrder       = ` ORDER BY name`
	getAllUnitOptionsQuery = `SELECT id, name, abbreviation FROM units WHERE is_active = 1 ORDER BY name`
	getUnitByIDQuery       = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE id = ? LIMIT 1`
	checkUnitNameQuery     = `SELECT id FROM units WHERE name = ? AND id != ? LIMIT 1`
	checkUnitUsedQuery     = `SELECT COUNT(*) FROM product_units WHERE unit_id = ?`
	createUnitQuery        = `INSERT INTO units (name, abbreviation) VALUES (?, ?)`
	getLastInsertIDQuery   = `SELECT LAST_INSERT_ID()`
	updateUnitQuery        = `UPDATE units SET name = ?, abbreviation = ?, updated_at = NOW() WHERE id = ?`
	deleteUnitQuery        = `DELETE FROM units WHERE id = ?`
	toggleUnitStatusQuery  = `UPDATE units SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *unitRepo) GetAll(req *dto.UnitListRequest) ([]*model.Unit, int64, error) {
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if req.Search != "" {
		if err := r.db.Raw(countUnitsSearchQuery, "%"+req.Search+"%").Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.Raw(countUnitsQuery).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	query := getAllUnitsQuery
	var args []any
	if req.Search != "" {
		query += ` AND name LIKE ?`
		args = append(args, "%"+req.Search+"%")
	}
	query += getAllUnitsOrder
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var units []*model.Unit
	if err := r.db.Raw(query, args...).Scan(&units).Error; err != nil {
		return nil, 0, err
	}
	return units, total, nil
}

func (r *unitRepo) GetOptions() ([]*dto.UnitActiveResponse, error) {
	var units []*dto.UnitActiveResponse
	if err := r.db.Raw(getAllUnitOptionsQuery).Scan(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

func (r *unitRepo) GetByID(id int) (*model.Unit, error) {
	var unit model.Unit
	err := r.db.Raw(getUnitByIDQuery, id).Scan(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepo) Create(req *dto.CreateUnitRequest) (int64, error) {
	err := r.db.Exec(createUnitQuery, req.Name, req.Abbreviation).Error
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

func (r *unitRepo) Update(req *dto.UpdateUnitRequest) error {
	return r.db.Exec(updateUnitQuery, req.Name, req.Abbreviation, req.ID).Error
}

func (r *unitRepo) Delete(req *dto.DeleteUnitRequest) error {
	return r.db.Exec(deleteUnitQuery, req.ID).Error
}

func (r *unitRepo) ToggleStatus(req *dto.ToggleStatusUnitRequest) error {
	return r.db.Exec(toggleUnitStatusQuery, req.ID).Error
}

func (r *unitRepo) CheckNameExists(name string, excludeID int) (bool, error) {
	var id int
	err := r.db.Raw(checkUnitNameQuery, name, excludeID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *unitRepo) CountProductUnitsByUnit(unitID int) (int, error) {
	var count int
	err := r.db.Raw(checkUnitUsedQuery, unitID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
