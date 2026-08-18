package repo

import (
	stderrors "errors"
	"fmt"

	model "pos_api/domain/expiry_batch/model"
	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	custom_errors "pos_api/errors"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const nearExpiryDays = 7

const (
	getWarningsQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       COALESCE(u.name, '') as unit_name,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		LEFT JOIN product_packages pp ON eb.package_id = pp.id
		LEFT JOIN units u ON pp.unit_id = u.id
		WHERE eb.status = 'active' AND eb.expired_date <= DATE_ADD(?, INTERVAL ? DAY)
	`
	getWarningsSearchClause = ` AND p.name LIKE ?`
	getWarningsOrderClause  = ` ORDER BY eb.expired_date ASC`

	getByProductQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       COALESCE(u.name, '') as unit_name,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		LEFT JOIN product_packages pp ON eb.package_id = pp.id
		LEFT JOIN units u ON pp.unit_id = u.id
		WHERE eb.product_id = ?
		ORDER BY eb.expired_date ASC
	`

	getExpiryBatchByIDQuery = `
		SELECT eb.id, eb.product_id, COALESCE(p.name, '') as product_name, eb.purchase_item_id,
		       COALESCE(u.name, '') as unit_name,
		       eb.qty, eb.expired_date, eb.status, eb.resolved_by, eb.resolved_at, eb.notes, eb.created_at
		FROM product_expiry_batches eb
		LEFT JOIN products p ON eb.product_id = p.id
		LEFT JOIN product_packages pp ON eb.package_id = pp.id
		LEFT JOIN units u ON pp.unit_id = u.id
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

	getPackageIDFromPurchaseItemQuery = `SELECT package_id FROM purchase_items WHERE id = ?`
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

		var packageIDPtr *int
		if err := tx.Raw(getPackageIDFromPurchaseItemQuery, batch.PurchaseItemID).Scan(&packageIDPtr).Error; err != nil {
			return err
		}
		packageID := 0
		if packageIDPtr != nil && *packageIDPtr > 0 {
			packageID = *packageIDPtr
		} else if resolved, ok := product_repo.ResolveDefaultPackageID(tx, batch.ProductID); ok {
			packageID = resolved
		}
		if packageID == 0 {
			return fmt.Errorf("produk %s tidak punya paket satuan yang bisa dipakai buat write-off", batch.ProductName)
		}

		mutationNotes := fmt.Sprintf("Write-off batch expired %s", batch.ExpiredDate.Format("2006-01-02"))
		if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
			ProductID:     batch.ProductID,
			PackageID:     packageID,
			Quantity:      batch.Qty,
			Direction:     model_product.StockOut,
			MutationType:  "expired",
			ReferenceType: "expiry_batch",
			ReferenceID:   batch.ID,
			Notes:         mutationNotes,
			UserID:        &userID,
		}); err != nil {
			if stderrors.Is(err, model_product.ErrInsufficientStock) {
				return &custom_errors.BadRequestError{Message: fmt.Sprintf(
					"Qty write-off (%.3f) melebihi stok %s yang tersedia -- tidak bisa diproses otomatis, periksa data batch ini", batch.Qty, batch.ProductName,
				)}
			}
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
