package repo

import (
	stderrors "errors"
	"fmt"
	"strconv"

	model_product "pos_api/domain/product/model"
	product_repo "pos_api/domain/product/repo"
	dto "pos_api/domain/supplier_return/dto"
	model "pos_api/domain/supplier_return/model"
	custom_errors "pos_api/errors"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	generateReturnCodeQuery = `SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(return_code, '-', -1) AS UNSIGNED)), 0) FROM supplier_returns WHERE DATE(return_date) = ?`
	createReturnQuery       = `INSERT INTO supplier_returns (return_code, purchase_id, supplier_id, supplier_name, return_date, total_return_amount, reason, status, user_id, notes) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?)`
	// package_id diisi dari purchase_items.package_id lewat purchase_item_id
	// (celah #14) -- FE tidak perlu diubah, purchase_item_id sudah dikirim.
	createReturnItemQuery      = `INSERT INTO supplier_return_items (return_id, purchase_item_id, product_id, product_name, package_id, quantity, unit, purchase_price, subtotal) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	updateReturnStatusQuery    = `UPDATE supplier_returns SET status = ?, notes = ?, updated_at = NOW() WHERE id = ?`
	approveReturnStatusQuery   = `UPDATE supplier_returns SET status = 'approved', updated_at = NOW() WHERE id = ?`
	reservePackageStockQuery   = `UPDATE product_packages SET reserved_qty = reserved_qty + ?, updated_at = NOW() WHERE id = ?`
	releasePackageStockQuery   = `UPDATE product_packages SET reserved_qty = GREATEST(reserved_qty - ?, 0), updated_at = NOW() WHERE id = ?`
	getPackageIDFromPurchaseItemQuery = `SELECT package_id FROM purchase_items WHERE id = ?`
	getReturnItemsQuery         = `SELECT sri.id, sri.product_id, sri.product_name, sri.package_id, sri.quantity, sri.unit, sri.purchase_price, sri.subtotal FROM supplier_return_items sri WHERE sri.return_id = ?`
	checkReturnApprovedQuery    = `SELECT status FROM supplier_returns WHERE id = ?`
	getPurchaseIDAndAmountQuery = `SELECT purchase_id, total_return_amount FROM supplier_returns WHERE id = ?`
	reducePurchaseDebtQuery     = `UPDATE purchases SET remaining_amount = GREATEST(remaining_amount - ?, 0), payment_status = CASE WHEN remaining_amount <= 0 THEN 'paid' WHEN paid_amount > 0 THEN 'partial' ELSE 'unpaid' END, updated_at = NOW() WHERE id = ?`
	getReturnByIDQuery          = `SELECT sr.id, sr.return_code, sr.purchase_id, sr.supplier_id, sr.supplier_name, sr.return_date, sr.total_return_amount, sr.reason, sr.status, u.full_name as user_name, sr.notes FROM supplier_returns sr LEFT JOIN users u ON sr.user_id = u.id WHERE sr.id = ?`
	getAllReturnsBase           = `SELECT sr.id, sr.return_code, sr.purchase_id, sr.supplier_id, sr.supplier_name, sr.return_date, sr.total_return_amount, sr.reason, sr.status, u.full_name as user_name, sr.notes FROM supplier_returns sr LEFT JOIN users u ON sr.user_id = u.id WHERE 1=1`
	countReturnsBase            = `SELECT COUNT(*) FROM supplier_returns sr WHERE 1=1`
	getPurchaseDateQuery        = `SELECT purchase_date FROM purchases WHERE id = ? LIMIT 1`
	getPurchaseStatusQuery      = `SELECT status FROM purchases WHERE id = ? LIMIT 1`
	getPurchaseItemQtyQuery     = `SELECT quantity FROM purchase_items WHERE id = ? AND purchase_id = ? LIMIT 1 FOR UPDATE`
	getTotalReturnedQtyQuery    = `SELECT COALESCE(SUM(sri.quantity), 0) FROM supplier_return_items sri JOIN supplier_returns sr ON sri.return_id = sr.id WHERE sri.purchase_item_id = ? AND sr.status IN ('pending', 'approved')`
	deleteReturnItemsQuery      = `DELETE FROM supplier_return_items WHERE return_id = ?`
	deleteReturnQuery           = `DELETE FROM supplier_returns WHERE id = ?`
)

// resolveReturnItemPackageID mengisi celah #14: kalau supplier_return_items
// lama belum punya package_id (data sebelum migrasi), cari dari
// purchase_items.package_id lewat purchase_item_id sebagai fallback.
func resolveReturnItemPackageID(tx *gorm.DB, item model.SupplierReturnItem) (int, error) {
	if item.PackageID != nil && *item.PackageID > 0 {
		return *item.PackageID, nil
	}
	var pkgID *int
	if err := tx.Raw(getPackageIDFromPurchaseItemQuery, item.PurchaseItemID).Scan(&pkgID).Error; err != nil {
		return 0, err
	}
	if pkgID != nil && *pkgID > 0 {
		return *pkgID, nil
	}
	if resolved, ok := product_repo.ResolveDefaultPackageID(tx, item.ProductID); ok {
		return resolved, nil
	}
	return 0, fmt.Errorf("produk %s tidak punya paket satuan yang bisa dipakai buat retur", item.ProductName)
}

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

func (r *supplierReturnRepo) GetPurchaseStatus(purchaseID int) (string, error) {
	var status string
	err := r.db.Raw(getPurchaseStatusQuery, purchaseID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func isDuplicateReturnCodeError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return stderrors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *supplierReturnRepo) Create(req *dto.CreateSupplierReturnRequest) (*model.SupplierReturnRow, error) {
	// GET_LOCK per tanggal -- lihat komentar panjang di transaction_repo.go's
	// createOnce(): retry-on-duplicate-key SENDIRIAN terbukti nyata belum cukup
	// di bawah 2 request yang benar-benar konkuren.
	lockName := fmt.Sprintf("returncode:%s", time_helper.GetTimeNow().Format("2006-01-02"))

	const maxCodeRetries = 5
	var returnID int
	var err error
	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		returnID = 0
		err = r.db.Connection(func(conn *gorm.DB) error {
			var locked int
			if err := conn.Raw(`SELECT GET_LOCK(?, 5)`, lockName).Scan(&locked).Error; err != nil {
				return err
			}
			if locked != 1 {
				return fmt.Errorf("gagal mendapatkan lock generate kode retur (timeout)")
			}
			defer conn.Exec(`SELECT RELEASE_LOCK(?)`, lockName)
			return r.createOnce(conn, req, &returnID)
		})
		if err == nil || !isDuplicateReturnCodeError(err) {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	data, err := r.GetByID(returnID)
	return data, err
}

func (r *supplierReturnRepo) createOnce(conn *gorm.DB, req *dto.CreateSupplierReturnRequest, returnID *int) error {
	return conn.Transaction(func(tx *gorm.DB) error {
		now := time_helper.GetTimeNow()
		var count int
		if err := tx.Raw(generateReturnCodeQuery, now.Format("2006-01-02")).Scan(&count).Error; err != nil {
			return err
		}
		code := fmt.Sprintf("RTR-%s-%03d", now.Format("20060102"), count+1)

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

		if err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(returnID).Error; err != nil {
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
				// strconv.FormatFloat(..., -1, ...) -- BUKAN %.0f -- supaya sisa pecahan
				// (mis. satuan kontinu seperti Kilogram) tampil apa adanya (0.5), bukan
				// dibulatkan ke 0 di pesan errornya (meski validasi angkanya sendiri sudah
				// benar, pesannya jadi menyesatkan / kelihatan seperti sisa 0 padahal bukan).
				return &custom_errors.BadRequestError{
					Message: fmt.Sprintf("Jumlah retur %s melebihi sisa yang bisa diretur (maks %s)", item.ProductName, strconv.FormatFloat(sisaQty, 'f', -1, 64)),
				}
			}

			var purchaseItemPackageID *int
			if err := tx.Raw(getPackageIDFromPurchaseItemQuery, item.PurchaseItemID).Scan(&purchaseItemPackageID).Error; err != nil {
				return err
			}
			packageID := purchaseItemPackageID
			if packageID == nil || *packageID == 0 {
				resolved, ok := product_repo.ResolveDefaultPackageID(tx, item.ProductID)
				if !ok {
					return fmt.Errorf("produk %s tidak punya paket satuan yang bisa dipakai buat retur", item.ProductName)
				}
				packageID = &resolved
			}

			subtotal := item.PurchasePrice * item.Quantity
			err = tx.Exec(createReturnItemQuery,
				*returnID, item.PurchaseItemID, item.ProductID, item.ProductName, packageID,
				item.Quantity, item.Unit, item.PurchasePrice, subtotal,
			).Error
			if err != nil {
				return err
			}

			// Reserve di level PACKAGE (bukan lagi products.reserved_qty) --
			// supaya bagian yang ditahan retur ini tidak ikut kepakai/kepinjam
			// jalur lain (celah #7/#13).
			if err = tx.Exec(reservePackageStockQuery, item.Quantity, *packageID).Error; err != nil {
				return err
			}
		}

		return nil
	})
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

		notes := fmt.Sprintf("Supplier Return #%d", id)
		for _, item := range items {
			packageID, err := resolveReturnItemPackageID(tx, item)
			if err != nil {
				return err
			}

			// Lepas reservasi DULU, baru kurangi stok. ApplyStockDelta
			// menghitung pool "boleh dipakai" = stock - reserved_qty --
			// kalau reservasi belum dilepas, qty yang mau dikurangi ini
			// masih dianggap "ditahan" dan akan DITOLAK sebagai stok tidak
			// cukup, padahal justru inilah stok yang mau dikonsumsi retur.
			if err := tx.Exec(releasePackageStockQuery, item.Quantity, packageID).Error; err != nil {
				return err
			}

			if _, err := product_repo.ApplyStockDelta(tx, product_repo.ApplyStockDeltaParams{
				ProductID:     item.ProductID,
				PackageID:     packageID,
				Quantity:      item.Quantity,
				Direction:     model_product.StockOut,
				MutationType:  "return",
				ReferenceType: "supplier_return",
				ReferenceID:   id,
				Notes:         notes,
				UserID:        &userID,
			}); err != nil {
				if stderrors.Is(err, model_product.ErrInsufficientStock) {
					return &custom_errors.BadRequestError{
						Message: fmt.Sprintf("Stok %s tidak mencukupi untuk retur", item.ProductName),
					}
				}
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
			packageID, err := resolveReturnItemPackageID(tx, item)
			if err != nil {
				return err
			}
			if err := tx.Exec(releasePackageStockQuery, item.Quantity, packageID).Error; err != nil {
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
