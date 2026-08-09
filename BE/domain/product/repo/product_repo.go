package repo

import (
	"fmt"

	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"
	request_helper "pos_api/helper/request"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

const (
	getAllProductsBase = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.stock, p.reserved_qty, p.min_stock,
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
		       p.purchase_price, p.selling_price, p.stock, p.reserved_qty, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.id = ? LIMIT 1`

	getProductByBarcodeQuery = `
		SELECT p.id, p.barcode, COALESCE(p.sku, '') as sku, p.name, p.category_id, COALESCE(c.name, '') as category_name,
		       p.purchase_price, p.selling_price, p.stock, p.reserved_qty, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name, COALESCE(u.abbreviation, '') as unit_abbreviation,
		       p.is_active, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM product_packages pp WHERE pp.product_id = p.id AND pp.is_default = 0) AS extra_packages,
		       (SELECT COUNT(*) FROM product_prices pr WHERE pr.product_id = p.id) AS price_tiers_count
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.barcode = ? LIMIT 1`

	searchProductsQuery = `
		SELECT p.id, p.barcode, p.name, p.selling_price, (p.stock - p.reserved_qty) as stock, p.min_stock,
		       COALESCE(p.unit_id, 0) as unit_id, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE p.is_active = 1 AND (p.name LIKE ? OR p.barcode LIKE ?)`

	getLowStockQuery = `
		SELECT p.id, p.name, (p.stock - p.reserved_qty) as stock, p.min_stock, COALESCE(u.name, '') as unit_name
		FROM products p
		LEFT JOIN units u ON u.id = p.unit_id
		WHERE (p.stock - p.reserved_qty) <= p.min_stock AND p.is_active = 1`

	getProductOptionsQuery      = `SELECT id, name FROM products WHERE is_active = 1 ORDER BY name`
	checkProductUsedQuery       = `SELECT COUNT(*) FROM transaction_items WHERE product_id = ?`
	checkProductPurchasedQuery  = `SELECT COUNT(*) FROM purchase_items WHERE product_id = ?`
	createProductQuery          = `INSERT INTO products (barcode, sku, name, category_id, purchase_price, selling_price, stock, min_stock, unit_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	insertAnchorPackageOnCreate = `INSERT INTO product_packages (product_id, unit_id, purchase_price, selling_price, is_default) VALUES (?, ?, ?, ?, 1)`
	getLastProductInsertIDQuery = `SELECT LAST_INSERT_ID()`
	updateProductQuery          = `UPDATE products SET barcode=?, sku=?, name=?, category_id=?, purchase_price=?, selling_price=?, stock=?, min_stock=?, updated_at=? WHERE id=?`
	updateAnchorPackagePrice    = `UPDATE product_packages SET purchase_price=?, selling_price=?, updated_at=? WHERE product_id=? AND is_default=1`
	deleteProductQuery          = `DELETE FROM products WHERE id = ?`
	toggleProductStatusQuery    = `UPDATE products SET is_active = NOT is_active, updated_at = ? WHERE id = ?`
	updateProductStockQuery     = `UPDATE products SET stock = stock + ?, updated_at = ? WHERE id = ?`
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
		conditions += ` AND (p.stock - p.reserved_qty) <= p.min_stock`
	}

	var total int64
	if err := r.db.Raw(countProductsBase+conditions, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := request_helper.NormalizePagination(req.Page, req.Limit, 10, 100)

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

// Create menyimpan produk baru sekaligus paket anchor-nya (satuan pencatatan stok,
// permanen sejak dibuat) dan satuan lain (req.Packages, opsional) dalam SATU transaksi —
// supaya user bisa isi semua satuan produk langsung di form Tambah Produk, satu kali
// simpan, tanpa produk pernah ada dalam keadaan "setengah jadi" kalau ada yang gagal.
func (r *productRepo) Create(req *dto.CreateRequest) (int64, error) {
	var id int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(createProductQuery,
			req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
			req.SellingPrice, req.Stock, req.MinStock, req.UnitID,
		).Error; err != nil {
			return err
		}

		if err := tx.Raw(getLastProductInsertIDQuery).Scan(&id).Error; err != nil {
			return err
		}

		if err := tx.Exec(insertAnchorPackageOnCreate, id, req.UnitID, req.PurchasePrice, req.SellingPrice).Error; err != nil {
			return err
		}
		var anchorID int64
		if err := tx.Raw(getLastProductInsertIDQuery).Scan(&anchorID).Error; err != nil {
			return err
		}

		// tempToReal: peta penanda sementara dari FE (temp_id) ke ID asli product_packages
		// yang baru dibuat. 0 selalu berarti paket anchor.
		tempToReal := map[int]int64{0: anchorID}
		for _, p := range req.Packages {
			refID, ok := tempToReal[p.RefTempID]
			if !ok {
				return fmt.Errorf("paket dengan temp_id %d merujuk paket yang belum dibuat (ref_temp_id %d)", p.TempID, p.RefTempID)
			}

			var pkgName *string
			if p.PackageName != "" {
				pkgName = &p.PackageName
			}
			if err := tx.Exec(insertProductPackageQuery,
				id, p.UnitID, pkgName, refID, p.Qty, p.RefQty, p.PurchasePrice, p.SellingPrice,
			).Error; err != nil {
				return err
			}
			var newID int64
			if err := tx.Raw(getLastProductInsertIDQuery).Scan(&newID).Error; err != nil {
				return err
			}
			tempToReal[p.TempID] = newID
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *productRepo) Update(req *dto.UpdateRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time_helper.GetTimeNow()
		if err := tx.Exec(updateProductQuery,
			req.Barcode, req.SKU, req.Name, req.CategoryID, req.PurchasePrice,
			req.SellingPrice, req.Stock, req.MinStock, now, req.ID,
		).Error; err != nil {
			return err
		}

		return tx.Exec(updateAnchorPackagePrice, req.PurchasePrice, req.SellingPrice, now, req.ID).Error
	})
}

func (r *productRepo) Delete(req *dto.DeleteRequest) error {
	err := r.db.Exec(deleteProductQuery, req.ID).Error
	return err
}

func (r *productRepo) ToggleStatus(req *dto.ToggleStatusRequest) error {
	err := r.db.Exec(toggleProductStatusQuery, time_helper.GetTimeNow(), req.ID).Error
	return err
}

func (r *productRepo) UpdateStock(id int, delta float64) error {
	err := r.db.Exec(updateProductStockQuery, delta, time_helper.GetTimeNow(), id).Error
	return err
}
