package repo

import (
	"time"

	dto "pos_api/domain/receivable/dto"
	model "pos_api/domain/receivable/model"
)

const (
	getAllReceivablesQuery    = `SELECT r.id, t.transaction_code, c.name as customer_name, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE 1=1`
	countReceivablesBase     = `SELECT COUNT(*) FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE 1=1`
	getReceivableByIDQuery   = `SELECT r.id, r.transaction_id, r.customer_id, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date, r.created_at, r.updated_at FROM receivables r WHERE r.id = ?`
	getReceivableDetailQuery = `SELECT r.id, t.transaction_code, c.name as customer_name, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE r.id = ?`
	getReceivableSummaryQuery = `SELECT r.customer_id, c.name as customer_name, SUM(r.total_amount) as total_receivable, SUM(r.paid_amount) as total_paid, SUM(r.remaining_amount) as total_remaining, COUNT(r.id) as count FROM receivables r LEFT JOIN customers c ON r.customer_id = c.id WHERE r.status != 'paid' GROUP BY r.customer_id, c.name ORDER BY total_remaining DESC`
	createPaymentQuery       = `INSERT INTO receivable_payments (receivable_id, payment_date, amount, payment_method, notes, user_id) VALUES (?, ?, ?, ?, ?, ?)`
	updateReceivableQuery    = `UPDATE receivables SET paid_amount = paid_amount + ?, remaining_amount = remaining_amount - ?, status = CASE WHEN remaining_amount - ? <= 0 THEN 'paid' WHEN paid_amount + ? > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
	getPaymentsQuery         = `SELECT rp.id, rp.payment_date, rp.amount, rp.payment_method, u.full_name as user_name, rp.notes FROM receivable_payments rp LEFT JOIN users u ON rp.user_id = u.id WHERE rp.receivable_id = ? ORDER BY rp.payment_date DESC`
)

func (r *receivableRepo) GetAll(req *dto.GetAllRequest) ([]*dto.ReceivableResponse, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		conditions += " AND (c.name LIKE ? OR t.transaction_code LIKE ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like)
	}
	if req.Status != "" {
		conditions += " AND r.status = ?"
		args = append(args, req.Status)
	}

	var total int64
	if err := r.db.Raw(countReceivablesBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := getAllReceivablesQuery + conditions + " ORDER BY r.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*dto.ReceivableResponse
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *receivableRepo) GetByID(id int) (*model.Receivable, error) {
	var rec model.Receivable
	if err := r.db.Raw(getReceivableByIDQuery, id).Scan(&rec).Error; err != nil {
		return nil, err
	}
	if rec.ID == 0 {
		return nil, nil
	}
	return &rec, nil
}

func (r *receivableRepo) GetDetailByID(id int) (*dto.ReceivableDetailResponse, error) {
	var item dto.ReceivableDetailResponse
	if err := r.db.Raw(getReceivableDetailQuery, id).Scan(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, nil
	}
	return &item, nil
}

func (r *receivableRepo) GetSummary() ([]*dto.ReceivableSummaryItem, error) {
	var dataDB []*dto.ReceivableSummaryItem
	if err := r.db.Raw(getReceivableSummaryQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *receivableRepo) GetPayments(receivableID int) ([]*dto.PaymentResponse, error) {
	var dataDB []*dto.PaymentResponse
	if err := r.db.Raw(getPaymentsQuery, receivableID).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *receivableRepo) CreatePayment(receivableID int, req *dto.PayRequest, userID int) error {
	return r.db.Exec(createPaymentQuery, receivableID, time.Now(), req.Amount, req.PaymentMethod, req.Notes, userID).Error
}

func (r *receivableRepo) UpdateAfterPayment(receivableID int, amount float64) error {
	return r.db.Exec(updateReceivableQuery, amount, amount, amount, amount, receivableID).Error
}
