package repo_product

import (
	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"
	request_helper "pos_api/helper/request"
)

const (
	getAllProductsBase = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.stock, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE 1=1`

	getProductByIDQuery = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.stock, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.id = ? LIMIT 1`

	getProductByBarcodeQuery = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.stock, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.barcode = ? LIMIT 1`

	searchProductsQuery = `
		SELECT p.id, p.barcode, p.name, p.selling_price, p.stock, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1 AND (p.name LIKE ? OR p.barcode LIKE ?)`

	getLowStockQuery = `
		SELECT p.id, p.name, p.stock, p.min_stock, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.stock <= p.min_stock AND p.is_active = 1`

	getProductOptionsQuery      = `SELECT id, name FROM products WHERE is_active = 1 ORDER BY name`
	checkProductUsedQuery       = `SELECT COUNT(*) FROM transaction_items WHERE product_id = ?`
	checkProductPurchasedQuery  = `SELECT COUNT(*) FROM purchase_items WHERE product_id = ?`
	createProductQuery          = `INSERT INTO products (barcode, sku, name, category_id, purchase_price, selling_price, stock, min_stock, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getLastProductInsertIDQuery = `SELECT LAST_INSERT_ID()`
	updateProductQuery          = `UPDATE products SET barcode=?, sku=?, name=?, category_id=?, purchase_price=?, selling_price=?, stock=?, min_stock=?, unit_id=?, updated_at=NOW() WHERE id=?`
	deleteProductQuery          = `DELETE FROM products WHERE id = ?`
	toggleProductStatusQuery    = `UPDATE products SET is_active = NOT is_active, updated_at = NOW() WHERE id = ?`
	updateProductStockQuery     = `UPDATE products SET stock = stock + ?, updated_at = NOW() WHERE id = ?`
	getAllProductsDefaultOrder  = ` ORDER BY p.name ASC`
	countProductsBase           = `SELECT COUNT(*) FROM products p WHERE 1=1`
)

func (r *productRepo) GetAll(req *dto.GetAllRequest) ([]*model.Product, int64, error) {
	var args []any
	conditions := ""

	if req.Search != "" {
		search := "%" + req.Search + "%"
		conditions += ` AND (p.name LIKE ? OR p.barcode LIKE ?)`
		args = append(args, search, search)
	}
	if req.CategoryID != nil {
		conditions += ` AND p.category_id = ?`
		args = append(args, *req.CategoryID)
	}
	if req.IsActive != nil {
		conditions += ` AND p.is_active = ?`
		args = append(args, *req.IsActive)
	}
	if req.LowStock {
		conditions += ` AND p.stock <= p.min_stock`
	}

	var total int64
	if err := r.db.Raw(countProductsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSortFields := map[string]string{
		"name":           "p.name",
		"selling_price":  "p.selling_price",
		"purchase_price": "p.purchase_price",
		"stock":          "p.stock",
		"is_active":      "p.is_active",
	}
	query := getAllProductsBase + conditions
	query += request_helper.BuildOrderClause(req.SortBy, req.SortOrder, allowedSortFields, getAllProductsDefaultOrder)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dataDB []*model.Product
	if err := r.db.Raw(query, args...).Scan(&dataDB).Error; err != nil {
		return nil, 0, err
	}
	return dataDB, total, nil
}

func (r *productRepo) GetOptions() ([]*model.ProductOption, error) {
	var dataDB []*model.ProductOption
	err := r.db.Raw(getProductOptionsQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productRepo) GetByID(id int) (*model.Product, error) {
	var dataDB model.Product
	err := r.db.Raw(getProductByIDQuery, id).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *productRepo) GetByBarcode(barcode string) (*model.Product, error) {
	var dataDB model.Product
	err := r.db.Raw(getProductByBarcodeQuery, barcode).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	if dataDB.ID == 0 {
		return nil, nil
	}
	return &dataDB, nil
}

func (r *productRepo) Search(req *dto.SearchRequest) ([]*model.ProductSearchResult, error) {
	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	query := searchProductsQuery
	var args []any
	search := "%" + req.Q + "%"
	query += " LIMIT ?"
	args = append(args, search, search, limit)

	var dataDB []*model.ProductSearchResult
	err := r.db.Raw(query, args...).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productRepo) GetLowStock() ([]*model.LowStockProduct, error) {
	var dataDB []*model.LowStockProduct
	err := r.db.Raw(getLowStockQuery).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productRepo) CountTransactionItems(productID int) (int, error) {
	var count int
	err := r.db.Raw(checkProductUsedQuery, productID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepo) CountPurchaseItems(productID int) (int, error) {
	var count int
	err := r.db.Raw(checkProductPurchasedQuery, productID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepo) Create(req *dto.CreateRequest) (int64, error) {
	err := r.db.Exec(createProductQuery,
		req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
		req.SellingPrice, req.Stock, req.MinStock, req.UnitID,
	).Error
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.Raw(getLastProductInsertIDQuery).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *productRepo) Update(req *dto.UpdateRequest) error {
	err := r.db.Exec(updateProductQuery,
		req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
		req.SellingPrice, req.Stock, req.MinStock, req.UnitID, req.ID,
	).Error
	return err
}

func (r *productRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteProductQuery, req.ID).Error
	return err
}

func (r *productRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleProductStatusQuery, req.ID).Error
	return err
}

func (r *productRepo) UpdateStock(id int, delta float64) error {
	err := r.db.Exec(updateProductStockQuery, delta, id).Error
	return err
}
