package repo

import (
	"fmt"

	model "pos_api/domain/expiry_batch/model"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const nearExpiryDays = 7

const (
	getWarningsQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		WHERE eb.status = 'active' AND eb.expired_date <= DATE_ADD(?, INTERVAL ? DAY)
	`
	getWarningsSearchClause = ` AND p.name LIKE ?`
	getWarningsOrderClause  = ` ORDER BY eb.expired_date ASC`

	getByProductQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		WHERE eb.product_id = ?
		ORDER BY eb.expired_date ASC
	`

	getExpiryBatchByIDQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		WHERE eb.id = ?
	`

	confirmExpiryBatchQuery = `
		UPDATE product_expiry_batches
		SET status = 'cleared', resolved_by = ?, resolved_at = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	writeOffExpiryBatchQuery = `
		UPDATE product_expiry_batches
		SET status = 'written_off', resolved_by = ?, resolved_at = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	getProductStockForWriteOffQuery = `SELECT stock FROM products WHERE id = ? FOR UPDATE`
	deductStockForWriteOffQuery     = `UPDATE products SET stock = GREATEST(stock - ?, 0), updated_at = ? WHERE id = ?`
	createExpiredStockMutationQuery = `INSERT INTO stock_mutations (product_id, mutation_type, quantity, stock_before, stock_after, reference_type, reference_id, notes, user_id) VALUES (?, 'expired', ?, ?, ?, 'expiry_batch', ?, ?, ?)`
)

func (r *expiryBatchRepo) GetWarnings(search string) ([]*model.ExpiryBatch, error) {
	query := getWarningsQuery
	args := []any{time_helper.ToSQLDate(time_helper.GetTimeNow()), nearExpiryDays}
	if search != "" {
		query += getWarningsSearchClause
		args = append(args, "%"+search+"%")
	}
	query += getWarningsOrderClause

	var rows []*model.ExpiryBatch
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *expiryBatchRepo) GetByProduct(productID int) ([]*model.ExpiryBatch, error) {
	var rows []*model.ExpiryBatch
	if err := r.db.Raw(getByProductQuery, productID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *expiryBatchRepo) GetByID(id int) (*model.ExpiryBatch, error) {
	var row model.ExpiryBatch
	result := r.db.Raw(getExpiryBatchByIDQuery, id).Scan(&row)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &row, nil
}

func (r *expiryBatchRepo) Confirm(id, userID int, notes string) error {
	now := time_helper.GetTimeNow()
	return r.db.Exec(confirmExpiryBatchQuery, userID, now, notes, now, id).Error
}

func (r *expiryBatchRepo) WriteOff(id, userID int, notes string) error {
	now := time_helper.GetTimeNow()
	return r.db.Transaction(func(tx *gorm.DB) error {
		batch, err := r.byIDForUpdate(tx, id)
		if err != nil {
			return err
		}

		var stockBefore float64
		if err := tx.Raw(getProductStockForWriteOffQuery, batch.ProductID).Scan(&stockBefore).Error; err != nil {
			return err
		}

		if err := tx.Exec(deductStockForWriteOffQuery, batch.Qty, now, batch.ProductID).Error; err != nil {
			return err
		}
		stockAfter := stockBefore - batch.Qty
		if stockAfter < 0 {
			stockAfter = 0
		}

		mutationNotes := fmt.Sprintf("Write-off batch expired %s", batch.ExpiredDate.Format("2006-01-02"))
		if err := tx.Exec(createExpiredStockMutationQuery,
			batch.ProductID, batch.Qty, stockBefore, stockAfter, batch.ID, mutationNotes, userID,
		).Error; err != nil {
			return err
		}

		return tx.Exec(writeOffExpiryBatchQuery, userID, now, notes, now, id).Error
	})
}

func (r *expiryBatchRepo) byIDForUpdate(tx *gorm.DB, id int) (*model.ExpiryBatch, error) {
	var row model.ExpiryBatch
	if err := tx.Raw(getExpiryBatchByIDQuery+" FOR UPDATE", id).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
