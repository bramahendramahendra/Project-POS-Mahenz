package repo

const (
	checkBarcodeExistsQuery = `SELECT id FROM products WHERE barcode = ? AND id != ? LIMIT 1`
	checkSkuExistsQuery     = `SELECT id FROM products WHERE sku = ? AND id != ? LIMIT 1`
	countSkuByCategoryQuery = `SELECT COUNT(*) FROM products WHERE category_id = ?`
)

func (r *productRepo) CheckBarcodeExists(barcode string, excludeID int) (bool, error) {
	var id int
	err := r.db.Raw(checkBarcodeExistsQuery, barcode, excludeID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *productRepo) CheckSkuExists(sku string, excludeID int) (bool, error) {
	var id int
	err := r.db.Raw(checkSkuExistsQuery, sku, excludeID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

func (r *productRepo) CountSkuByCategory(categoryID int) (int, error) {
	var count int
	err := r.db.Raw(countSkuByCategoryQuery, categoryID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
