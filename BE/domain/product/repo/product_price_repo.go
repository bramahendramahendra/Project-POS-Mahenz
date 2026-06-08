package repo_product

import (
	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"

	"gorm.io/gorm"
)

const (
	getProductPricesQuery    = `SELECT id, product_id, tier_name, min_qty, price FROM product_prices WHERE product_id = ? ORDER BY min_qty`
	deleteProductPricesQuery = `DELETE FROM product_prices WHERE product_id = ?`
	insertProductPriceQuery  = `INSERT INTO product_prices (product_id, tier_name, min_qty, price) VALUES (?, ?, ?, ?)`
)

func (r *productPriceRepo) GetByProduct(productID int) ([]*model.ProductPrice, error) {
	var dataDB []*model.ProductPrice
	err := r.db.Raw(getProductPricesQuery, productID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productPriceRepo) Save(productID int, prices []dto.ProductPriceRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(deleteProductPricesQuery, productID).Error
		if err != nil {
			return err
		}
		for _, p := range prices {
			err = tx.Exec(insertProductPriceQuery, productID, p.TierName, p.MinQty, p.Price).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
