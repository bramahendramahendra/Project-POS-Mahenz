package repo_product

import "strings"

const (
	checkBarcodeExistsQuery = `SELECT id FROM products WHERE barcode = ? AND id != ? LIMIT 1`
	checkSkuExistsQuery     = `SELECT id FROM products WHERE sku = ? AND id != ? LIMIT 1`
	countSkuByCategoryQuery = `SELECT COUNT(*) FROM products WHERE category_id = ?`
)

func (r *productRepo) CheckBarcodeExists(barcode string, excludeID int) (bool, error) {
	if strings.TrimSpace(barcode) == "" {
		return false, nil
	}
	var id int
	result := r.db.Raw(checkBarcodeExistsQuery, barcode, excludeID).Scan(&id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *productRepo) CheckSkuExists(sku string, excludeID int) (bool, error) {
	if strings.TrimSpace(sku) == "" {
		return false, nil
	}
	var id int
	result := r.db.Raw(checkSkuExistsQuery, sku, excludeID).Scan(&id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *productRepo) CountSkuByCategory(categoryID int) (int, error) {
	var count int
	if err := r.db.Raw(countSkuByCategoryQuery, categoryID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
