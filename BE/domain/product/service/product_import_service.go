package service

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	dto "pos_api/domain/product/dto"
	"pos_api/errors"
	log_helper "pos_api/helper/log"

	"github.com/xuri/excelize/v2"
)

// resolveCategoryID mencari/membuat kategori berdasarkan nama, dengan cache per-run
// agar nama kategori yang sama di banyak baris file import tidak memicu query/insert berulang.
func (s *productService) resolveCategoryID(cache map[string]int, name string) (*int, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if id, ok := cache[key]; ok {
		return &id, nil
	}

	cat, err := s.repoCategory.GetByName(name)
	if err != nil {
		return nil, err
	}

	var id int
	if cat == nil {
		newID, createErr := s.createCategoryWithCode(name, "")
		if createErr != nil {
			return nil, createErr
		}
		id = int(newID)
	} else {
		id = cat.ID
	}

	cache[key] = id
	return &id, nil
}

// normalizeBarcode mengembalikan angka desimal biasa jika Excel menyimpan
// barcode numerik dalam notasi ilmiah (mis. "8.991002100016E+12"), karena
// kolom Barcode tanpa format Text akan dibaca excelize sesuai tampilannya.
func normalizeBarcode(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" || !strings.ContainsAny(s, "eE") {
		return s
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (s *productService) ImportPreview(file *multipart.FileHeader) (data dto.ImportPreviewResponse, err error) {
	src, openErr := file.Open()
	if openErr != nil {
		return data, &errors.InternalServerError{Message: "Gagal membuka file"}
	}
	defer src.Close()

	f, readErr := excelize.OpenReader(src)
	if readErr != nil {
		return data, &errors.BadRequestError{Message: "Gagal membaca file Excel"}
	}
	defer f.Close()

	validCategories := make(map[string]bool)
	if cats, e := s.repoCategory.GetOptions(); e == nil {
		for _, c := range cats {
			validCategories[strings.ToLower(strings.TrimSpace(c.Name))] = true
		}
	}

	unitIDMap := make(map[string]int)
	if units, e := s.repoUnit.GetOptions(); e == nil {
		for _, u := range units {
			unitIDMap[strings.ToLower(strings.TrimSpace(u.Name))] = u.ID
			unitIDMap[strings.ToLower(strings.TrimSpace(u.Abbreviation))] = u.ID
		}
	}

	sheetProduk := "Produk"
	if idx, _ := f.GetSheetIndex(sheetProduk); idx == -1 {
		sheetProduk = f.GetSheetName(0)
	}
	produkRows, rowErr := f.GetRows(sheetProduk)
	if rowErr != nil || len(produkRows) < 2 {
		return dto.ImportPreviewResponse{
			Rows:   []dto.ImportPreviewRow{},
			Grosir: []dto.ImportPreviewGrosirRow{},
		}, nil
	}

	headerProduk := produkRows[0]
	colIdx := func(headers []string, name string) int {
		for i, h := range headers {
			if strings.EqualFold(strings.TrimSpace(h), name) {
				return i
			}
		}
		return -1
	}
	getCell := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	toFloat := func(s string) float64 {
		s = strings.ReplaceAll(s, ",", "")
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	toInt := func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	}

	colNo := colIdx(headerProduk, "No")
	colNama := colIdx(headerProduk, "Produk")
	colBarcode := colIdx(headerProduk, "Barcode")
	colKategori := colIdx(headerProduk, "Kategori")
	colHargaBeli := colIdx(headerProduk, "Harga Beli")
	colHargaJual := colIdx(headerProduk, "Harga Jual")
	colStok := colIdx(headerProduk, "Stok")
	colStokMin := colIdx(headerProduk, "Stok Minimum")
	colSatuan := colIdx(headerProduk, "Satuan")

	seenBarcodes := make(map[string]bool)
	validNos := make(map[int]bool)

	var previewRows []dto.ImportPreviewRow
	for _, row := range produkRows[1:] {
		no := toInt(getCell(row, colNo))
		if no <= 0 {
			continue
		}

		errs := []string{}
		warns := []string{}

		nama := getCell(row, colNama)
		barcode := normalizeBarcode(getCell(row, colBarcode))
		kategori := getCell(row, colKategori)
		satuan := getCell(row, colSatuan)
		hargaBeli := toFloat(getCell(row, colHargaBeli))
		hargaJual := toFloat(getCell(row, colHargaJual))
		stok := toFloat(getCell(row, colStok))
		stokMin := toFloat(getCell(row, colStokMin))

		if nama == "" {
			errs = append(errs, "Nama produk wajib diisi")
		}
		satuanID := unitIDMap[strings.ToLower(satuan)]
		if satuan == "" {
			errs = append(errs, "Satuan wajib diisi")
		} else if satuanID == 0 {
			errs = append(errs, fmt.Sprintf("Satuan \"%s\" tidak ditemukan di master data", satuan))
		}
		if hargaJual <= 0 {
			errs = append(errs, "Harga jual harus lebih dari 0")
		}
		if hargaJual > 0 && hargaBeli > 0 && hargaJual < hargaBeli {
			errs = append(errs, "Harga jual tidak boleh lebih rendah dari harga beli")
		}
		if stok < 0 {
			errs = append(errs, "Stok tidak boleh negatif")
		}
		if stokMin < 0 {
			errs = append(errs, "Stok minimum tidak boleh negatif")
		}

		if kategori == "" {
			warns = append(warns, "Kategori kosong — produk akan masuk tanpa kategori")
		} else if !validCategories[strings.ToLower(kategori)] {
			errs = append(errs, fmt.Sprintf("Kategori \"%s\" tidak ditemukan di master data", kategori))
		}

		if barcode == "" {
			gen, genErr := s.GenerateBarcode()
			if genErr == nil {
				barcode = gen.Barcode
				warns = append(warns, "Barcode kosong — di-generate otomatis")
			} else {
				errs = append(errs, "Gagal generate barcode")
			}
		} else {
			barcodeKey := strings.ToLower(barcode)
			if seenBarcodes[barcodeKey] {
				errs = append(errs, fmt.Sprintf("Barcode \"%s\" duplikat dalam file", barcode))
			} else {
				exists, checkErr := s.repo.CheckBarcodeExists(barcode, 0)
				if checkErr == nil && exists {
					errs = append(errs, fmt.Sprintf("Barcode \"%s\" sudah digunakan produk lain", barcode))
				}
			}
		}
		seenBarcodes[strings.ToLower(barcode)] = true

		valid := len(errs) == 0
		if valid {
			validNos[no] = true
		}

		margin := 0
		if hargaJual > 0 && hargaBeli >= 0 {
			margin = int(math.Round(((hargaJual - hargaBeli) / hargaJual) * 100))
		}

		previewRows = append(previewRows, dto.ImportPreviewRow{
			No:          no,
			Nama:        nama,
			Barcode:     barcode,
			Kategori:    kategori,
			HargaBeli:   hargaBeli,
			HargaJual:   hargaJual,
			Margin:      margin,
			Stok:        stok,
			StokMinimum: stokMin,
			Satuan:      satuan,
			SatuanID:    satuanID,
			Valid:       valid,
			Errors:      errs,
			Warnings:    warns,
		})
	}

	sheetGrosir := "Grosir"
	if idx, _ := f.GetSheetIndex(sheetGrosir); idx == -1 {
		sheetGrosir = f.GetSheetName(1)
	}
	grosirRows, _ := f.GetRows(sheetGrosir)

	var previewGrosir []dto.ImportPreviewGrosirRow
	if len(grosirRows) >= 2 {
		headerGrosir := grosirRows[0]
		gColNoProduk := colIdx(headerGrosir, "No Produk")
		gColNamaPaket := colIdx(headerGrosir, "Nama Paket")
		gColSatuan := colIdx(headerGrosir, "Satuan")
		gColKonversi := colIdx(headerGrosir, "Konversi")
		gColHargaBeli := colIdx(headerGrosir, "Harga Beli")
		gColHargaJual := colIdx(headerGrosir, "Harga Jual")

		for _, row := range grosirRows[1:] {
			noProduk := toInt(getCell(row, gColNoProduk))
			if noProduk <= 0 {
				continue
			}
			errs := []string{}
			namaPaket := getCell(row, gColNamaPaket)
			satuan := getCell(row, gColSatuan)
			konversi := toFloat(getCell(row, gColKonversi))
			hargaBeli := toFloat(getCell(row, gColHargaBeli))
			hargaJual := toFloat(getCell(row, gColHargaJual))

			if !validNos[noProduk] {
				errs = append(errs, fmt.Sprintf("No produk %d tidak ditemukan atau tidak valid di sheet Produk", noProduk))
			}
			satuanID := unitIDMap[strings.ToLower(satuan)]
			if satuan == "" {
				errs = append(errs, "Satuan wajib diisi")
			} else if satuanID == 0 {
				errs = append(errs, fmt.Sprintf("Satuan \"%s\" tidak ditemukan di master data", satuan))
			}
			if konversi <= 0 {
				errs = append(errs, "Konversi harus lebih dari 0")
			}
			if hargaJual <= 0 {
				errs = append(errs, "Harga jual harus lebih dari 0")
			}

			previewGrosir = append(previewGrosir, dto.ImportPreviewGrosirRow{
				NoProduk:  noProduk,
				NamaPaket: namaPaket,
				Satuan:    satuan,
				SatuanID:  satuanID,
				Konversi:  konversi,
				HargaBeli: hargaBeli,
				HargaJual: hargaJual,
				Valid:     len(errs) == 0,
				Errors:    errs,
			})
		}
	}

	if previewRows == nil {
		previewRows = []dto.ImportPreviewRow{}
	}
	if previewGrosir == nil {
		previewGrosir = []dto.ImportPreviewGrosirRow{}
	}

	return dto.ImportPreviewResponse{
		Rows:   previewRows,
		Grosir: previewGrosir,
	}, nil
}

func (s *productService) ImportBulk(bulkReq dto.BulkImportRequest) (data dto.BulkImportResult, err error) {
	data.Failed = []dto.BulkImportFailed{}

	noToProductID := make(map[int]int)
	defaultPackages := make(map[int]dto.PackageRequest)
	categoryCache := make(map[string]int)

	for i, row := range bulkReq.Rows {
		rowNum := i + 2

		addFailed := func(alasan string) {
			data.Failed = append(data.Failed, dto.BulkImportFailed{
				Baris:  rowNum,
				Data:   row,
				Alasan: alasan,
			})
		}

		if strings.TrimSpace(row.Nama) == "" {
			addFailed("Nama produk kosong")
			continue
		}

		satuanKey := strings.ToLower(strings.TrimSpace(row.Satuan))
		if satuanKey == "" {
			addFailed("Satuan kosong")
			continue
		}
		resolvedUnitID := row.SatuanID
		if resolvedUnitID == 0 {
			addFailed(fmt.Sprintf("Satuan \"%s\" tidak ditemukan di master data", row.Satuan))
			continue
		}

		req := &dto.CreateRequest{
			Barcode:       strings.TrimSpace(row.Barcode),
			Name:          strings.TrimSpace(row.Nama),
			PurchasePrice: row.HargaBeli,
			SellingPrice:  row.HargaJual,
			Stock:         row.Stok,
			MinStock:      row.StokMinimum,
			UnitID:        resolvedUnitID,
		}

		kategori := strings.TrimSpace(row.Kategori)
		if kategori != "" {
			categoryID, catErr := s.resolveCategoryID(categoryCache, kategori)
			if catErr != nil {
				addFailed(fmt.Sprintf("Gagal memproses kategori: %s", kategori))
				continue
			}
			req.CategoryID = categoryID
		}

		if req.Barcode == "" {
			gen, genErr := s.GenerateBarcode()
			if genErr != nil {
				addFailed("Gagal generate barcode")
				continue
			}
			req.Barcode = gen.Barcode
		} else {
			exists, checkErr := s.repo.CheckBarcodeExists(req.Barcode, 0)
			if checkErr != nil {
				addFailed("Gagal memeriksa barcode")
				continue
			}
			if exists {
				addFailed(fmt.Sprintf("Barcode sudah digunakan: %s", req.Barcode))
				continue
			}
		}

		if req.CategoryID != nil {
			skuResp, skuErr := s.GenerateSku(*req.CategoryID)
			if skuErr == nil {
				req.SKU = skuResp.SKU
			}
		}

		productID, createErr := s.repo.Create(req)
		if createErr != nil {
			addFailed("Gagal menyimpan produk")
			continue
		}

		if row.No > 0 {
			noToProductID[row.No] = int(productID)
		}

		defaultPackages[int(productID)] = dto.PackageRequest{
			UnitID:        resolvedUnitID,
			ConversionQty: 1,
			PurchasePrice: row.HargaBeli,
			SellingPrice:  row.HargaJual,
			IsDefault:     true,
		}

		data.Success++
	}

	grosirByProduct := make(map[int][]dto.PackageRequest)
	for _, g := range bulkReq.Grosir {
		productID, ok := noToProductID[g.NoProduk]
		if !ok || g.SatuanID == 0 {
			continue
		}
		grosirByProduct[productID] = append(grosirByProduct[productID], dto.PackageRequest{
			UnitID:        g.SatuanID,
			PackageName:   strings.TrimSpace(g.NamaPaket),
			ConversionQty: g.Konversi,
			PurchasePrice: g.HargaBeli,
			SellingPrice:  g.HargaJual,
			IsDefault:     false,
		})
	}

	for productID, defaultPkg := range defaultPackages {
		allPkgs := []dto.PackageRequest{defaultPkg}
		if grosirPkgs, ok := grosirByProduct[productID]; ok {
			allPkgs = append(allPkgs, grosirPkgs...)
		}
		if err := s.repo.SavePackages(productID, allPkgs); err != nil {
			entry := log_helper.FromBackground("ImportBulk", "product_import",
				fmt.Sprintf("Gagal menyimpan paket produk ID %d: %v", productID, err))
			log_helper.LogError(entry)
		}
	}

	return data, nil
}
