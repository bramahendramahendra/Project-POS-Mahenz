package repo

import (
	dto "pos_api/domain/supplier/dto"
	model "pos_api/domain/supplier/model"
)

const (
	countSuppliersQuery            = `SELECT COUNT(*) FROM suppliers WHERE 1=1`
	countSuppliersSearchQuery      = `SELECT COUNT(*) FROM suppliers WHERE (name LIKE ? OR supplier_code LIKE ? OR phone LIKE ?)`
	getAllSuppliersQuery           = `SELECT id, supplier_code, name, address, phone, email, contact_person, notes, is_active, created_at FROM suppliers WHERE 1=1`
	getAllSuppliersOrder           = ` ORDER BY name ASC`
	getAllSupplierOptionsQuery     = `SELECT id, supplier_code, name FROM suppliers WHERE is_active = 1 ORDER BY name`
	getSupplierByIDQuery           = `SELECT id, supplier_code, name, address, phone, email, contact_person, notes, is_active, created_at FROM suppliers WHERE id = ? LIMIT 1`
	getSupplierPurchasesQuery      = `SELECT id, purchase_code, purchase_date, total_amount, payment_status, remaining_amount FROM purchases WHERE supplier_id = ? ORDER BY purchase_date DESC LIMIT 10`
	getSupplierReturnsQuery        = `SELECT id, return_code, return_date, total_return_amount AS total_return, reason, status FROM supplier_returns WHERE supplier_id = ? ORDER BY return_date DESC LIMIT 10`
	generateSupplierCodeQuery      = `SELECT COUNT(*) FROM suppliers`
	checkSupplierCodeExistsQuery   = `SELECT id FROM suppliers WHERE supplier_code = ? LIMIT 1`
	checkSupplierNameExistsQuery   = `SELECT id FROM suppliers WHERE name = ? AND id != ? LIMIT 1`
	countPurchasesBySupplierQuery  = `SELECT COUNT(*) FROM purchases WHERE supplier_id = ?`
	countActiveDebtBySupplierQuery = `SELECT COUNT(*) FROM purchases WHERE supplier_id = ? AND payment_status != 'paid'`
	createSupplierQuery            = `INSERT INTO suppliers (supplier_code, name, address, phone, email, contact_person, notes) VALUES (?, ?, ?, ?, ?, ?, ?)`
	getLastSupplierInsertIDQuery   = `SELECT LAST_INSERT_ID()`
	updateSupplierQuery            = `UPDATE suppliers SET name=?, address=?, phone=?, email=?, contact_person=?, notes=?, updated_at=NOW() WHERE id=?`
	deleteSupplierQuery            = `DELETE FROM suppliers WHERE id = ?`
	toggleSupplierStatusQuery      = `UPDATE suppliers SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
)

func (r *supplierRepo) GetAll(req *dto.SupplierListRequest) ([]*model.Supplier, int64, error) {
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if req.Search != "" {
		search := "%" + req.Search + "%"
		countQuery := countSuppliersSearchQuery
		countArgs := []any{search, search, search}
		if req.IsActive != nil {
			countQuery += ` AND is_active = ?`
			countArgs = append(countArgs, *req.IsActive)
		}
		if err := r.db.Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	} else {
		countQuery := countSuppliersQuery
		var countArgs []any
		if req.IsActive != nil {
			countQuery += ` AND is_active = ?`
			countArgs = append(countArgs, *req.IsActive)
		}
		if err := r.db.Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	query := getAllSuppliersQuery
	var args []any
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query += ` AND (name LIKE ? OR supplier_code LIKE ? OR phone LIKE ?)`
		args = append(args, search, search, search)
	}
	if req.IsActive != nil {
		query += ` AND is_active = ?`
		args = append(args, *req.IsActive)
	}
	query += getAllSuppliersOrder
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Supplier
	err := r.db.Raw(query, args...).Scan(&dataDB).Error
	if err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *supplierRepo) GetOptions() ([]*dto.SupplierOptionResponse, error) {
	var dataDB []*dto.SupplierOptionResponse
	err := r.db.Raw(getAllSupplierOptionsQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *supplierRepo) GetByID(id int) (*model.Supplier, error) {
	var dataDB model.Supplier
	err := r.db.Raw(getSupplierByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *supplierRepo) Create(req *dto.CreateSupplierRequest, code string) (int64, error) {
	err := r.db.Exec(createSupplierQuery, code, req.Name, req.Address, req.Phone, req.Email, req.ContactPerson, req.Notes).Error
	if err != nil {
		return 0, err
	}

	var id int64
	errGet := r.db.Raw(getLastSupplierInsertIDQuery).Scan(&id).Error
	if errGet != nil {
		return 0, errGet
	}
	return id, nil
}

func (r *supplierRepo) Update(req *dto.UpdateSupplierRequest) error {
	return r.db.Exec(updateSupplierQuery, req.Name, req.Address, req.Phone, req.Email, req.ContactPerson, req.Notes, req.ID).Error
}

func (r *supplierRepo) Delete(req *dto.DeleteSupplierRequest) error {
	return r.db.Exec(deleteSupplierQuery, req.ID).Error
}

func (r *supplierRepo) ToggleStatus(req *dto.ToggleStatusSupplierRequest) error {
	return r.db.Exec(toggleSupplierStatusQuery, req.ID).Error
}

func (r *supplierRepo) GetPurchaseHistory(supplierID int) ([]dto.SupplierPurchaseItem, error) {
	var items []dto.SupplierPurchaseItem
	err := r.db.Raw(getSupplierPurchasesQuery, supplierID).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *supplierRepo) GetReturnHistory(supplierID int) ([]dto.SupplierReturnHistoryItem, error) {
	var items []dto.SupplierReturnHistoryItem
	if err := r.db.Raw(getSupplierReturnsQuery, supplierID).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *supplierRepo) CheckCodeExists(code string) (bool, error) {
	var id int
	if err := r.db.Raw(checkSupplierCodeExistsQuery, code).Scan(&id).Error; err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *supplierRepo) CheckNameExists(name string, excludeID int) (bool, error) {
	var id int
	if err := r.db.Raw(checkSupplierNameExistsQuery, name, excludeID).Scan(&id).Error; err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *supplierRepo) GetCount() (int, error) {
	var count int
	if err := r.db.Raw(generateSupplierCodeQuery).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *supplierRepo) CountPurchasesBySupplier(supplierID int) (int, error) {
	var count int
	if err := r.db.Raw(countPurchasesBySupplierQuery, supplierID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *supplierRepo) CountActiveDebtBySupplier(supplierID int) (int, error) {
	var count int
	if err := r.db.Raw(countActiveDebtBySupplierQuery, supplierID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
