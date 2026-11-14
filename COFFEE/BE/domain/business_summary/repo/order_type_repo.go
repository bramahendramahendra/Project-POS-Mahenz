package repo

import (
	"fmt"

	dto "coffee_api/domain/business_summary/dto"
)

// Coffee (keputusan 117): penjualan per jenis pesanan (makan di tempat / bawa pulang).
const orderTypesQuery = `
	SELECT COALESCE(order_type, '') as order_type, COALESCE(SUM(total_amount),0) as total, COUNT(*) as count
	FROM transactions WHERE DATE(transaction_date) BETWEEN ? AND ? AND status = 'completed'
	GROUP BY order_type ORDER BY total DESC`

func (r *businessSummaryRepo) GetOrderTypes(filter dto.DateRangeFilter) ([]dto.OrderTypeItem, error) {
	items := []dto.OrderTypeItem{}
	if err := r.db.Raw(orderTypesQuery, filter.StartDate, filter.EndDate).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("GetOrderTypes: %w", err)
	}
	return items, nil
}
