package repo_product

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"

	"gorm.io/gorm"
)

const (
	getProductPackagesQuery       = `SELECT pp.id, pp.product_id, pp.unit_id, COALESCE(u.name, '') AS unit_name, COALESCE(u.abbreviation, '') AS abbreviation, COALESCE(pp.package_name, '') AS package_name, pp.conversion_qty, pp.purchase_price, pp.selling_price, pp.is_default FROM product_packages pp JOIN units u ON u.id = pp.unit_id WHERE pp.product_id = ?`
	deleteProductPackagesQuery    = `DELETE FROM product_packages WHERE product_id = ?`
	insertProductPackageQuery     = `INSERT INTO product_packages (product_id, unit_id, package_name, conversion_qty, purchase_price, selling_price, is_default) VALUES (?, ?, ?, ?, ?, ?, ?)`
	deleteProductPackageByIDQuery = `DELETE FROM product_packages WHERE id = ? AND product_id = ?`
)

func (r *productPackageRepo) GetByProduct(productID int) ([]*model_product.ProductPackage, error) {
	var dataDB []*model_product.ProductPackage
	err := r.db.Raw(getProductPackagesQuery, productID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productPackageRepo) Save(productID int, packages []dto_product.ProductPackageRequest) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(deleteProductPackagesQuery, productID).Error
		if err != nil {
			return err
		}
		for _, p := range packages {
			var pkgName *string
			if p.PackageName != "" {
				pkgName = &p.PackageName
			}
			err = tx.Exec(insertProductPackageQuery,
				productID, p.UnitID, pkgName, p.ConversionQty, p.PurchasePrice, p.SellingPrice, p.IsDefault,
			).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *productPackageRepo) DeleteOne(id, productID int) error {
	err := r.db.Exec(deleteProductPackageByIDQuery, id, productID).Error
	return err
}
