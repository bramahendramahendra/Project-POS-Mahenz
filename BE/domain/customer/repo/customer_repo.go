package repo

import (
	request_helper "pos_api/helper/request"
	dto "pos_api/domain/customer/dto"
	model "pos_api/domain/customer/model"
)

const (
	countCustomersBase         = `SELECT COUNT(*) FROM customers WHERE 1=1`
	getAllCustomersQuery       = `SELECT id, customer_code, name, phone, address, credit_limit, notes, is_active, created_at FROM customers WHERE 1=1`
	getActiveCustomerListQuery = `SELECT id, name, customer_code, credit_limit FROM customers WHERE is_active = 1 ORDER BY name`
	getCustomerByIDQuery       = `SELECT id, customer_code, name, phone, address, credit_limit, notes, is_active, created_at FROM customers WHERE id = ? LIMIT 1`
	checkCustomerHasReceivable = `SELECT COUNT(*) FROM receivables WHERE customer_id = ? AND status != 'paid'`
	generateCustomerCodeQuery  = `SELECT COUNT(*) FROM customers`
	createCustomerQuery        = `INSERT INTO customers (customer_code, name, phone, address, credit_limit, notes) VALUES (?, ?, ?, ?, ?, ?)`
	updateCustomerQuery        = `UPDATE customers SET name=?, phone=?, address=?, credit_limit=?, notes=?, updated_at=NOW() WHERE id=?`
	deleteCustomerQuery        = `DELETE FROM customers WHERE id = ?`
	toggleCustomerStatusQuery  = `UPDATE customers SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
	getLastCustomerInsertID    = `SELECT LAST_INSERT_ID()`
)

func (r *customerRepo) GetAll(req *dto.GetAllRequest) ([]*model.Customer, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += " AND (name LIKE ? OR customer_code LIKE ? OR phone LIKE ?)"
		args = append(args, search, search, search)
	}
	if req.IsActive != nil {
		conditions += " AND is_active = ?"
		args = append(args, *req.IsActive)
	}

	var total int64
	if err := r.db.Raw(countCustomersBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

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

	query := getAllCustomersQuery + conditions + " ORDER BY " + sortCol + " " + sortDir + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Customer
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *customerRepo) GetOptions() ([]*model.Customer, error) {
	var dataDB []*model.Customer
	err := r.db.Raw(getActiveCustomerListQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *customerRepo) GetByID(id int) (*model.Customer, error) {
	var dataDB model.Customer
	err := r.db.Raw(getCustomerByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *customerRepo) GetCount() (int, error) {
	var count int
	err := r.db.Raw(generateCustomerCodeQuery).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *customerRepo) CountActiveReceivables(customerID int) (int, error) {
	var count int
	err := r.db.Raw(checkCustomerHasReceivable, customerID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *customerRepo) Create(req *dto.CreateRequest, code string) (int64, error) {
	err := r.db.Exec(createCustomerQuery, code, req.Name, req.Phone, req.Address, req.CreditLimit, req.Notes).Error
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.Raw(getLastCustomerInsertID).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *customerRepo) Update(req *dto.UpdateRequest) error {
	err := r.db.Exec(updateCustomerQuery, req.Name, req.Phone, req.Address, req.CreditLimit, req.Notes, req.ID).Error
	return err
}

func (r *customerRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteCustomerQuery, req.ID).Error
	return err
}

func (r *customerRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleCustomerStatusQuery, req.ID).Error
	return err
}
