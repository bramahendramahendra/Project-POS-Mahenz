package repo

import (
	stderrors "errors"
	"fmt"
	"math"
	"time"

	model_product "coffee_api/domain/product/model"
	product_repo "coffee_api/domain/product/repo"
	queue_repo "coffee_api/domain/queue_number/repo"
	"coffee_api/domain/transaction/dto"
	"coffee_api/domain/transaction/model"
	time_helper "coffee_api/helper/time"
)

// Status pesanan terbuka (migrasi 008, keputusan 27, 105, 107).
const (
	StatusOpen      = "open"
	StatusCancelled = "cancelled"
	StatusMerged    = "merged"

	StockActionReturn = "return"
	StockActionWaste  = "waste"
)

// ErrNotOpen: pesanan sudah dibayar, dibatalkan, atau digabung (diperiksa saat baris dikunci).
var ErrNotOpen = stderrors.New("pesanan bukan pesanan terbuka")

// LockedBill: data pesanan terbuka yang dikunci (FOR UPDATE) di awal setiap operasi.
type LockedBill struct {
	ID              int
	TransactionCode string
	Status          string
	OrderType       string
	TableID         *int
	// Tanggal pesanan (YYYY-MM-DD) dan pembuatnya, untuk aturan batas hari (tahap 4.9).
	BillDate        string
	TransactionDate time.Time
	CreatedBy       *int
}

// IsOverdue: pesanan terbuka dari hari sebelumnya (tertunda, keputusan 45, 113).
func (b *LockedBill) IsOverdue() bool {
	return b.BillDate < time_helper.GetTimeNow().Format("2006-01-02")
}

