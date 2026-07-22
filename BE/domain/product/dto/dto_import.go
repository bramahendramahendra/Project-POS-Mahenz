package dto

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

// GrosirImportRow: satuan lain untuk sebuah produk dari sheet "Grosir". Qty & RefQty
// sama seperti di form manual — "Qty x Satuan = RefQty x Satuan Dasar produk". Ini
// memungkinkan import juga menambahkan satuan yang LEBIH KECIL dari satuan dasar
// (mis. Qty=12 Batang = RefQty=1 Pack), bukan cuma lebih besar seperti sebelumnya.
type GrosirImportRow struct {
	NoProduk  int     `json:"no_produk"`
	NamaPaket string  `json:"nama_paket"`
	Satuan    string  `json:"satuan"`
	SatuanID  int     `json:"satuan_id"`
	Qty       float64 `json:"qty"`
	RefQty    float64 `json:"ref_qty"`
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
	Valid       bool     `json:"valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
}

type ImportPreviewGrosirRow struct {
	NoProduk  int      `json:"no_produk"`
	NamaPaket string   `json:"nama_paket"`
	Satuan    string   `json:"satuan"`
	SatuanID  int      `json:"satuan_id"`
	Qty       float64  `json:"qty"`
	RefQty    float64  `json:"ref_qty"`
	HargaBeli float64  `json:"harga_beli"`
	HargaJual float64  `json:"harga_jual"`
	Valid     bool     `json:"valid"`
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
