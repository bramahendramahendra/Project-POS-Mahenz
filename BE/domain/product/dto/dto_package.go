package dto

type (
	// REQUEST
	PackageByProductRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// CreatePackageRequest: paket baru, wajib direferensikan ke paket lain yang sudah ada
	// (kecuali paket anchor pertama, yang dibuat otomatis saat produk dibuat — lihat
	// product_service.go Create). RefPackageID di sini selalu wajib diisi karena endpoint
	// ini cuma dipakai untuk menambah paket grosir/satuan lain, bukan paket anchor.
	CreatePackageRequest struct {
		ProductID     int     `json:"-"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		PackageName   string  `json:"package_name" validate:"max=100"`
		RefPackageID  int     `json:"ref_package_id" validate:"required,min=1"`
		Qty           float64 `json:"qty" validate:"required,min=0.001"`
		RefQty        float64 `json:"ref_qty" validate:"required,min=0.001"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"min=0"`
	}

	UpdatePackageRequest struct {
		ID            int     `json:"-"`
		ProductID     int     `json:"-"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		PackageName   string  `json:"package_name" validate:"max=100"`
		RefPackageID  int     `json:"ref_package_id" validate:"required,min=1"`
		Qty           float64 `json:"qty" validate:"required,min=0.001"`
		RefQty        float64 `json:"ref_qty" validate:"required,min=0.001"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"min=0"`
	}

	DeletePackageRequest struct {
		ID        int `uri:"id" validate:"required,min=1"`
		PackageID int `uri:"package_id" validate:"required,min=1"`
	}

	PackageUriRequest struct {
		ID        int `uri:"id" validate:"required,min=1"`
		PackageID int `uri:"package_id" validate:"required,min=1"`
	}

	// RESPONSE
	PackageResponse struct {
		ID           int      `json:"id"`
		ProductID    int      `json:"product_id"`
		UnitID       int      `json:"unit_id"`
		UnitName     string   `json:"unit_name"`
		Abbreviation string   `json:"abbreviation"`
		PackageName  string   `json:"package_name"`
		RefPackageID *int     `json:"ref_package_id"`
		Qty          float64  `json:"qty"`
		RefQty       *float64 `json:"ref_qty"`
		// ResolvedFactor: 1 paket ini = ResolvedFactor x satuan anchor produk. Dihitung
		// server-side dengan menelusuri rantai ref_package_id, supaya FE (kasir, pembelian,
		// laporan) tidak perlu tahu struktur pohonnya sama sekali — cukup pakai angka ini.
		ResolvedFactor float64 `json:"resolved_factor"`
		PurchasePrice  float64 `json:"purchase_price"`
		SellingPrice   float64 `json:"selling_price"`
		IsDefault      bool    `json:"is_default"`
	}
)
