package repo

import (
	dto "pos_api/domain/product_unit/dto"
	model "pos_api/domain/product_unit/model"
)

const (
	getAllUnitsQuery       = `SELECT id, name, abbreviation, is_active, created_at FROM units ORDER BY name`
	getActiveUnitsQuery   = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE is_active = 1 ORDER BY name`
	getUnitByIDQuery      = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE id = ? LIMIT 1`
	checkUnitNameQuery    = `SELECT id FROM units WHERE name = ? AND id != ? LIMIT 1`
	checkUnitUsedQuery    = `SELECT COUNT(*) FROM product_units WHERE unit_id = ?`
	createUnitQuery       = `INSERT INTO units (name, abbreviation) VALUES (?, ?)`
	getLastInsertIDQuery  = `SELECT LAST_INSERT_ID()`
	updateUnitQuery       = `UPDATE units SET name = ?, abbreviation = ?, updated_at = NOW() WHERE id = ?`
	deleteUnitQuery       = `DELETE FROM units WHERE id = ?`
	toggleUnitStatusQuery = `UPDATE units SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *unitRepo) GetAll() ([]*model.Unit, error) {
	var units []*model.Unit
	err := r.db.Raw(getAllUnitsQuery).Scan(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (r *unitRepo) GetActive() ([]*model.Unit, error) {
	var units []*model.Unit
	err := r.db.Raw(getActiveUnitsQuery).Scan(&units).Error
	if err != nil {
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
