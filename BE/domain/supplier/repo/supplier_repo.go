package repo

import (
	dto "pos_api/domain/supplier/dto"
	model "pos_api/domain/supplier/model"
)

const (
	countSuppliersQuery = `SELECT COUNT(*) FROM suppliers WHERE 1=1`
	getAllSuppliersQuery = `SELECT id, supplier_code, name, address, phone, email, contact_person, notes, is_active, created_at FROM suppliers WHERE 1=1`
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

func (r *supplierRepo) GetAll(req *dto.GetAllRequest) ([]*model.Supplier, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND (name LIKE ? OR supplier_code LIKE ? OR phone LIKE ?)`
		args = append(args, search, search, search)
	}
	if req.IsActive != nil {
		conditions += ` AND is_active = ?`
		args = append(args, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countSuppliersQuery+conditions, args...).Scan(&total).Error; err != nil {
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

	allowedSortColumns := map[string]string{
		"name":      "name",
		"is_active": "is_active",
	}
	sortCol := "name"
	if col, ok := allowedSortColumns[req.SortBy]; ok {
		sortCol = col
	}
	sortDir := "ASC"
	if req.SortOrder == "desc" {
		sortDir = "DESC"
	}

	query := getAllSuppliersQuery + conditions + " ORDER BY " + sortCol + " " + sortDir + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Supplier
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *supplierRepo) GetOptions() ([]*model.SupplierOption, error) {
	var dataDB []*model.SupplierOption
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

func (r *supplierRepo) Create(req *dto.CreateRequest, code string) (int64, error) {
	err := r.db.Exec(createSupplierQuery, code, req.Name, req.Address, req.Phone, req.Email, req.ContactPerson, req.Notes).Error
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.Raw(getLastSupplierInsertIDQuery).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *supplierRepo) Update(req *dto.UpdateRequest) error {
	err := r.db.Exec(updateSupplierQuery, req.Name, req.Address, req.Phone, req.Email, req.ContactPerson, req.Notes, req.ID).Error
	return err
}

func (r *supplierRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteSupplierQuery, req.ID).Error
	return err
}

func (r *supplierRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleSupplierStatusQuery, req.ID).Error
	return err
}

func (r *supplierRepo) GetPurchaseHistory(supplierID int) ([]*model.SupplierPurchase, error) {
	var dataDB []*model.SupplierPurchase
	err := r.db.Raw(getSupplierPurchasesQuery, supplierID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *supplierRepo) GetReturnHistory(supplierID int) ([]*model.SupplierReturn, error) {
	var dataDB []*model.SupplierReturn
	err := r.db.Raw(getSupplierReturnsQuery, supplierID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *supplierRepo) CheckCodeExists(code string) (bool, error) {
	var id int
	err := r.db.Raw(checkSupplierCodeExistsQuery, code).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *supplierRepo) CheckNameExists(name string, excludeID int) (bool, error) {
	var id int
	err := r.db.Raw(checkSupplierNameExistsQuery, name, excludeID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *supplierRepo) GetCount() (int, error) {
	var count int
	err := r.db.Raw(generateSupplierCodeQuery).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *supplierRepo) CountPurchasesBySupplier(supplierID int) (int, error) {
	var count int
	err := r.db.Raw(countPurchasesBySupplierQuery, supplierID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *supplierRepo) CountActiveDebtBySupplier(supplierID int) (int, error) {
	var count int
	err := r.db.Raw(countActiveDebtBySupplierQuery, supplierID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
