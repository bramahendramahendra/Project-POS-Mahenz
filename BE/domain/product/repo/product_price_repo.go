package repo

import (
	model "pos_api/domain/product/model"
)

const (
	getProductPricesQuery = `SELECT id, product_id, tier_name, min_qty, price FROM product_prices WHERE product_id = ? ORDER BY min_qty`
)

func (r *productRepo) GetPricesByProduct(productID int) ([]*model.ProductPrice, error) {
	var dataDB []*model.ProductPrice
	err := r.db.Raw(getProductPricesQuery, productID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}
