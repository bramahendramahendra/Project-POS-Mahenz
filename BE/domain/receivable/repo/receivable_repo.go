package repo

import (
	"time"

	dto "pos_api/domain/receivable/dto"
	model "pos_api/domain/receivable/model"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"
)

const (
	getAllReceivablesQuery   = `SELECT r.id, COALESCE(r.transaction_id,0) as transaction_id, t.transaction_code, c.name as customer_name, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date, r.created_at FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE 1=1`
	countReceivablesBase     = `SELECT COUNT(*) FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE 1=1`
	getReceivableByIDQuery   = `SELECT r.id, r.transaction_id, r.customer_id, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date, r.created_at, r.updated_at FROM receivables r WHERE r.id = ?`
	getReceivableDetailQuery = `SELECT r.id, COALESCE(r.transaction_id,0) as transaction_id, t.transaction_code, COALESCE(r.customer_id,0) as customer_id, c.name as customer_name, r.total_amount, r.paid_amount, r.remaining_amount, r.status, r.due_date, COALESCE(r.notes,'') as notes, r.created_at FROM receivables r LEFT JOIN transactions t ON r.transaction_id = t.id LEFT JOIN customers c ON r.customer_id = c.id WHERE r.id = ?`
	getReceivableItemsQuery  = `SELECT ti.product_name, ti.unit, ti.quantity, ti.price, ti.subtotal, ti.discount_item FROM transaction_items ti WHERE ti.transaction_id = ? ORDER BY ti.id`
	getReceivableStatsQuery  = `
		SELECT
			COALESCE(SUM(CASE WHEN r.status NOT IN ('paid','void') THEN r.remaining_amount ELSE 0 END),0) AS total_remaining,
			COUNT(DISTINCT CASE WHEN r.status NOT IN ('paid','void') THEN r.customer_id END) AS customer_count,
			COUNT(CASE WHEN r.status NOT IN ('paid','void') AND r.due_date IS NOT NULL AND r.due_date < CURDATE() THEN 1 END) AS overdue_count,
			COALESCE(SUM(CASE WHEN r.status NOT IN ('paid','void') AND r.due_date IS NOT NULL AND r.due_date < CURDATE() THEN r.remaining_amount ELSE 0 END),0) AS overdue_remaining
		FROM receivables r`
	getReceivableSummaryQuery = `SELECT r.customer_id, c.name as customer_name, SUM(r.total_amount) as total_receivable, SUM(r.paid_amount) as total_paid, SUM(r.remaining_amount) as total_remaining, COUNT(r.id) as count FROM receivables r LEFT JOIN customers c ON r.customer_id = c.id WHERE r.status NOT IN ('paid', 'void') GROUP BY r.customer_id, c.name ORDER BY total_remaining DESC`
	createPaymentQuery        = `INSERT INTO receivable_payments (receivable_id, payment_date, amount, notes, user_id) VALUES (?, ?, ?, ?, ?)`
	updateReceivableQuery     = `UPDATE receivables SET paid_amount = paid_amount + ?, remaining_amount = remaining_amount - ?, status = CASE WHEN remaining_amount <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = ? WHERE id = ?`
	getPaymentsQuery          = `SELECT rp.id, rp.payment_date, rp.amount, rp.payment_method, u.full_name as user_name, rp.notes FROM receivable_payments rp LEFT JOIN users u ON rp.user_id = u.id WHERE rp.receivable_id = ? ORDER BY rp.payment_date DESC`
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
	// Filter jatuh tempo. Hanya bermakna untuk piutang yang belum lunas & belum
	// dibatalkan (punya due_date terisi).
	switch req.DueTerm {
	case "overdue":
		conditions += " AND r.due_date IS NOT NULL AND r.due_date < CURDATE() AND r.status NOT IN ('paid','void')"
	case "not_due":
		conditions += " AND r.due_date IS NOT NULL AND r.due_date >= CURDATE() AND r.status NOT IN ('paid','void')"
	}

	var total int64
	if err := r.db.Raw(countReceivablesBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 20, 100)

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

	// Rincian barang dari transaksi asal (kalau transaksi induk masih ada).
	item.Items = []dto.TransactionItemBrief{}
	if item.TransactionID > 0 {
		var items []dto.TransactionItemBrief
		if err := r.db.Raw(getReceivableItemsQuery, item.TransactionID).Scan(&items).Error; err != nil {
			return nil, err
		}
		if items != nil {
			item.Items = items
		}
	}

	// Riwayat cicilan pembayaran.
	payments, err := r.GetPayments(id)
	if err != nil {
		return nil, err
	}
	item.Payments = make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		item.Payments = append(item.Payments, *p)
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

func (r *receivableRepo) GetStats() (*dto.ReceivableStats, error) {
	var stats dto.ReceivableStats
	if err := r.db.Raw(getReceivableStatsQuery).Scan(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *receivableRepo) GetPayments(receivableID int) ([]*dto.PaymentResponse, error) {
	var dataDB []*dto.PaymentResponse
	if err := r.db.Raw(getPaymentsQuery, receivableID).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *receivableRepo) CreatePayment(receivableID int, req *dto.PayRequest, userID int) error {
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return err
	}
	return r.db.Exec(createPaymentQuery, receivableID, paymentDate, req.Amount, req.Notes, userID).Error
}

func (r *receivableRepo) UpdateAfterPayment(receivableID int, amount float64) error {
	return r.db.Exec(updateReceivableQuery, amount, amount, time_helper.GetTimeNow(), receivableID).Error
}
