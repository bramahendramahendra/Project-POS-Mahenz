package dto_product

type ImportResult struct {
	Success int                 `json:"success"`
	Failed  int                 `json:"failed"`
	Errors  []ImportErrorDetail `json:"errors"`
}

type ImportErrorDetail struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type BulkImportRequest struct {
	Rows   []BulkImportRow   `json:"rows"`
	Grosir []GrosirImportRow `json:"grosir"`
}

type BulkImportRow struct {
	No          int     `json:"no"`
	Nama        string  `json:"nama"`
	Barcode     string  `json:"barcode"`
	Kategori    string  `json:"kategori"`
	HargaBeli   float64 `json:"harga_beli"`
	HargaJual   float64 `json:"harga_jual"`
	Stok        float64 `json:"stok"`
	StokMinimum float64 `json:"stok_minimum"`
	Satuan      string  `json:"satuan"`
	SatuanID    int     `json:"satuan_id"`
}

type GrosirImportRow struct {
	NoProduk  int     `json:"no_produk"`
	NamaPaket string  `json:"nama_paket"`
	Satuan    string  `json:"satuan"`
	SatuanID  int     `json:"satuan_id"`
	Konversi  float64 `json:"konversi"`
	HargaBeli float64 `json:"harga_beli"`
	HargaJual float64 `json:"harga_jual"`
}

type ImportPreviewRow struct {
	No          int      `json:"no"`
	Nama        string   `json:"nama"`
	Barcode     string   `json:"barcode"`
	Kategori    string   `json:"kategori"`
	HargaBeli   float64  `json:"harga_beli"`
	HargaJual   float64  `json:"harga_jual"`
	Margin      int      `json:"margin"`
	Stok        float64  `json:"stok"`
	StokMinimum float64  `json:"stok_minimum"`
	Satuan      string   `json:"satuan"`
	SatuanID    int      `json:"satuan_id"`
	Valid        bool     `json:"valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
}

type ImportPreviewGrosirRow struct {
	NoProduk  int      `json:"no_produk"`
	NamaPaket string   `json:"nama_paket"`
	Satuan    string   `json:"satuan"`
	SatuanID  int      `json:"satuan_id"`
	Konversi  float64  `json:"konversi"`
	HargaBeli float64  `json:"harga_beli"`
	HargaJual float64  `json:"harga_jual"`
	Valid      bool     `json:"valid"`
	Errors    []string `json:"errors"`
}

type ImportPreviewResponse struct {
	Rows   []ImportPreviewRow       `json:"rows"`
	Grosir []ImportPreviewGrosirRow `json:"grosir"`
}

type BulkImportResult struct {
	Success int                `json:"success"`
	Failed  []BulkImportFailed `json:"failed"`
}

type BulkImportFailed struct {
	Baris  int           `json:"baris"`
	Data   BulkImportRow `json:"data"`
	Alasan string        `json:"alasan"`
}
