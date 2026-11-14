package repo

import (
	dto "coffee_api/domain/report/dto"
)

// Laporan Coffee (tahap 4.10, keputusan 114-117). Rentang tanggal dalam format NormalizeDateRange
// ("YYYY-MM-DD 00:00:00" s.d. "YYYY-MM-DD 23:59:59").

const (
	// Harga beli rata-rata tertimbang per bahan dari seluruh pembelian yang tidak di-void (keputusan 114).
	ingredientAvgPriceSub = `SELECT pi.ingredient_id, SUM(pi.quantity * pi.purchase_price) / NULLIF(SUM(pi.quantity), 0) AS avg_price
		FROM purchase_items pi JOIN purchases p ON p.id = pi.purchase_id
		WHERE pi.item_type = 'ingredient' AND p.status <> 'void' GROUP BY pi.ingredient_id`

	salesBreakdownQuery = `SELECT COALESCE(ti.menu_item_id, 0) AS menu_item_id, ti.product_name AS menu_name,
		COALESCE(mc.name, '') AS category_name, COALESCE(ti.variant_name, '') AS variant_name,
		SUM(ti.quantity) AS qty,
		COALESCE(SUM(ti.subtotal * (t.total_amount - t.tax) / NULLIF(t.subtotal, 0)), 0) AS revenue
		FROM transaction_items ti
		JOIN transactions t ON t.id = ti.transaction_id
		LEFT JOIN menu_items mi ON mi.id = ti.menu_item_id
		LEFT JOIN menu_categories mc ON mc.id = mi.category_id
		WHERE t.status = 'completed' AND t.transaction_date BETWEEN ? AND ?
		GROUP BY ti.menu_item_id, ti.product_name, mc.name, ti.variant_name
		ORDER BY revenue DESC`
	salesByPaymentQuery = `SELECT payment_method AS ` + "`key`" + `, COUNT(*) AS count, COALESCE(SUM(total_amount), 0) AS total
		FROM transactions WHERE status = 'completed' AND transaction_date BETWEEN ? AND ?
		GROUP BY payment_method ORDER BY total DESC`
	salesByOrderTypeQuery = `SELECT COALESCE(order_type, '') AS ` + "`key`" + `, COUNT(*) AS count, COALESCE(SUM(total_amount), 0) AS total
		FROM transactions WHERE status = 'completed' AND transaction_date BETWEEN ? AND ?
		GROUP BY order_type ORDER BY total DESC`

	// Per bahan: pembelian bersih (purchase - purchase_void), terpakai dan terbuang yang tidak dibatalkan,
	// nilai pembelian dari nota pada tanggal pembelian, stok saat ini, harga rata-rata.
	ingredientReportQuery = `SELECT i.id AS ingredient_id, i.name, COALESCE(u.abbreviation, '') AS unit_abbr,
		COALESCE(SUM(CASE WHEN m.type IN ('purchase', 'purchase_void') THEN m.qty_change END), 0) AS purchased_qty,
		COALESCE(pv.purchase_value, 0) AS purchase_value,
		COALESCE(SUM(CASE WHEN m.type = 'usage' AND m.cancelled_at IS NULL THEN -m.qty_change END), 0) AS used_qty,
		COALESCE(SUM(CASE WHEN m.type = 'waste' AND m.cancelled_at IS NULL THEN -m.qty_change END), 0) AS wasted_qty,
		i.stock, COALESCE(ap.avg_price, 0) AS avg_price
		FROM ingredients i
		LEFT JOIN units u ON u.id = i.unit_id
		LEFT JOIN ingredient_mutations m ON m.ingredient_id = i.id AND m.created_at BETWEEN ? AND ?
		LEFT JOIN (` + ingredientAvgPriceSub + `) ap ON ap.ingredient_id = i.id
		LEFT JOIN (SELECT pi.ingredient_id, SUM(pi.subtotal) AS purchase_value
			FROM purchase_items pi JOIN purchases p ON p.id = pi.purchase_id
			WHERE pi.item_type = 'ingredient' AND p.status <> 'void' AND p.purchase_date BETWEEN ? AND ?
			GROUP BY pi.ingredient_id) pv ON pv.ingredient_id = i.id
		GROUP BY i.id, i.name, u.abbreviation, pv.purchase_value, i.stock, ap.avg_price
		ORDER BY i.name`

	// Barang terbuang: bahan (mutasi waste yang tidak dibatalkan) dan barang jadi (mutasi stok waste dari
	// void/batal/hapus item, keputusan 100, 105). Nilai bahan = jumlah x harga rata-rata; barang jadi =
	// jumlah satuan dasar x harga beli kemasan dasar.
	wasteListQuery = `SELECT * FROM (
		SELECT DATE_FORMAT(m.created_at, '%Y-%m-%d %H:%i:%s') AS created_at, 'ingredient' AS kind, i.name AS item_name,
			-m.qty_change AS qty, COALESCE(u.abbreviation, '') AS unit, -m.qty_change * COALESCE(ap.avg_price, 0) AS value,
			COALESCE(m.notes, '') AS reason, COALESCE(us.full_name, '') AS user_name
		FROM ingredient_mutations m
		JOIN ingredients i ON i.id = m.ingredient_id
		LEFT JOIN units u ON u.id = i.unit_id
		LEFT JOIN users us ON us.id = m.created_by
		LEFT JOIN (` + ingredientAvgPriceSub + `) ap ON ap.ingredient_id = m.ingredient_id
		WHERE m.type = 'waste' AND m.cancelled_at IS NULL AND m.created_at BETWEEN ? AND ?
		UNION ALL
		SELECT DATE_FORMAT(sm.created_at, '%Y-%m-%d %H:%i:%s'), 'product', p.name,
			sm.quantity, COALESCE(un.name, ''), sm.quantity * COALESCE(pp.purchase_price, p.purchase_price, 0),
			COALESCE(sm.notes, ''), COALESCE(us.full_name, '')
		FROM stock_mutations sm
		JOIN products p ON p.id = sm.product_id
		LEFT JOIN product_packages pp ON pp.product_id = p.id AND pp.is_default = 1
		LEFT JOIN units un ON un.id = pp.unit_id
		LEFT JOIN users us ON us.id = sm.user_id
		WHERE sm.mutation_type = 'waste' AND sm.created_at BETWEEN ? AND ?
	) w ORDER BY w.created_at DESC`

	// Biaya bahan periode (keputusan 114): (terpakai + terbuang, tidak dibatalkan) x harga rata-rata.
	ingredientCostQuery = `SELECT COALESCE(SUM(-m.qty_change * COALESCE(ap.avg_price, 0)), 0)
		FROM ingredient_mutations m
		LEFT JOIN (` + ingredientAvgPriceSub + `) ap ON ap.ingredient_id = m.ingredient_id
		WHERE m.type IN ('usage', 'waste') AND m.cancelled_at IS NULL AND m.created_at BETWEEN ? AND ?`
)

