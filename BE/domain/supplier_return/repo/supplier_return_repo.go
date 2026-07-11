package repo

import (
	"fmt"
	"time"

	dto "pos_api/domain/supplier_return/dto"
	model "pos_api/domain/supplier_return/model"
	custom_errors "pos_api/errors"
	request_helper "pos_api/helper/request"

	"gorm.io/gorm"
)

const (
	generateReturnCodeQuery       = `SELECT COUNT(*) FROM supplier_returns WHERE DATE(return_date) = ?`
	createReturnQuery             = `INSERT INTO supplier_returns (return_code, purchase_id, supplier_id, supplier_name, return_date, total_return_amount, reason, status, user_id, notes) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?)`
	createReturnItemQuery         = `INSERT INTO supplier_return_items (return_id, purchase_item_id, product_id, product_name, quantity, unit, purchase_price, subtotal) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	updateReturnStatusQuery       = `UPDATE supplier_returns SET status = ?, notes = ?, updated_at = NOW() WHERE id = ?`
	approveReturnStatusQuery      = `UPDATE supplier_returns SET status = 'approved', updated_at = NOW() WHERE id = ?`
	reduceStockQuery              = `UPDATE products SET stock = stock - ?, updated_at = NOW() WHERE id = ?`
	reserveStockQuery             = `UPDATE products SET reserved_qty = reserved_qty + ?, updated_at = NOW() WHERE id = ?`
	releaseReservedStockQuery     = `UPDATE products SET reserved_qty = GREATEST(reserved_qty - ?, 0), updated_at = NOW() WHERE id = ?`
	getReturnItemsQuery           = `SELECT sri.id, sri.product_id, sri.product_name, sri.quantity, sri.unit, sri.purchase_price, sri.subtotal FROM supplier_return_items sri WHERE sri.return_id = ?`
	checkReturnApprovedQuery      = `SELECT status FROM supplier_returns WHERE id = ?`
	getPurchaseIDAndAmountQuery   = `SELECT purchase_id, total_return_amount FROM supplier_returns WHERE id = ?`
	reducePurchaseDebtQuery       = `UPDATE purchases SET remaining_amount = GREATEST(remaining_amount - ?, 0), payment_status = CASE WHEN remaining_amount <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
	getReturnByIDQuery            = `SELECT sr.id, sr.return_code, sr.purchase_id, sr.supplier_id, sr.supplier_name, sr.return_date, sr.total_return_amount, sr.reason, sr.status, u.full_name as user_name, sr.notes FROM supplier_returns sr LEFT JOIN users u ON sr.user_id = u.id WHERE sr.id = ?`
	getAllReturnsBase  = `SELECT sr.id, sr.return_code, sr.purchase_id, sr.supplier_id, sr.supplier_name, sr.return_date, sr.total_return_amount, sr.reason, sr.status, u.full_name as user_name, sr.notes FROM supplier_returns sr LEFT JOIN users u ON sr.user_id = u.id WHERE 1=1`
	countReturnsBase  = `SELECT COUNT(*) FROM supplier_returns sr WHERE 1=1`
	getPurchaseDateQuery          = `SELECT purchase_date FROM purchases WHERE id = ? LIMIT 1`
	getPurchaseItemQtyQuery       = `SELECT quantity FROM purchase_items WHERE id = ? AND purchase_id = ? LIMIT 1 FOR UPDATE`
	getTotalReturnedQtyQuery      = `SELECT COALESCE(SUM(sri.quantity), 0) FROM supplier_return_items sri JOIN supplier_returns sr ON sri.return_id = sr.id WHERE sri.purchase_item_id = ? AND sr.status IN ('pending', 'approved')`
	deleteReturnItemsQuery        = `DELETE FROM supplier_return_items WHERE return_id = ?`
	deleteReturnQuery             = `DELETE FROM supplier_returns WHERE id = ?`
	getProductStockForUpdateQuery = `SELECT stock FROM products WHERE id = ? LIMIT 1 FOR UPDATE`
	createStockMutationQuery      = `INSERT INTO stock_mutations (product_id, mutation_type, quantity, stock_before, stock_after, reference_type, reference_id, notes, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func (r *supplierReturnRepo) GetAll(req *dto.SupplierReturnListRequest) ([]*model.SupplierReturnRow, int64, error) {
	var args []any
	conditions := ""

	if req.StartDate != "" {
		conditions += " AND DATE(sr.return_date) >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		conditions += " AND DATE(sr.return_date) <= ?"
		args = append(args, req.EndDate)
	}
	if req.SupplierID != nil {
		conditions += " AND sr.supplier_id = ?"
		args = append(args, *req.SupplierID)
	}
	if req.Status != "" {
		conditions += " AND sr.status = ?"
		args = append(args, req.Status)
	}

	var total int64
	if err := r.db.Raw(countReturnsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	allowedSortFields := map[string]string{
		"return_date":         "sr.return_date",
		"total_return_amount": "sr.total_return_amount",
		"supplier_name":       "sr.supplier_name",
		"status":              "sr.status",
	}
	const defaultOrder = " ORDER BY sr.return_date DESC, sr.id DESC"

	query := getAllReturnsBase + conditions + request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, defaultOrder) + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.SupplierReturnRow
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *supplierReturnRepo) GetByID(id int) (*model.SupplierReturnRow, error) {
	var dataDB model.SupplierReturnRow
	if err := r.db.Raw(getReturnByIDQuery, id).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}

	items, err := r.GetItems(id)
	if err != nil {
		return nil, err
	}
	dataDB.Items = items

	return &dataDB, nil
}

func (r *supplierReturnRepo) GetStatus(id int) (string, error) {
	var status string
	err := r.db.Raw(checkReturnApprovedQuery, id).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *supplierReturnRepo) GetItems(returnID int) ([]model.SupplierReturnItem, error) {
	var dataDB []model.SupplierReturnItem
	if err := r.db.Raw(getReturnItemsQuery, returnID).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *supplierReturnRepo) GetPurchaseDate(purchaseID int) (string, error) {
	var purchaseDate string
	err := r.db.Raw(getPurchaseDateQuery, purchaseID).Scan(&purchaseDate).Error
	if err != nil {
		return "", err
	}
	return purchaseDate, nil
}

func (r *supplierReturnRepo) Create(req *dto.CreateSupplierReturnRequest) (*model.SupplierReturnRow, error) {
	var returnID int

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		today := time.Now().Format("2006-01-02")
		var count int
		if err := tx.Raw(generateReturnCodeQuery, today).Scan(&count).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("RTR-%s-%03d", time.Now().Format("20060102"), count+1)

		var totalAmount float64
		for _, item := range req.Items {
			totalAmount += item.PurchasePrice * item.Quantity
		}

		err := tx.Exec(createReturnQuery,
			code, req.PurchaseID, req.SupplierID, req.SupplierName,
			req.ReturnDate, totalAmount, req.Reason, req.UserID, req.Notes,
		).Error
		if err != nil {
			return err
		}

		if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&returnID).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			var purchaseQty float64
			if err := tx.Raw(getPurchaseItemQtyQuery, item.PurchaseItemID, req.PurchaseID).Scan(&purchaseQty).Error; err != nil {
				return &custom_errors.NotFoundError{Message: "Item pembelian tidak ditemukan"}
			}

			var alreadyReturned float64
			if err := tx.Raw(getTotalReturnedQtyQuery, item.PurchaseItemID).Scan(&alreadyReturned).Error; err != nil {
				return err
			}

			sisaQty := purchaseQty - alreadyReturned
			if item.Quantity > sisaQty {
				return &custom_errors.BadRequestError{
					Message: fmt.Sprintf("Jumlah retur %s melebihi sisa yang bisa diretur (maks %.0f)", item.ProductName, sisaQty),
				}
			}

			subtotal := item.PurchasePrice * item.Quantity
			err = tx.Exec(createReturnItemQuery,
				returnID, item.PurchaseItemID, item.ProductID, item.ProductName,
				item.Quantity, item.Unit, item.PurchasePrice, subtotal,
			).Error
			if err != nil {
				return err
			}

			if err = tx.Exec(reserveStockQuery, item.Quantity, item.ProductID).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}
	data, err := r.GetByID(returnID)
	return data, err
}

func (r *supplierReturnRepo) UpdateStatus(id int, status, notes string) error {
	err := r.db.Exec(updateReturnStatusQuery, status, notes, id).Error
	return err
}

func (r *supplierReturnRepo) ApproveWithStockReduction(id int, userID int) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(approveReturnStatusQuery, id).Error
		if err != nil {
			return err
		}

		var dataReturn model.SupplierReturnPurchaseRef
		err = tx.Raw(getPurchaseIDAndAmountQuery, id).Scan(&dataReturn).Error
		if err != nil {
			return err
		}

		items, err := r.GetItems(id)
		if err != nil {
			return err
		}

		for _, item := range items {
			var stockBefore float64
			err := tx.Raw(getProductStockForUpdateQuery, item.ProductID).Scan(&stockBefore).Error
			if err != nil {
				return err
			}

			if stockBefore < item.Quantity {
				return &custom_errors.BadRequestError{
					Message: fmt.Sprintf("stok %s tidak mencukupi untuk retur (stok saat ini: %.0f)", item.ProductName, stockBefore),
				}
			}

			err = tx.Exec(reduceStockQuery, item.Quantity, item.ProductID).Error
			if err != nil {
				return err
			}

			if err = tx.Exec(releaseReservedStockQuery, item.Quantity, item.ProductID).Error; err != nil {
				return err
			}

			stockAfter := stockBefore - item.Quantity
			notes := fmt.Sprintf("Supplier Return #%d", id)
			err = tx.Exec(createStockMutationQuery,
				item.ProductID, "return", item.Quantity, stockBefore, stockAfter,
				"supplier_return", id, notes, userID,
			).Error
			if err != nil {
				return err
			}
		}

		err = tx.Exec(reducePurchaseDebtQuery, dataReturn.TotalReturnAmount, dataReturn.PurchaseID).Error
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// ReleaseReservedStock melepas reserved_qty yang dibuat saat retur ini dibuat (status pending),
// dipanggil saat retur ditolak (rejected) atau dihapus, tanpa mengubah stock fisik.
func (r *supplierReturnRepo) ReleaseReservedStock(id int) error {
	items, err := r.GetItems(id)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Exec(releaseReservedStockQuery, item.Quantity, item.ProductID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *supplierReturnRepo) Delete(req *dto.GetSupplierReturnByIDRequest) error {
	if err := r.ReleaseReservedStock(req.ID); err != nil {
		return err
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(deleteReturnItemsQuery, req.ID).Error
		if err != nil {
			return err
		}
		err = tx.Exec(deleteReturnQuery, req.ID).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