const (
	lockBillQuery = `SELECT id, transaction_code, status, COALESCE(order_type, '') AS order_type, table_id,
		DATE_FORMAT(transaction_date, '%Y-%m-%d') AS bill_date, transaction_date, created_by
		FROM transactions WHERE id = ? LIMIT 1 FOR UPDATE`
	lockTableQuery     = `SELECT name, is_active FROM dining_tables WHERE id = ? LIMIT 1 FOR UPDATE`
	openOnTableQuery   = `SELECT COALESCE(MAX(id), 0) FROM transactions WHERE table_id = ? AND status = 'open' AND id <> ?`
	recalcOpenQuery    = `UPDATE transactions SET subtotal = (SELECT COALESCE(SUM(subtotal), 0) FROM transaction_items WHERE transaction_id = ?), total_amount = (SELECT COALESCE(SUM(subtotal), 0) FROM transaction_items WHERE transaction_id = ?), updated_at = ? WHERE id = ?`
	getItemQuery       = `SELECT id, transaction_id, menu_item_id, variant_id, COALESCE(product_id, 0) AS product_id, product_name, COALESCE(variant_name, '') AS variant_name, COALESCE(notes, '') AS notes, quantity, unit, price, purchase_price, subtotal, discount_item, conversion_qty, unit_id FROM transaction_items WHERE id = ? AND transaction_id = ? LIMIT 1`
	deleteItemQuery    = `DELETE FROM transaction_items WHERE id = ?`
	reduceItemQuery    = `UPDATE transaction_items SET quantity = ?, discount_item = ?, subtotal = ? WHERE id = ?`
	insertCancelQuery  = `INSERT INTO transaction_item_cancellations (transaction_id, item_name, variant_name, quantity, price, stock_action, reason, user_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getCancelsQuery    = `SELECT c.item_name, COALESCE(c.variant_name, '') AS variant_name, c.quantity, c.price, COALESCE(c.stock_action, '') AS stock_action, COALESCE(c.reason, '') AS reason, COALESCE(u.full_name, '') AS user_name, c.created_at FROM transaction_item_cancellations c LEFT JOIN users u ON u.id = c.user_id WHERE c.transaction_id = ? ORDER BY c.id`
	moveTableQuery     = `UPDATE transactions SET table_id = ?, table_name = ?, updated_at = ? WHERE id = ?`
	moveItemsQuery     = `UPDATE transaction_items SET transaction_id = ? WHERE transaction_id = ?`
	markMergedQuery    = `UPDATE transactions SET status = 'merged', merged_into_id = ?, subtotal = 0, total_amount = 0, updated_at = ? WHERE id = ?`
	markCancelledQuery = `UPDATE transactions SET status = 'cancelled', cancelled_by = ?, cancelled_at = ?, cancel_reason = ?, updated_at = ? WHERE id = ?`
	payOpenQuery       = `UPDATE transactions SET status = 'completed', user_id = ?, settled_by = ?, shift_id = ?, transaction_date = ?, subtotal = ?, discount = ?, tax = ?, total_amount = ?, payment_method = ?, payment_amount = ?, change_amount = ?, customer_id = ?, is_credit = ?, updated_at = ? WHERE id = ?`
	openListQuery      = `SELECT t.id, t.transaction_code, COALESCE(t.order_type, '') AS order_type, t.table_id, COALESCE(t.table_name, '') AS table_name,
		COALESCE(t.queue_number, 0) AS queue_number, COALESCE(t.order_name, '') AS order_name, t.subtotal,
		DATE_FORMAT(t.transaction_date, '%Y-%m-%d') AS bill_date,
		(SELECT COALESCE(SUM(quantity), 0) FROM transaction_items i WHERE i.transaction_id = t.id) AS item_count,
		COALESCE(u.full_name, '') AS created_by_name, t.created_at
		FROM transactions t LEFT JOIN users u ON u.id = t.created_by
		WHERE t.status = 'open' ORDER BY t.created_at, t.id`
	getOpenItemsForUpdate = `SELECT id FROM transaction_items WHERE transaction_id = ? FOR UPDATE`
)

// LockBill mengunci baris pesanan; nil bila tidak ada.
func (r *transactionRepo) LockBill(id int) (*LockedBill, error) {
	var b LockedBill
	res := r.db.Raw(lockBillQuery, id).Scan(&b)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &b, nil
}

// LockTableForOrder mengunci baris meja (menyerialkan pemakaian meja) lalu memeriksa apakah ada pesanan
// terbuka lain di meja itu (keputusan 103). occupiedBy = id pesanan terbuka lain, 0 bila kosong.
func (r *transactionRepo) LockTableForOrder(tableID, excludeBillID int) (name string, active, found bool, occupiedBy int, err error) {
	var row struct {
		Name     string
		IsActive bool
	}
	res := r.db.Raw(lockTableQuery, tableID).Scan(&row)
	if res.Error != nil {
		return "", false, false, 0, res.Error
	}
	if res.RowsAffected == 0 {
		return "", false, false, 0, nil
	}
	if err := r.db.Raw(openOnTableQuery, tableID, excludeBillID).Scan(&occupiedBy).Error; err != nil {
		return "", false, false, 0, err
	}
	return row.Name, row.IsActive, true, occupiedBy, nil
}

// ListOpen: semua pesanan terbuka, terlama dulu (semua kasir, keputusan 106).
func (r *transactionRepo) ListOpen() ([]dto.OpenBillSummary, error) {
	var rows []dto.OpenBillSummary
	if err := r.db.Raw(openListQuery).Scan(&rows).Error; err != nil {
		return nil, err
	}
	today := time_helper.GetTimeNow().Format("2006-01-02")
	for i := range rows {
		if rows[i].QueueNumber > 0 {
			rows[i].QueueLabel = queue_repo.FormatQueueNumber(rows[i].QueueNumber)
		}
		rows[i].IsOverdue = rows[i].BillDate < today
	}
	if rows == nil {
		rows = []dto.OpenBillSummary{}
	}
	return rows, nil
}

func (r *transactionRepo) recalcOpen(id int) error {
	return r.db.Exec(recalcOpenQuery, id, id, time_helper.GetTimeNow(), id).Error
}

func (r *transactionRepo) getCancellations(id int) ([]dto.ItemCancellationResponse, error) {
	var rows []dto.ItemCancellationResponse
	if err := r.db.Raw(getCancelsQuery, id).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *transactionRepo) GetItem(billID, itemID int) (*model.TransactionItem, error) {
	var item model.TransactionItem
	res := r.db.Raw(getItemQuery, itemID, billID).Scan(&item)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &item, nil
}

// AddOpenItems menambah baris ke pesanan terbuka (stok barang jadi dipotong, keputusan 40).
func (r *transactionRepo) AddOpenItems(bill *LockedBill, items []dto.CreateTransactionItemRequest, userID int) error {
	for i := range items {
		if err := r.insertItem(bill.ID, bill.TransactionCode, &items[i], userID); err != nil {
			return err
		}
	}
	return r.recalcOpen(bill.ID)
}

// restoreItemStock: barang jadi kembali ke stok (mutasi void); bila waste, langsung dikurangi lagi
// sebagai mutasi waste beserta alasannya (keputusan 100, 105). Racikan tidak menyentuh stok.
func (r *transactionRepo) restoreItemStock(item *model.TransactionItem, qty float64, action, reason string, billID, userID int, label string) error {
	if item.ProductID == 0 {
		return nil
	}
	packageID := 0
	if item.UnitID != nil && *item.UnitID > 0 {
		packageID = *item.UnitID
	} else if resolved, ok := product_repo.ResolveDefaultPackageID(r.db, item.ProductID); ok {
		packageID = resolved
	}
	if packageID == 0 {
		return fmt.Errorf("produk %s tidak punya satuan stok", item.ProductName)
	}
	if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
		ProductID: item.ProductID, PackageID: packageID, Quantity: qty, Direction: model_product.StockIn,
		MutationType: "void", ReferenceType: "transaction", ReferenceID: billID,
		Notes: fmt.Sprintf("%s pesanan ID %d", label, billID), UserID: &userID,
	}); err != nil {
		return product_repo.WrapStockError(err, item.ProductID)
	}
	if action != StockActionWaste {
		return nil
	}
	if _, err := product_repo.ApplyStockDelta(r.db, product_repo.ApplyStockDeltaParams{
		ProductID: item.ProductID, PackageID: packageID, Quantity: qty, Direction: model_product.StockOut,
		MutationType: MutationWaste, ReferenceType: "transaction", ReferenceID: billID,
		Notes: fmt.Sprintf("Dibuang (%s pesanan ID %d): %s", label, billID, reason), UserID: &userID,
	}); err != nil {
		return product_repo.WrapStockError(err, item.ProductID)
	}
	return nil
}

func (r *transactionRepo) insertCancellation(billID int, item *model.TransactionItem, qty float64, action, reason string, userID int, now time.Time) error {
	var variantName, stockAction, reasonPtr *string
	if item.VariantName != "" {
		variantName = &item.VariantName
	}
	if item.ProductID > 0 {
		stockAction = &action
	}
	if reason != "" {
		reasonPtr = &reason
	}
	return r.db.Exec(insertCancelQuery, billID, item.ProductName, variantName, qty, item.Price, stockAction, reasonPtr, userID, now).Error
}

// RemoveOpenItem menghapus sebagian/seluruh jumlah satu baris pesanan terbuka. Diskon item dikurangi
// sebanding jumlah yang tersisa.
func (r *transactionRepo) RemoveOpenItem(bill *LockedBill, item *model.TransactionItem, qty float64, action, reason string, userID int) error {
	now := time_helper.GetTimeNow()
	if qty >= item.Quantity {
		qty = item.Quantity
		if err := r.db.Exec(deleteItemQuery, item.ID).Error; err != nil {
			return err
		}
	} else {
		left := item.Quantity - qty
		discount := math.Round(item.DiscountItem*left/item.Quantity*100) / 100
		if err := r.db.Exec(reduceItemQuery, left, discount, item.Price*left-discount, item.ID).Error; err != nil {
			return err
		}
	}
	if err := r.insertCancellation(bill.ID, item, qty, action, reason, userID, now); err != nil {
		return err
	}
	if err := r.restoreItemStock(item, qty, action, reason, bill.ID, userID, "Hapus item"); err != nil {
		return err
	}
	return r.recalcOpen(bill.ID)
}

func (r *transactionRepo) MoveOpenTable(billID, tableID int, tableName string) error {
	return r.db.Exec(moveTableQuery, tableID, tableName, time_helper.GetTimeNow(), billID).Error
}

// MergeOpen memindahkan semua item sumber ke tujuan (keputusan 107). Stok tidak berubah.
func (r *transactionRepo) MergeOpen(targetID, sourceID int) error {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(moveItemsQuery, targetID, sourceID).Error; err != nil {
		return err
	}
	if err := r.db.Exec(markMergedQuery, targetID, now, sourceID).Error; err != nil {
		return err
	}
	return r.recalcOpen(targetID)
}

// CancelOpen membatalkan pesanan terbuka: setiap item dicatat, barang jadi kembali ke stok atau dibuang.
func (r *transactionRepo) CancelOpen(bill *LockedBill, wasteItemIDs map[int]bool, reason string, userID int) error {
	now := time_helper.GetTimeNow()
	if err := r.db.Exec(getOpenItemsForUpdate, bill.ID).Error; err != nil {
		return err
	}
	items, err := r.GetItems(bill.ID)
	if err != nil {
		return err
	}
	for i := range items {
		item := &items[i]
		action := StockActionReturn
		if wasteItemIDs[item.ID] {
			action = StockActionWaste
		}
		if err := r.insertCancellation(bill.ID, item, item.Quantity, action, reason, userID, now); err != nil {
			return err
		}
		if err := r.restoreItemStock(item, item.Quantity, action, reason, bill.ID, userID, "Batal"); err != nil {
			return err
		}
	}
	return r.db.Exec(markCancelledQuery, userID, now, reason, now, bill.ID).Error
}

// PayOpen menyelesaikan pesanan terbuka. ownerID = pemilik kas penampung (user_id transaksi), settledBy =
// pelaku penyelesaian, txDate = tanggal penjualan (saat dibayar, atau tanggal pesanan tertunda, keputusan
// 106, 113). Utang membuat piutang seperti Toko.
func (r *transactionRepo) PayOpen(billID int, req *dto.OpenPayRequest, ownerID, settledBy int, txDate, now time.Time) error {
	if err := r.db.Exec(payOpenQuery, ownerID, settledBy, req.ShiftID, txDate, req.Subtotal, req.Discount, req.Tax, req.TotalAmount,
		req.PaymentMethod, req.PaymentAmount, req.ChangeAmount, req.CustomerID, req.IsCredit, now, billID).Error; err != nil {
		return err
	}
	if req.IsCredit && req.CustomerID != nil && req.TotalAmount > 0 {
		return r.db.Exec(createReceivableQuery, billID, *req.CustomerID, req.TotalAmount, req.TotalAmount).Error
	}
	return nil
}