type SalesBreakdownRow struct {
	MenuItemID   int
	MenuName     string
	CategoryName string
	VariantName  string
	Qty          float64
	Revenue      float64
}

func (r *reportRepo) GetSalesBreakdownRows(params dto.FilterParams) ([]SalesBreakdownRow, error) {
	var rows []SalesBreakdownRow
	err := r.db.Raw(salesBreakdownQuery, params.DateFrom, params.DateTo).Scan(&rows).Error
	return rows, err
}

func (r *reportRepo) GetSalesByPayment(params dto.FilterParams) ([]dto.SalesGroupRow, error) {
	rows := []dto.SalesGroupRow{}
	err := r.db.Raw(salesByPaymentQuery, params.DateFrom, params.DateTo).Scan(&rows).Error
	return rows, err
}

func (r *reportRepo) GetSalesByOrderType(params dto.FilterParams) ([]dto.SalesGroupRow, error) {
	rows := []dto.SalesGroupRow{}
	err := r.db.Raw(salesByOrderTypeQuery, params.DateFrom, params.DateTo).Scan(&rows).Error
	return rows, err
}

func (r *reportRepo) GetIngredientReport(params dto.FilterParams) ([]dto.IngredientReportRow, error) {
	rows := []dto.IngredientReportRow{}
	dateFrom, dateTo := params.DateFrom[:10], params.DateTo[:10]
	err := r.db.Raw(ingredientReportQuery, params.DateFrom, params.DateTo, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *reportRepo) GetWasteList(params dto.FilterParams) ([]dto.WasteRow, error) {
	rows := []dto.WasteRow{}
	err := r.db.Raw(wasteListQuery, params.DateFrom, params.DateTo, params.DateFrom, params.DateTo).Scan(&rows).Error
	return rows, err
}

func (r *reportRepo) GetIngredientCost(params dto.FilterParams) (float64, error) {
	var total float64
	err := r.db.Raw(ingredientCostQuery, params.DateFrom, params.DateTo).Scan(&total).Error
	return total, err
}
