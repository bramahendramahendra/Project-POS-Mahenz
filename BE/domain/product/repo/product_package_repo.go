package repo

import (
	dto_product "pos_api/domain/product/dto"
	model_product "pos_api/domain/product/model"
	time_helper "pos_api/helper/time"
)

const (
	getProductPackagesQuery = `
		SELECT pp.id, pp.product_id, pp.unit_id, COALESCE(u.name, '') AS unit_name, COALESCE(u.abbreviation, '') AS abbreviation,
		       COALESCE(pp.package_name, '') AS package_name, pp.ref_package_id, pp.qty, pp.ref_qty, pp.purchase_price, pp.selling_price, pp.is_default
		FROM product_packages pp
		JOIN units u ON u.id = pp.unit_id
		WHERE pp.product_id = ?`
	insertProductPackageQuery       = `INSERT INTO product_packages (product_id, unit_id, package_name, ref_package_id, qty, ref_qty, purchase_price, selling_price, is_default) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)`
	updateProductPackageQuery       = `UPDATE product_packages SET unit_id = ?, package_name = ?, ref_package_id = ?, qty = ?, ref_qty = ?, purchase_price = ?, selling_price = ?, updated_at = ? WHERE id = ? AND product_id = ? AND is_default = 0`
	deleteProductPackageByIDQuery   = `DELETE FROM product_packages WHERE id = ? AND product_id = ? AND is_default = 0`
	countPackagesReferencingIDQuery = `SELECT COUNT(*) FROM product_packages WHERE ref_package_id = ?`
)

func (r *productRepo) GetPackagesByProduct(productID int) ([]*model_product.ProductPackage, error) {
	var dataDB []*model_product.ProductPackage
	err := r.db.Raw(getProductPackagesQuery, productID).Scan(&dataDB).Error
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

func (r *productRepo) CreatePackage(productID int, req *dto_product.CreatePackageRequest) (int64, error) {
	var pkgName *string
	if req.PackageName != "" {
		pkgName = &req.PackageName
	}

	err := r.db.Exec(insertProductPackageQuery,
		productID, req.UnitID, pkgName, req.RefPackageID, req.Qty, req.RefQty, req.PurchasePrice, req.SellingPrice,
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

func (r *productRepo) UpdatePackage(id, productID int, req *dto_product.UpdatePackageRequest) error {
	var pkgName *string
	if req.PackageName != "" {
		pkgName = &req.PackageName
	}

	return r.db.Exec(updateProductPackageQuery,
		req.UnitID, pkgName, req.RefPackageID, req.Qty, req.RefQty, req.PurchasePrice, req.SellingPrice, time_helper.GetTimeNow(), id, productID,
	).Error
}

func (r *productRepo) CountPackagesReferencing(packageID int) (int, error) {
	var count int
	err := r.db.Raw(countPackagesReferencingIDQuery, packageID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepo) DeletePackage(id, productID int) error {
	return r.db.Exec(deleteProductPackageByIDQuery, id, productID).Error
}
