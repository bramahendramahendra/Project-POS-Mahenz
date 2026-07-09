package handler

import (
	"fmt"

	dto "pos_api/domain/product/dto"
	service "pos_api/domain/product/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ProductImportHandler struct {
	service service.ProductServiceInterface
}

func NewProductImportHandler(service service.ProductServiceInterface) *ProductImportHandler {
	return &ProductImportHandler{service: service}
}

// POST /products/import-preview
func (h *ProductImportHandler) ImportPreview(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(&errors.BadRequestError{Message: "File tidak ditemukan"})
		return
	}

	result, svcErr := h.service.ImportPreview(file)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Preview berhasil",
		Data:    result,
	})
}

// POST /products/import-bulk
func (h *ProductImportHandler) ImportBulk(c *gin.Context) {
	req, err := binder.BindJSON[dto.BulkImportRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	result, svcErr := h.service.ImportBulk(req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Import selesai",
		Data:    result,
	})
}

// POST /products/import-template
func (h *ProductImportHandler) DownloadImportTemplate(c *gin.Context) {
	categoryNames, err := h.service.GetCategoryNames()
	if err != nil {
		c.Error(err)
		return
	}
	unitInfos, err := h.service.GetUnitInfos()
	if err != nil {
		c.Error(err)
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	// ── Sheet 1: Produk ──────────────────────────────────────────────────────
	sheetProduk := "Produk"
	f.SetSheetName("Sheet1", sheetProduk)

	prodHeaders := []string{"No", "Produk", "Barcode", "Kategori", "Harga Beli", "Harga Jual", "Margin", "Stok", "Stok Minimum", "Satuan"}
	prodExample := []any{0, "Contoh Produk A", "", "Minuman", 5000, 7000, "", 100, 10, "pcs"}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
	})
	noteStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Italic: true, Color: "#7F7F7F", Size: 9},
	})
	currencyStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 3,
	})
	textStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 49, // "@" (Text) — cegah Excel mengubah barcode panjang jadi notasi ilmiah
	})
	marginHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#808080"}, Pattern: 1},
	})
	marginCellStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Font:      &excelize.Font{Color: "#595959"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		NumFmt:    9,
	})

	for col, h := range prodHeaders {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetProduk, cell, h)
	}
	for col, val := range prodExample {
		cell, _ := excelize.CoordinatesToCellName(col+1, 2)
		f.SetCellValue(sheetProduk, cell, val)
	}
	f.SetCellStyle(sheetProduk, "A1", "F1", headerStyle)
	f.SetCellStyle(sheetProduk, "G1", "G1", marginHeaderStyle)
	f.SetCellStyle(sheetProduk, "H1", "J1", headerStyle)
	f.SetCellStyle(sheetProduk, "C2", "C1000", textStyle)
	f.SetCellStyle(sheetProduk, "E2", "E1000", currencyStyle)
	f.SetCellStyle(sheetProduk, "F2", "F1000", currencyStyle)
	f.SetCellStyle(sheetProduk, "G2", "G1000", marginCellStyle)
	for row := 2; row <= 1000; row++ {
		cellG, _ := excelize.CoordinatesToCellName(7, row)
		f.SetCellFormula(sheetProduk, cellG, fmt.Sprintf(`IFERROR(ROUND((F%d-E%d)/F%d,2),0)`, row, row, row))
	}

	f.SetCellValue(sheetProduk, "A3", "* No: nomor unik dalam file ini (dipakai sheet Grosir). Barcode opsional (di-generate otomatis jika kosong). Kolom Margin otomatis.")
	f.SetCellStyle(sheetProduk, "A3", "J3", noteStyle)
	f.MergeCell(sheetProduk, "A3", "J3")

	prodColWidths := map[string]float64{
		"A": 6, "B": 26, "C": 18, "D": 20,
		"E": 14, "F": 14, "G": 10, "H": 10, "I": 16, "J": 14,
	}
	for col, w := range prodColWidths {
		f.SetColWidth(sheetProduk, col, col, w)
	}

	dvNo := excelize.NewDataValidation(true)
	dvNo.Sqref = "A4:A1000"
	dvNo.SetRange(1, 999999, excelize.DataValidationTypeWhole, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvNo.SetError(excelize.DataValidationErrorStyleStop, "No tidak valid", "Kolom 'No' harus diisi dengan angka bulat minimal 1.")
	f.AddDataValidation(sheetProduk, dvNo)

	dvNama := excelize.NewDataValidation(true)
	dvNama.Sqref = "B4:B1000"
	dvNama.SetRange(1, 200, excelize.DataValidationTypeTextLength, excelize.DataValidationOperatorBetween)
	dvNama.SetError(excelize.DataValidationErrorStyleStop, "Nama tidak valid", "Nama produk wajib diisi, maksimal 200 karakter.")
	f.AddDataValidation(sheetProduk, dvNama)

	dvBarcode := excelize.NewDataValidation(true)
	dvBarcode.Sqref = "C4:C1000"
	dvBarcode.SetRange(0, 100, excelize.DataValidationTypeTextLength, excelize.DataValidationOperatorBetween)
	dvBarcode.SetError(excelize.DataValidationErrorStyleStop, "Barcode tidak valid", "Barcode maksimal 100 karakter. Kosongkan jika ingin di-generate otomatis.")
	f.AddDataValidation(sheetProduk, dvBarcode)

	if len(categoryNames) > 0 {
		lastCatRow := len(categoryNames) + 1
		dvKategori := excelize.NewDataValidation(true)
		dvKategori.Sqref = "D4:D1000"
		dvKategori.Formula1 = fmt.Sprintf("Kategori!$A$2:$A$%d", lastCatRow)
		dvKategori.ShowDropDown = true
		dvKategori.SetError(excelize.DataValidationErrorStyleStop, "Kategori tidak valid", "Pilih kategori dari daftar yang tersedia di sheet Kategori.")
		f.AddDataValidation(sheetProduk, dvKategori)
	}

	dvHargaBeli := excelize.NewDataValidation(true)
	dvHargaBeli.Sqref = "E4:E1000"
	dvHargaBeli.SetRange(0, 9999999999999.99, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvHargaBeli.SetError(excelize.DataValidationErrorStyleStop, "Harga beli tidak valid", "Harga beli harus berupa angka >= 0.")
	f.AddDataValidation(sheetProduk, dvHargaBeli)

	dvHargaJual := excelize.NewDataValidation(true)
	dvHargaJual.Sqref = "F4:F1000"
	dvHargaJual.SetRange(0.01, 9999999999999.99, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvHargaJual.SetError(excelize.DataValidationErrorStyleStop, "Harga jual tidak valid", "Harga jual harus berupa angka > 0.")
	f.AddDataValidation(sheetProduk, dvHargaJual)

	dvStok := excelize.NewDataValidation(true)
	dvStok.Sqref = "H4:H1000"
	dvStok.SetRange(0, 9999999999999.999, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvStok.SetError(excelize.DataValidationErrorStyleStop, "Stok tidak valid", "Stok harus berupa angka >= 0.")
	f.AddDataValidation(sheetProduk, dvStok)

	dvStokMin := excelize.NewDataValidation(true)
	dvStokMin.Sqref = "I4:I1000"
	dvStokMin.SetRange(0, 9999999999999.999, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvStokMin.SetError(excelize.DataValidationErrorStyleStop, "Stok minimum tidak valid", "Stok minimum harus berupa angka >= 0.")
	f.AddDataValidation(sheetProduk, dvStokMin)

	if len(unitInfos) > 0 {
		lastUnitRow := len(unitInfos) + 1
		dvSatuan := excelize.NewDataValidation(true)
		dvSatuan.Sqref = "J4:J1000"
		dvSatuan.Formula1 = fmt.Sprintf("Satuan!$A$2:$A$%d", lastUnitRow)
		dvSatuan.ShowDropDown = true
		dvSatuan.SetError(excelize.DataValidationErrorStyleStop, "Satuan tidak valid", "Pilih satuan dari daftar yang tersedia di sheet Satuan.")
		f.AddDataValidation(sheetProduk, dvSatuan)
	}

	// ── Sheet 2: Grosir ──────────────────────────────────────────────────────
	sheetGrosir := "Grosir"
	f.NewSheet(sheetGrosir)

	grosirHeaders := []string{"No Produk", "Nama Paket", "Satuan", "Konversi", "Ref Harga Beli", "Harga Beli", "Ref Harga Jual", "Harga Jual"}
	grosirExample := []any{0, "1 Dus", "Dus", 12, "", 55000, "", 75000}

	refHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#808080"}, Pattern: 1},
	})
	refCellStyle, _ := f.NewStyle(&excelize.Style{
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Font:   &excelize.Font{Color: "#7F7F7F"},
		NumFmt: 3,
	})
	grosirCurrencyStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 3})

	for col, h := range grosirHeaders {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetGrosir, cell, h)
	}
	for col, val := range grosirExample {
		cell, _ := excelize.CoordinatesToCellName(col+1, 2)
		f.SetCellValue(sheetGrosir, cell, val)
	}
	f.SetCellStyle(sheetGrosir, "A1", "D1", headerStyle)
	f.SetCellStyle(sheetGrosir, "E1", "E1", refHeaderStyle)
	f.SetCellStyle(sheetGrosir, "F1", "F1", headerStyle)
	f.SetCellStyle(sheetGrosir, "G1", "G1", refHeaderStyle)
	f.SetCellStyle(sheetGrosir, "H1", "H1", headerStyle)
	f.SetCellStyle(sheetGrosir, "F2", "F500", grosirCurrencyStyle)
	f.SetCellStyle(sheetGrosir, "H2", "H500", grosirCurrencyStyle)
	f.SetCellStyle(sheetGrosir, "E2", "E500", refCellStyle)
	f.SetCellStyle(sheetGrosir, "G2", "G500", refCellStyle)
	for row := 2; row <= 500; row++ {
		cellA, _ := excelize.CoordinatesToCellName(1, row)
		cellD, _ := excelize.CoordinatesToCellName(4, row)
		cellE, _ := excelize.CoordinatesToCellName(5, row)
		cellG, _ := excelize.CoordinatesToCellName(7, row)
		f.SetCellFormula(sheetGrosir, cellE, fmt.Sprintf(`IFERROR(VLOOKUP(%s,Produk!$A:$J,5,0)*%s,"")`, cellA, cellD))
		f.SetCellFormula(sheetGrosir, cellG, fmt.Sprintf(`IFERROR(VLOOKUP(%s,Produk!$A:$J,6,0)*%s,"")`, cellA, cellD))
	}

	f.SetCellValue(sheetGrosir, "A3", "* No Produk: isi dengan nilai kolom 'No' dari sheet Produk. Satuan dipilih dari dropdown. Kolom Ref otomatis.")
	f.SetCellStyle(sheetGrosir, "A3", "H3", noteStyle)
	f.MergeCell(sheetGrosir, "A3", "H3")

	dvGNoProduk := excelize.NewDataValidation(true)
	dvGNoProduk.Sqref = "A4:A500"
	dvGNoProduk.SetRange(1, 999999, excelize.DataValidationTypeWhole, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvGNoProduk.SetError(excelize.DataValidationErrorStyleStop, "No produk tidak valid", "Isi dengan angka 'No' dari sheet Produk (bulat >= 1).")
	f.AddDataValidation(sheetGrosir, dvGNoProduk)

	dvGNamaPaket := excelize.NewDataValidation(true)
	dvGNamaPaket.Sqref = "B4:B500"
	dvGNamaPaket.SetRange(0, 100, excelize.DataValidationTypeTextLength, excelize.DataValidationOperatorBetween)
	dvGNamaPaket.SetError(excelize.DataValidationErrorStyleStop, "Nama paket tidak valid", "Nama paket maksimal 100 karakter.")
	f.AddDataValidation(sheetGrosir, dvGNamaPaket)

	if len(unitInfos) > 0 {
		lastUnitRow := len(unitInfos) + 1
		dvGSatuan := excelize.NewDataValidation(true)
		dvGSatuan.Sqref = "C4:C500"
		dvGSatuan.Formula1 = fmt.Sprintf("Satuan!$A$2:$A$%d", lastUnitRow)
		dvGSatuan.ShowDropDown = true
		dvGSatuan.SetError(excelize.DataValidationErrorStyleStop, "Satuan tidak valid", "Pilih satuan dari daftar yang tersedia di sheet Satuan.")
		f.AddDataValidation(sheetGrosir, dvGSatuan)
	}

	dvGKonversi := excelize.NewDataValidation(true)
	dvGKonversi.Sqref = "D4:D500"
	dvGKonversi.SetRange(0.001, 9999999999999.999, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvGKonversi.SetError(excelize.DataValidationErrorStyleStop, "Konversi tidak valid", "Konversi harus berupa angka > 0.")
	f.AddDataValidation(sheetGrosir, dvGKonversi)

	dvGHargaBeli := excelize.NewDataValidation(true)
	dvGHargaBeli.Sqref = "F4:F500"
	dvGHargaBeli.SetRange(0, 9999999999999.99, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvGHargaBeli.SetError(excelize.DataValidationErrorStyleStop, "Harga beli tidak valid", "Harga beli harus berupa angka >= 0.")
	f.AddDataValidation(sheetGrosir, dvGHargaBeli)

	dvGHargaJual := excelize.NewDataValidation(true)
	dvGHargaJual.Sqref = "H4:H500"
	dvGHargaJual.SetRange(0.01, 9999999999999.99, excelize.DataValidationTypeDecimal, excelize.DataValidationOperatorGreaterThanOrEqual)
	dvGHargaJual.SetError(excelize.DataValidationErrorStyleStop, "Harga jual tidak valid", "Harga jual harus berupa angka > 0.")
	f.AddDataValidation(sheetGrosir, dvGHargaJual)

	for _, col := range []string{"A", "B", "C", "D", "F", "H"} {
		f.SetColWidth(sheetGrosir, col, col, 16)
	}
	for _, col := range []string{"E", "G"} {
		f.SetColWidth(sheetGrosir, col, col, 18)
	}

	// ── Sheet 3: Kategori ────────────────────────────────────────────────────
	sheetKategori := "Kategori"
	f.NewSheet(sheetKategori)
	f.SetCellValue(sheetKategori, "A1", "nama_kategori")

	catHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E2EFDA"}, Pattern: 1},
	})
	f.SetCellStyle(sheetKategori, "A1", "A1", catHeaderStyle)
	for i, name := range categoryNames {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetCellValue(sheetKategori, cell, name)
	}
	f.SetColWidth(sheetKategori, "A", "A", 24)

	// ── Sheet 4: Satuan ──────────────────────────────────────────────────────
	sheetSatuan := "Satuan"
	f.NewSheet(sheetSatuan)
	f.SetCellValue(sheetSatuan, "A1", "nama_satuan")
	f.SetCellValue(sheetSatuan, "B1", "singkatan")
	f.SetCellStyle(sheetSatuan, "A1", "B1", catHeaderStyle)
	for i, u := range unitInfos {
		cellA, _ := excelize.CoordinatesToCellName(1, i+2)
		cellB, _ := excelize.CoordinatesToCellName(2, i+2)
		f.SetCellValue(sheetSatuan, cellA, u.Name)
		f.SetCellValue(sheetSatuan, cellB, u.Abbreviation)
	}
	f.SetColWidth(sheetSatuan, "A", "A", 20)
	f.SetColWidth(sheetSatuan, "B", "B", 14)

	if idx, err := f.GetSheetIndex(sheetProduk); err == nil {
		f.SetActiveSheet(idx)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=template_import_produk.xlsx")

	if err := f.Write(c.Writer); err != nil {
		c.Error(&errors.InternalServerError{Message: "Gagal generate template"})
	}
}
