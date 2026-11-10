package repo

import (
	"fmt"
	"pos_api/domain/customer_balance/dto"
	"pos_api/domain/customer_balance/model"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const (
	getBalanceQuery       = `SELECT balance FROM customers WHERE id = ? LIMIT 1`
	updateBalanceDeduct   = `UPDATE customers SET balance = balance + ?, updated_at = ? WHERE id = ? AND balance + ? >= 0`
	updateBalanceCredit   = `UPDATE customers SET balance = balance + ?, updated_at = ? WHERE id = ?`
	insertMutationQuery   = `INSERT INTO customer_balance_mutations (customer_id, amount, balance_after, type, reference_type, reference_id, notes, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	getMutationByRef      = `SELECT id, customer_id, amount, balance_after, type, reference_type, reference_id, notes, user_id, created_at FROM customer_balance_mutations WHERE reference_type = ? AND reference_id = ? AND type = ? LIMIT 1`
	getHistoryQuery       = `SELECT m.id, m.customer_id, m.amount, m.balance_after, m.type, m.reference_type, m.reference_id, m.notes, COALESCE(u.full_name, '') as user_name, m.created_at FROM customer_balance_mutations m LEFT JOIN users u ON m.user_id = u.id WHERE m.customer_id = ? ORDER BY m.created_at DESC LIMIT ? OFFSET ?`
	countHistoryQuery     = `SELECT COUNT(*) FROM customer_balance_mutations WHERE customer_id = ?`
	getBalanceAfterUpdate = `SELECT balance FROM customers WHERE id = ? LIMIT 1`
)

func (r *customerBalanceRepo) GetBalance(customerID int) (float64, error) {
	var balance float64
	if err := r.db.Raw(getBalanceQuery, customerID).Scan(&balance).Error; err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *customerBalanceRepo) GetMutationByRef(refType string, refID int, mutationType string) (*model.CustomerBalanceMutation, error) {
	var m model.CustomerBalanceMutation
	if err := r.db.Raw(getMutationByRef, refType, refID, mutationType).Scan(&m).Error; err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *customerBalanceRepo) GetHistory(req *dto.HistoryRequest) ([]*dto.MutationResponse, int64, error) {
	var total int64
	if err := r.db.Raw(countHistoryQuery, req.CustomerID).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 20, 100)

	var data []*dto.MutationResponse
	if err := r.db.Raw(getHistoryQuery, req.CustomerID, limit, offset).Scan(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

// Topup menambah saldo customer. Harus dipanggil di dalam DB transaction dengan FOR UPDATE.
func (r *customerBalanceRepo) Topup(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error) {
	return r.mutate(customerID, amount, "topup", refType, refID, notes, userID)
}

// Deduct mengurangi saldo customer (usage). Harus dipanggil di dalam DB transaction.
func (r *customerBalanceRepo) Deduct(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error) {
	return r.mutate(customerID, -amount, "usage", refType, refID, notes, userID)
}

// Refund mengembalikan saldo ke customer. Harus dipanggil di dalam DB transaction.
func (r *customerBalanceRepo) Refund(customerID int, amount float64, refType string, refID *int, notes string, userID int) (float64, error) {
	return r.mutate(customerID, -amount, "refund", refType, refID, notes, userID)
}

// Adjust koreksi manual saldo (bisa + atau -).
func (r *customerBalanceRepo) Adjust(customerID int, amount float64, notes string, userID int) (float64, error) {
	return r.mutate(customerID, amount, "adjustment", "manual", nil, notes, userID)
}

// mutate adalah helper internal yang:
// 1. Atomic UPDATE balance (tanpa FOR UPDATE lock untuk menghindari deadlock)
// 2. Read balance baru
// 3. Insert mutasi
// Returns: balance setelah mutasi
func (r *customerBalanceRepo) mutate(customerID int, amount float64, mutType string, refType string, refID *int, notes string, userID int) (float64, error) {
	// Untuk deduction (amount negatif), pakai conditional update agar tidak bisa negatif
	// Untuk credit (amount positif), langsung update
	var query string
	if amount < 0 {
		// Deduct: pastikan saldo cukup via WHERE condition
		query = updateBalanceDeduct
	} else {
		query = updateBalanceCredit
	}

	var result *gorm.DB
	if amount < 0 {
		result = r.db.Exec(query, amount, time_helper.GetTimeNow(), customerID, amount)
	} else {
		result = r.db.Exec(query, amount, time_helper.GetTimeNow(), customerID)
	}
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("saldo tidak mencukupi")
	}

	// Baca balance terbaru
	var newBalance float64
	if err := r.db.Raw(getBalanceAfterUpdate, customerID).Scan(&newBalance).Error; err != nil {
		return 0, err
	}

	// Insert mutation
	var refTypePtr *string
	if refType != "" {
		refTypePtr = &refType
	}
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}

	if err := r.db.Exec(insertMutationQuery, customerID, amount, newBalance, mutType, refTypePtr, refID, notesPtr, userID).Error; err != nil {
		return 0, err
	}

	return newBalance, nil
}
