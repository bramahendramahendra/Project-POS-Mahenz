package repo

import (
	"fmt"

	dto "pos_api/domain/capital_reconciliation/dto"
	custom_errors "pos_api/errors"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const (
	// Daftar baris tinjauan + join info transaksi & produk untuk tampilan.
	listBaseQuery = `
		SELECT cr.id AS review_id, cr.transaction_item_id, cr.transaction_id,
		       COALESCE(t.transaction_code, '') AS transaction_code,
		       t.transaction_date,
		       COALESCE(cr.product_id, 0) AS product_id,
		       COALESCE(ti.product_name, '') AS product_name,
		       COALESCE(ti.unit, '') AS unit,
		       COALESCE(ti.quantity, 0) AS quantity,
		       cr.sell_price, cr.old_purchase_price, cr.suggested_price,
		       cr.reason, cr.status, cr.resolved_price, COALESCE(cr.note, '') AS note,
		       COALESCE(u.full_name, '') AS resolved_by, cr.resolved_at
		FROM transaction_cost_review cr
		LEFT JOIN transaction_items ti ON ti.id = cr.transaction_item_id
		LEFT JOIN transactions t ON t.id = cr.transaction_id
		LEFT JOIN users u ON u.id = cr.resolved_by
		WHERE 1=1`

	listCountBase = `SELECT COUNT(*) FROM transaction_cost_review cr
		LEFT JOIN transaction_items ti ON ti.id = cr.transaction_item_id
		WHERE 1=1`

	getReviewByIDQuery = `
		SELECT id, transaction_item_id, transaction_id, product_id,
		       old_purchase_price, suggested_price, sell_price, status
		FROM transaction_cost_review WHERE id = ? LIMIT 1 FOR UPDATE`

	auditInsertQuery = `
		INSERT INTO transaction_items_cogs_audit
		  (transaction_item_id, transaction_id, product_id, old_purchase_price, new_purchase_price, anchor_cost, factor, note)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	updateItemCostQuery = `UPDATE transaction_items SET purchase_price = ? WHERE id = ?`

	resolveReviewQuery = `
		UPDATE transaction_cost_review
		SET status = 'resolved', resolved_price = ?, note = ?, resolved_by = ?, resolved_at = ?, updated_at = NOW()
		WHERE id = ?`

	skipReviewQuery = `
		UPDATE transaction_cost_review
		SET status = 'skipped', note = ?, resolved_by = ?, resolved_at = ?, updated_at = NOW()
		WHERE id = ?`
)

// reviewRow menampung hasil scan getReviewByIDQuery.
type reviewRow struct {
	ID                int
	TransactionItemID int
	TransactionID     int
	ProductID         *int
	OldPurchasePrice  float64
	SuggestedPrice    float64
	SellPrice         float64
	Status            string
}

func (r *capitalReconciliationRepo) GetList(req *dto.ListRequest) ([]*dto.ListItem, int64, error) {
	cond := ""
	var args []any
	var countArgs []any

	status := req.Status
	if status == "" {
		status = "pending"
	}
	if status != "all" {
		cond += " AND cr.status = ?"
		args = append(args, status)
		countArgs = append(countArgs, status)
	}
	if req.Search != "" {
		cond += " AND (ti.product_name LIKE ? OR t.transaction_code LIKE ?)"
		like := "%" + req.Search + "%"
		args = append(args, like, like)
		// countBase tidak join transactions; tambahkan filter yang tersedia saja
	}

	// count (search di count hanya pakai product_name supaya tidak butuh join transactions)
	countCond := ""
	if status != "all" {
		countCond += " AND cr.status = ?"
	}
	if req.Search != "" {
		countCond += " AND ti.product_name LIKE ?"
		countArgs = append(countArgs, "%"+req.Search+"%")
	}
	var total int64
	if err := r.db.Raw(listCountBase+countCond, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)
	query := listBaseQuery + cond + " ORDER BY t.transaction_date DESC, cr.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var items []*dto.ListItem
	if err := r.db.Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []*dto.ListItem{}
	}
	return items, total, nil
}

func (r *capitalReconciliationRepo) GetSummary() (*dto.SummaryResponse, error) {
	type cnt struct {
		Status string
		Total  int
	}
	var rows []cnt
	if err := r.db.Raw(`SELECT status, COUNT(*) AS total FROM transaction_cost_review GROUP BY status`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := &dto.SummaryResponse{}
	for _, rw := range rows {
		switch rw.Status {
		case "pending":
			out.Pending = rw.Total
		case "resolved":
			out.Resolved = rw.Total
		case "skipped":
			out.Skipped = rw.Total
		}
		out.Total += rw.Total
	}
	return out, nil
}

// Resolve — admin menetapkan modal benar untuk 1 baris:
// update transaction_items.purchase_price + catat audit + tandai review resolved.
func (r *capitalReconciliationRepo) Resolve(req *dto.ResolveRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var rv reviewRow
		if err := tx.Raw(getReviewByIDQuery, req.ReviewID).Scan(&rv).Error; err != nil {
			return err
		}
		if rv.ID == 0 {
			return &custom_errors.NotFoundError{Message: "Baris tinjauan tidak ditemukan"}
		}
		if rv.Status != "pending" {
			return &custom_errors.BadRequestError{Message: "Baris ini sudah ditinjau sebelumnya"}
		}

		pid := any(nil)
		if rv.ProductID != nil {
			pid = *rv.ProductID
		}

		// Audit: catat perubahan modal (anchor/factor 0 karena input manual).
		note := fmt.Sprintf("koreksi modal manual (Rekonsiliasi Modal) oleh user %d", req.UserID)
		if req.Note != "" {
			note = req.Note + " (" + note + ")"
		}
		if err := tx.Exec(auditInsertQuery,
			rv.TransactionItemID, rv.TransactionID, pid,
			rv.OldPurchasePrice, req.CorrectCost, 0, 0, note,
		).Error; err != nil {
			return fmt.Errorf("audit: %w", err)
		}

		// Update modal baris transaksi.
		if err := tx.Exec(updateItemCostQuery, req.CorrectCost, rv.TransactionItemID).Error; err != nil {
			return fmt.Errorf("update modal: %w", err)
		}

		// Tandai review resolved.
		now := time_helper.GetTimeNow()
		if err := tx.Exec(resolveReviewQuery, req.CorrectCost, req.Note, req.UserID, now, req.ReviewID).Error; err != nil {
			return fmt.Errorf("update review: %w", err)
		}
		return nil
	})
}

// Skip — tandai baris sudah ditinjau tanpa mengubah modal.
func (r *capitalReconciliationRepo) Skip(req *dto.SkipRequest) error {
	var rv reviewRow
	if err := r.db.Raw(getReviewByIDQuery, req.ReviewID).Scan(&rv).Error; err != nil {
		return err
	}
	if rv.ID == 0 {
		return &custom_errors.NotFoundError{Message: "Baris tinjauan tidak ditemukan"}
	}
	if rv.Status != "pending" {
		return &custom_errors.BadRequestError{Message: "Baris ini sudah ditinjau sebelumnya"}
	}
	now := time_helper.GetTimeNow()
	return r.db.Exec(skipReviewQuery, req.Note, req.UserID, now, req.ReviewID).Error
}
