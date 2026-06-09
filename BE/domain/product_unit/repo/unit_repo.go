package repo

import (
	dto "pos_api/domain/product_unit/dto"
	model "pos_api/domain/product_unit/model"
)

const (
	countUnitsQuery = `SELECT COUNT(*) FROM units WHERE 1=1`
	getAllUnitsQuery               = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE 1=1`
	getAllUnitsOrder               = ` ORDER BY name`
	getAllUnitOptionsQuery         = `SELECT id, name, abbreviation FROM units WHERE is_active = 1 ORDER BY name`
	getUnitByIDQuery               = `SELECT id, name, abbreviation, is_active, created_at FROM units WHERE id = ? LIMIT 1`
	checkUnitNameQuery             = `SELECT id FROM units WHERE name = ? AND id != ? LIMIT 1`
	checkUnitUsedQuery             = `SELECT COUNT(*) FROM product_units WHERE unit_id = ?`
	checkActiveProductsByUnitQuery = `SELECT COUNT(*) FROM product_units pu JOIN products p ON p.id = pu.product_id WHERE pu.unit_id = ? AND p.is_active = 1`
	createUnitQuery                = `INSERT INTO units (name, abbreviation) VALUES (?, ?)`
	getLastInsertIDQuery           = `SELECT LAST_INSERT_ID()`
	updateUnitQuery                = `UPDATE units SET name = ?, abbreviation = ?, updated_at = NOW() WHERE id = ?`
	deleteUnitQuery                = `DELETE FROM units WHERE id = ?`
	toggleUnitStatusQuery          = `UPDATE units SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *unitRepo) GetAll(req *dto.GetAllRequest) ([]*model.Unit, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND name LIKE ?`
		args = append(args, search)
	}

	var total int64
	if err := r.db.Raw(countUnitsQuery+conditions, args...).Scan(&total).Error; err != nil {
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

	query := getAllUnitsQuery + conditions + getAllUnitsOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Unit
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *unitRepo) GetOptions() ([]*model.UnitOption, error) {
	var dataDB []*model.UnitOption
	err := r.db.Raw(getAllUnitOptionsQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *unitRepo) GetByID(id int) (*model.Unit, error) {
	var dataDB model.Unit
	err := r.db.Raw(getUnitByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *unitRepo) Create(req *dto.CreateRequest) (int64, error) {
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

func (r *unitRepo) Update(req *dto.UpdateRequest) error {
	err := r.db.Exec(updateUnitQuery, req.Name, req.Abbreviation, req.ID).Error
	return err
}

func (r *unitRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteUnitQuery, req.ID).Error
	return err
}

func (r *unitRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleUnitStatusQuery, req.ID).Error
	return err
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

func (r *unitRepo) CountActiveProductsByUnit(unitID int) (int, error) {
	var count int
	err := r.db.Raw(checkActiveProductsByUnitQuery, unitID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
