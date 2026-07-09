package repo

import (
	request_helper "pos_api/helper/request"
	dto "pos_api/domain/stock_mutation/dto"
)

const (
	getAllMutationsQuery = `
		SELECT sm.id, sm.product_id, p.name as product_name, sm.mutation_type,
		       sm.quantity, sm.stock_before, sm.stock_after, sm.reference_type,
		       sm.reference_id, sm.notes, COALESCE(u.full_name, '') as user_name, sm.created_at
		FROM stock_mutations sm
		LEFT JOIN products p ON sm.product_id = p.id
		LEFT JOIN users u ON sm.user_id = u.id
		WHERE 1=1`

	countMutationsBase = `SELECT COUNT(*) FROM stock_mutations sm WHERE 1=1`

	getMutationsByProductQuery = `
		SELECT sm.id, sm.mutation_type, sm.quantity, sm.stock_before, sm.stock_after,
		       sm.reference_type, sm.reference_id, sm.notes, COALESCE(u.full_name, '') as user_name, sm.created_at
		FROM stock_mutations sm
		LEFT JOIN users u ON sm.user_id = u.id
		WHERE sm.product_id = ?
		ORDER BY sm.created_at DESC`
)

func (r *stockMutationRepo) GetAll(req *dto.GetAllRequest) ([]*dto.StockMutationResponse, int64, error) {
	var args []any
	conditions := ""

	if req.ProductID != nil {
		conditions += " AND sm.product_id = ?"
		args = append(args, *req.ProductID)
	}
	if req.MutationType != "" {
		conditions += " AND sm.mutation_type = ?"
		args = append(args, req.MutationType)
	}
	if req.ReferenceType != "" {
		conditions += " AND sm.reference_type = ?"
		args = append(args, req.ReferenceType)
	}
	if req.DateFrom != "" {
		conditions += " AND DATE(sm.created_at) >= ?"
		args = append(args, req.DateFrom)
	}
	if req.DateTo != "" {
		conditions += " AND DATE(sm.created_at) <= ?"
		args = append(args, req.DateTo)
	}

	var total int64
	if err := r.db.Raw(countMutationsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

	query := getAllMutationsQuery + conditions + " ORDER BY sm.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*dto.StockMutationResponse
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *stockMutationRepo) GetByProduct(productID int) ([]*dto.StockMutationByProductResponse, error) {
	var dataDB []*dto.StockMutationByProductResponse
	if err := r.db.Raw(getMutationsByProductQuery, productID).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	return dataDB, nil
}
