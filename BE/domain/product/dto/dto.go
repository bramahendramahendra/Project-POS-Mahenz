package dto

type (
	// REQUEST
	CreateRequest struct {
		Barcode       string  `json:"barcode" validate:"required,max=100"`
		SKU           string  `json:"sku" validate:"required,max=50"`
		Name          string  `json:"name" validate:"required,max=200"`
		CategoryID    *int    `json:"category_id" validate:"required"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"required,min=0"`
		Stock         float64 `json:"stock" validate:"min=0"`
		MinStock      float64 `json:"min_stock" validate:"min=0"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		// Packages: satuan lain yang diisi bersamaan saat produk dibuat (opsional).
		// Diproses berurutan sesuai urutan array ini — lihat PackageDraftRequest.
		Packages []PackageDraftRequest `json:"packages" validate:"dive"`
	}

	// PackageDraftRequest: satu baris "satuan lain" yang diisi di form Tambah Produk,
	// sebelum produk (dan paket anchor-nya) benar-benar tersimpan. TempID adalah
	// penanda sementara dari FE (bukan ID asli product_packages) — dipakai draft lain
	// untuk saling merujuk dalam satu payload. RefTempID = 0 berarti merujuk paket
	// anchor (satuan dasar produk); selain itu wajib TempID draft LAIN yang muncul
	// LEBIH DULU di array Packages (diproses berurutan, tidak boleh maju/referensi ke depan).
	PackageDraftRequest struct {
		TempID        int     `json:"temp_id" validate:"required,min=1"`
		UnitID        int     `json:"unit_id" validate:"required,min=1"`
		PackageName   string  `json:"package_name" validate:"max=100"`
		RefTempID     int     `json:"ref_temp_id"`
		Qty           float64 `json:"qty" validate:"required,min=0.001"`
		RefQty        float64 `json:"ref_qty" validate:"required,min=0.001"`
		PurchasePrice float64 `json:"purchase_price" validate:"min=0"`
		SellingPrice  float64 `json:"selling_price" validate:"min=0"`
	}

	// RESPONSE
	GetLowStockResponse struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Stock    float64 `json:"stock"`
		MinStock float64 `json:"min_stock"`
		UnitName string  `json:"unit_name"`
	}

	GetUnitInfoResponse struct {
		Name         string
		Abbreviation string
	}
)
