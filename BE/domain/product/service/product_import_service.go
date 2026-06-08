package service

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	dto "pos_api/domain/product/dto"
	"pos_api/errors"

	"github.com/xuri/excelize/v2"
)

func (s *productService) ImportFromFile(file *multipart.FileHeader) (data dto.ImportResult, err error) {
	src, openErr := file.Open()
	if openErr != nil {
		return data, &errors.InternalServerError{Message: "Gagal membuka file"}
	}
	defer src.Close()

	f, readErr := excelize.OpenReader(src)
	if readErr != nil {
		return data, &errors.BadRequestError{Message: "File tidak dapat dibaca sebagai Excel"}
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return data, &errors.BadRequestError{Message: "File tidak memiliki sheet"}
	}

	rows, rowErr := f.GetRows(sheets[0])
	if rowErr != nil {
		return data, &errors.InternalServerError{Message: "Gagal membaca baris file"}
	}

	data.Errors = []dto.ImportErrorDetail{}

	if len(rows) <= 1 {
		return data, nil
	}

	for i, row := range rows[1:] {
		rowNum := i + 2

		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		name := getCol(0)
		sellingPriceStr := getCol(4)

		if name == "" || sellingPriceStr == "" {
			data.Failed++
			data.Errors = append(data.Errors, dto.ImportErrorDetail{
				Row:     rowNum,
				Message: "Kolom name dan selling_price wajib diisi",
			})
			continue
		}

		sellingPrice, parseErr := strconv.ParseFloat(sellingPriceStr, 64)
		if parseErr != nil {
			data.Failed++
			data.Errors = append(data.Errors, dto.ImportErrorDetail{
				Row:     rowNum,
				Message: fmt.Sprintf("selling_price tidak valid: %s", sellingPriceStr),
			})
			continue
		}

		req := &dto.ProductRequest{
			Barcode:      getCol(1),
			Name:         name,
			SellingPrice: sellingPrice,
		}

		if v := getCol(3); v != "" {
			if pp, e := strconv.ParseFloat(v, 64); e == nil {
				req.PurchasePrice = pp
			}
		}
		if v := getCol(5); v != "" {
			if st, e := strconv.ParseFloat(v, 64); e == nil {
				req.Stock = st
			}
		}
		if v := getCol(6); v != "" {
			if ms, e := strconv.ParseFloat(v, 64); e == nil {
				req.MinStock = ms
			}
		}

		if categoryName := getCol(2); categoryName != "" {
			cat, catErr := s.repoCategory.GetByName(categoryName)
			if catErr != nil {
				data.Failed++
				data.Errors = append(data.Errors, dto.ImportErrorDetail{
					Row:     rowNum,
					Message: fmt.Sprintf("Gagal mencari kategori: %s", categoryName),
				})
				continue
			}
			if cat == nil {
				newID, createErr := s.createCategoryWithCode(categoryName, "")
				if createErr != nil {
					data.Failed++
					data.Errors = append(data.Errors, dto.ImportErrorDetail{
						Row:     rowNum,
						Message: fmt.Sprintf("Gagal membuat kategori: %s", categoryName),
					})
					continue
				}
				id := int(newID)
				req.CategoryID = &id
			} else {
				req.CategoryID = &cat.ID
			}
		}

		if req.Barcode != "" {
			exists, checkErr := s.repo.CheckBarcodeExists(req.Barcode, 0)
			if checkErr != nil {
				data.Failed++
				data.Errors = append(data.Errors, dto.ImportErrorDetail{
					Row:     rowNum,
					Message: "Gagal memeriksa barcode",
				})
				continue
			}
			if exists {
				data.Failed++
				data.Errors = append(data.Errors, dto.ImportErrorDetail{
					Row:     rowNum,
					Message: fmt.Sprintf("Barcode sudah digunakan: %s", req.Barcode),
				})
				continue
			}
		}

		if _, createErr := s.repo.Create(req); createErr != nil {
			data.Failed++
			data.Errors = append(data.Errors, dto.ImportErrorDetail{
				Row:     rowNum,
				Message: "Gagal menyimpan produk",
			})
			continue
		}

		data.Success++
	}

	return data, nil
}

func (s *productService) ImportBulk(bulkReq dto.BulkImportRequest) (data dto.BulkImportResult, err error) {
	data.Failed = []dto.BulkImportFailed{}

	noToProductID := make(map[int]int)
	defaultPackages := make(map[int]dto.ProductPackageRequest)

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

		req := &dto.ProductRequest{
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
			cat, catErr := s.repoCategory.GetByName(kategori)
			if catErr != nil {
				addFailed(fmt.Sprintf("Gagal mencari kategori: %s", kategori))
				continue
			}
			if cat == nil {
				newID, createErr := s.createCategoryWithCode(kategori, "")
				if createErr != nil {
					addFailed(fmt.Sprintf("Gagal membuat kategori: %s", kategori))
					continue
				}
				id := int(newID)
				req.CategoryID = &id
			} else {
				req.CategoryID = &cat.ID
			}
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

		defaultPackages[int(productID)] = dto.ProductPackageRequest{
			UnitID:        resolvedUnitID,
			ConversionQty: 1,
			PurchasePrice: row.HargaBeli,
			SellingPrice:  row.HargaJual,
			IsDefault:     true,
		}

		data.Success++
	}

	grosirByProduct := make(map[int][]dto.ProductPackageRequest)
	for _, g := range bulkReq.Grosir {
		productID, ok := noToProductID[g.NoProduk]
		if !ok || g.SatuanID == 0 {
			continue
		}
		grosirByProduct[productID] = append(grosirByProduct[productID], dto.ProductPackageRequest{
			UnitID:        g.SatuanID,
			PackageName:   strings.TrimSpace(g.NamaPaket),
			ConversionQty: g.Konversi,
			PurchasePrice: g.HargaBeli,
			SellingPrice:  g.HargaJual,
			IsDefault:     false,
		})
	}

	for productID, defaultPkg := range defaultPackages {
		allPkgs := []dto.ProductPackageRequest{defaultPkg}
		if grosirPkgs, ok := grosirByProduct[productID]; ok {
			allPkgs = append(allPkgs, grosirPkgs...)
		}
		_ = s.repo.SavePackages(productID, allPkgs)
	}

	return data, nil
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

		var errs, warns []string

		nama := getCell(row, colNama)
		barcode := getCell(row, colBarcode)
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
			var errs []string
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
