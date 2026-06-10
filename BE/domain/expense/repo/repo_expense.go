package repo

import (
	dto "pos_api/domain/expense/dto"
	model "pos_api/domain/expense/model"
)

const (
	countExpensesBase   = `SELECT COUNT(*) FROM expenses e WHERE 1=1`
	getAllExpensesQuery  = `SELECT e.id, e.expense_date, e.category, e.description, e.amount, e.payment_method, e.user_id, COALESCE(u.full_name, '') as user_name, e.notes FROM expenses e LEFT JOIN users u ON e.user_id = u.id WHERE 1=1`
	getAllExpensesOrder  = ` ORDER BY e.expense_date DESC, e.id DESC`
	getExpenseByIDQuery = `SELECT e.id, e.expense_date, e.category, e.description, e.amount, e.payment_method, e.user_id, COALESCE(u.full_name, '') as user_name, e.notes FROM expenses e LEFT JOIN users u ON e.user_id = u.id WHERE e.id = ? LIMIT 1`
	createExpenseQuery  = `INSERT INTO expenses (expense_date, category, description, amount, payment_method, user_id, notes) VALUES (?, ?, ?, ?, ?, ?, ?)`
	getLastExpenseID    = `SELECT LAST_INSERT_ID()`
	updateExpenseQuery  = `UPDATE expenses SET expense_date=?, category=?, description=?, amount=?, payment_method=?, notes=?, updated_at=NOW() WHERE id=?`
	deleteExpenseQuery  = `DELETE FROM expenses WHERE id = ?`
)

func (r *expenseRepo) GetAll(req *dto.GetAllRequest) ([]*model.Expense, int64, error) {
	var args []any
	conditions := ""

	if req.StartDate != "" {
		conditions += " AND DATE(e.expense_date) >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		conditions += " AND DATE(e.expense_date) <= ?"
		args = append(args, req.EndDate)
	}
	if req.Category != "" {
		conditions += " AND e.category = ?"
		args = append(args, req.Category)
	}
	if req.UserID != nil {
		conditions += " AND e.user_id = ?"
		args = append(args, *req.UserID)
	}

	var total int64
	if err := r.db.Raw(countExpensesBase+conditions, args...).Scan(&total).Error; err != nil {
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

	query := getAllExpensesQuery + conditions + getAllExpensesOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Expense
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *expenseRepo) GetByID(id int) (*model.Expense, error) {
	var dataDB model.Expense
	if err := r.db.Raw(getExpenseByIDQuery, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *expenseRepo) Create(req *dto.CreateRequest, userID int) (int64, error) {
	if err := r.db.Exec(createExpenseQuery,
		req.ExpenseDate, req.Category, req.Description,
		req.Amount, req.PaymentMethod, userID, req.Notes,
	).Error; err != nil {
		return 0, err
	}
	var id int64
	if err := r.db.Raw(getLastExpenseID).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *expenseRepo) Update(req *dto.UpdateRequest) error {
	return r.db.Exec(updateExpenseQuery,
		req.ExpenseDate, req.Category, req.Description,
		req.Amount, req.PaymentMethod, req.Notes, req.ID,
	).Error
}

func (r *expenseRepo) Delete(req *dto.DeleteRequest) error {
	return r.db.Exec(deleteExpenseQuery, req.ID).Error
}

func (r *expenseRepo) UpdateFromSync(id int, data map[string]interface{}) error {
	allowed := []string{"expense_date", "category", "description", "amount", "payment_method", "notes"}
	updates := make(map[string]interface{}, len(allowed))
	for _, key := range allowed {
		if val, ok := data[key]; ok {
			updates[key] = val
		}
	}
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = "NOW()"
	return r.db.Table("expenses").Where("id = ?", id).Updates(updates).Error
}
