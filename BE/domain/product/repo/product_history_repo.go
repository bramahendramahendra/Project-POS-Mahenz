package repo

import (
	"math"
	"time"

	dto "pos_api/domain/product/dto"
)

const (
	purchaseHistoryQuery = `
		SELECT p.purchase_date, p.purchase_code, p.invoice_number,
		       COALESCE(s.name, '') AS supplier_name,
		       pi.unit, pi.quantity, pi.purchase_price, pi.subtotal
		FROM purchase_items pi
		JOIN purchases p ON pi.purchase_id = p.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		WHERE pi.product_id = ? AND p.status = 'active'
		ORDER BY p.purchase_date DESC
		LIMIT ? OFFSET ?`

	purchaseHistoryCountQuery = `
		SELECT COUNT(*)
		FROM purchase_items pi
		JOIN purchases p ON pi.purchase_id = p.id
		WHERE pi.product_id = ? AND p.status = 'active'`

	purchaseHistorySummaryQuery = `
		SELECT COUNT(DISTINCT pi.purchase_id),
		       COALESCE(SUM(pi.quantity), 0),
		       COALESCE(SUM(pi.subtotal), 0)
		FROM purchase_items pi
		JOIN purchases p ON pi.purchase_id = p.id
		WHERE pi.product_id = ? AND p.status = 'active'`

	saleHistoryQuery = `
		SELECT t.transaction_date, t.transaction_code,
		       COALESCE(c.name, '') AS customer_name,
		       ti.unit, ti.quantity, ti.price, ti.subtotal, ti.discount_item,
		       t.status
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE ti.product_id = ?`

	saleHistoryCountQuery = `
		SELECT COUNT(*)
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE ti.product_id = ?`

	saleHistorySummaryQuery = `
		SELECT COUNT(*),
		       COALESCE(SUM(ti.quantity), 0),
		       COALESCE(SUM(ti.subtotal), 0),
		       COALESCE(MIN(t.transaction_date), NOW()),
		       COALESCE(MAX(t.transaction_date), NOW())
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		WHERE ti.product_id = ? AND t.status = 'completed'`
)

func (r *productRepo) GetPurchaseHistory(productID, page, limit int) (*dto.PurchaseHistoryResponse, error) {
	offset := (page - 1) * limit

	// Count
	var total int64
	if err := r.db.Raw(purchaseHistoryCountQuery, productID).Scan(&total).Error; err != nil {
		return nil, err
	}

	// Summary
	var totalNotes int
	var totalQty, totalValue float64
	if err := r.db.Raw(purchaseHistorySummaryQuery, productID).Row().Scan(&totalNotes, &totalQty, &totalValue); err != nil {
		return nil, err
	}
	avgPrice := float64(0)
	if totalQty > 0 {
		avgPrice = math.Round(totalValue / totalQty)
	}

	// Items
	rows, err := r.db.Raw(purchaseHistoryQuery, productID, limit, offset).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.PurchaseHistoryItem
	for rows.Next() {
		var item dto.PurchaseHistoryItem
		if err := rows.Scan(
			&item.PurchaseDate, &item.PurchaseCode, &item.InvoiceNumber,
			&item.SupplierName, &item.Unit, &item.Quantity,
			&item.PurchasePrice, &item.Subtotal,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if items == nil {
		items = []dto.PurchaseHistoryItem{}
	}

	return &dto.PurchaseHistoryResponse{
		Summary: dto.PurchaseHistorySummary{
			TotalNotes:   totalNotes,
			TotalQty:     totalQty,
			TotalValue:   totalValue,
			AveragePrice: avgPrice,
		},
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *productRepo) GetSaleHistory(productID, page, limit int, status string) (*dto.SaleHistoryResponse, error) {
	offset := (page - 1) * limit

	// Build dynamic WHERE clause for status filter
	statusFilter := ""
	args := []interface{}{productID}
	if status != "" {
		statusFilter = " AND t.status = ?"
		args = append(args, status)
	}

	// Count
	var total int64
	countQuery := saleHistoryCountQuery + statusFilter
	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	// Summary (always for completed only)
	var totalTx int
	var totalQty, totalRevenue float64
	var minDate, maxDate time.Time
	if err := r.db.Raw(saleHistorySummaryQuery, productID).Row().Scan(&totalTx, &totalQty, &totalRevenue, &minDate, &maxDate); err != nil {
		return nil, err
	}
	avgPerDay := float64(0)
	if totalTx > 0 {
		days := maxDate.Sub(minDate).Hours()/24 + 1
		if days > 0 {
			avgPerDay = math.Round(totalQty/days*10) / 10
		}
	}

	// Items
	itemQuery := saleHistoryQuery + statusFilter + " ORDER BY t.transaction_date DESC LIMIT ? OFFSET ?"
	itemArgs := append(args, limit, offset)
	rows, err := r.db.Raw(itemQuery, itemArgs...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.SaleHistoryItem
	for rows.Next() {
		var item dto.SaleHistoryItem
		if err := rows.Scan(
			&item.TransactionDate, &item.TransactionCode,
			&item.CustomerName, &item.Unit, &item.Quantity,
			&item.Price, &item.Subtotal, &item.DiscountItem,
			&item.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if items == nil {
		items = []dto.SaleHistoryItem{}
	}

	return &dto.SaleHistoryResponse{
		Summary: dto.SaleHistorySummary{
			TotalTransactions: totalTx,
			TotalQty:          totalQty,
			TotalRevenue:      totalRevenue,
			AveragePerDay:     avgPerDay,
		},
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
