package service

import (
	stderrors "errors"
	"fmt"
	"strings"

	dto "coffee_api/domain/supplier_purchase/dto"
	repo "coffee_api/domain/supplier_purchase/repo"
	"coffee_api/errors"

	"gorm.io/gorm"
)

// normalizeItems memeriksa dan melengkapi baris pembelian sesuai jenisnya (keputusan 92):
//   - product (bawaan): wajib product_id, tanpa ingredient_id, satuan wajib diisi.
//   - ingredient: wajib ingredient_id bahan aktif, tanpa product_id/package_id;
//     satuan diisi singkatan satuan stok bahan dan konversi 1.
//
// Baris yang sama (produk+kemasan atau bahan yang sama) tidak boleh dobel.
func (s *purchaseService) normalizeItems(items []dto.PurchaseRequest) error {
	type key struct {
		itemType string
		id       int
		pkgID    int
	}
	seen := make(map[key]bool, len(items))
	for i := range items {
		item := &items[i]
		if item.ItemType == "" {
			item.ItemType = repo.ItemTypeProduct
		}
		k := key{itemType: item.ItemType}
		switch item.ItemType {
		case repo.ItemTypeIngredient:
			if item.IngredientID <= 0 {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: bahan baku wajib dipilih", i+1)}
			}
			if item.ProductID != 0 || item.PackageID != nil {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: baris bahan baku tidak memakai produk atau kemasan", i+1)}
			}
			info, err := s.repo.GetIngredientInfo(item.IngredientID)
			if err != nil {
				return err
			}
			if info == nil {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: bahan baku tidak ditemukan", i+1)}
			}
			if !info.IsActive {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: bahan baku %s tidak aktif", i+1, info.Name)}
			}
			item.Unit = info.UnitAbbr
			item.ConversionQty = 1
			k.id = item.IngredientID
		default:
			if item.ProductID <= 0 {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: barang jadi wajib dipilih", i+1)}
			}
			if item.IngredientID != 0 {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: baris barang jadi tidak memakai bahan baku", i+1)}
			}
			if strings.TrimSpace(item.Unit) == "" {
				return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: satuan wajib diisi", i+1)}
			}
			k.id = item.ProductID
			if item.PackageID != nil {
				k.pkgID = *item.PackageID
			}
		}
		if seen[k] {
			msg := "produk dengan satuan ini sudah dipilih di baris lain"
			if item.ItemType == repo.ItemTypeIngredient {
				msg = "bahan baku ini sudah dipilih di baris lain"
			}
			return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: %s", i+1, msg)}
		}
		seen[k] = true
	}
	return nil
}

func validateDiscountAmount(items []dto.PurchaseRequest, discountAmount float64) error {
	var subtotal float64
	for _, item := range items {
		subtotal += item.PurchasePrice * item.Quantity
	}
	if discountAmount > subtotal {
		return &errors.BadRequestError{Message: fmt.Sprintf(
			"Diskon (%.2f) tidak boleh lebih besar dari subtotal (%.2f)", discountAmount, subtotal,
		)}
	}
	return nil
}

func (s *purchaseService) GetAll(req *dto.GetAllRequest) (data []dto.PurchaseResponse, total int64, err error) {
	data = []dto.PurchaseResponse{}
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.PurchaseResponse{
			ID:              v.ID,
			PurchaseCode:    v.PurchaseCode,
			InvoiceNumber:   v.InvoiceNumber,
			SupplierID:      v.SupplierID,
			SupplierName:    v.SupplierName,
			PurchaseDate:    v.PurchaseDate,
			DiscountAmount:  v.DiscountAmount,
			TotalAmount:     v.TotalAmount,
			PaidAmount:      v.PaidAmount,
			RemainingAmount: v.RemainingAmount,
			PaymentStatus:   v.PaymentStatus,
			Status:          v.Status,
			UserName:        v.UserName,
			Notes:           v.Notes,
		})
	}

	return data, total, nil
}

func (s *purchaseService) GetByID(id int) (data dto.PurchaseResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}

	items := make([]dto.PurchaseItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		items = append(items, dto.PurchaseItemResponse{
			ID:             v.ID,
			ItemType:       v.ItemType,
			ProductID:      v.ProductID,
			ProductName:    v.ProductName,
			IngredientID:   v.IngredientID,
			IngredientName: v.IngredientName,
			PackageID:      v.PackageID,
			Quantity:       v.Quantity,
			Unit:           v.Unit,
			ConversionQty:  v.ConversionQty,
			PurchasePrice:  v.PurchasePrice,
			Subtotal:       v.Subtotal,
		})
	}

	data = dto.PurchaseResponse{
		ID:              dataDB.ID,
		PurchaseCode:    dataDB.PurchaseCode,
		InvoiceNumber:   dataDB.InvoiceNumber,
		SupplierID:      dataDB.SupplierID,
		SupplierName:    dataDB.SupplierName,
		PurchaseDate:    dataDB.PurchaseDate,
		DiscountAmount:  dataDB.DiscountAmount,
		TotalAmount:     dataDB.TotalAmount,
		PaidAmount:      dataDB.PaidAmount,
		RemainingAmount: dataDB.RemainingAmount,
		PaymentStatus:   dataDB.PaymentStatus,
		Status:          dataDB.Status,
		UserName:        dataDB.UserName,
		Notes:           dataDB.Notes,
		Items:           items,
	}

	return data, nil
}

func (s *purchaseService) GenerateCode() (data dto.GenerateCodeResponse, err error) {
	dataDB, err := s.repo.GenerateCode()
	if err != nil {
		return data, err
	}

	data = dto.GenerateCodeResponse{PurchaseCode: dataDB}

	return data, nil
}

func (s *purchaseService) Create(req *dto.CreateRequest) (data dto.PurchaseResponse, err error) {
	if err := s.normalizeItems(req.Items); err != nil {
		return data, err
	}

	if err := validateDiscountAmount(req.Items, req.DiscountAmount); err != nil {
		return data, err
	}

	if req.PaymentMethod != "" {
		valid, err := s.repo.IsValidPaymentMethod(req.PaymentMethod)
		if err != nil {
			return data, err
		}
		if !valid {
			return data, &errors.BadRequestError{Message: "Metode pembayaran tidak valid"}
		}
	}

	dataDB, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	itemsDB := make([]dto.PurchaseItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		itemsDB = append(itemsDB, dto.PurchaseItemResponse{
			ID:             v.ID,
			ItemType:       v.ItemType,
			ProductID:      v.ProductID,
			ProductName:    v.ProductName,
			IngredientID:   v.IngredientID,
			IngredientName: v.IngredientName,
			Quantity:       v.Quantity,
			Unit:           v.Unit,
			ConversionQty:  v.ConversionQty,
			PurchasePrice:  v.PurchasePrice,
			Subtotal:       v.Subtotal,
		})
	}

	data = dto.PurchaseResponse{
		ID:              dataDB.ID,
		PurchaseCode:    dataDB.PurchaseCode,
		InvoiceNumber:   dataDB.InvoiceNumber,
		SupplierID:      dataDB.SupplierID,
		SupplierName:    dataDB.SupplierName,
		PurchaseDate:    dataDB.PurchaseDate,
		DiscountAmount:  dataDB.DiscountAmount,
		TotalAmount:     dataDB.TotalAmount,
		PaidAmount:      dataDB.PaidAmount,
		RemainingAmount: dataDB.RemainingAmount,
		PaymentStatus:   dataDB.PaymentStatus,
		Status:          dataDB.Status,
		UserName:        dataDB.UserName,
		Notes:           dataDB.Notes,
		Items:           itemsDB,
	}

	return data, nil
}

func (s *purchaseService) Update(req *dto.UpdateRequest) (data dto.PurchaseResponse, err error) {
	existing, err := s.repo.GetRawByID(req.ID)
	if err != nil {
		return data, err
	}
	if existing == nil {
		return data, &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}
	if existing.Status == "void" {
		return data, &errors.BadRequestError{Message: "PO sudah di-void, tidak bisa diedit"}
	}

	if err := s.normalizeItems(req.Items); err != nil {
		return data, err
	}

	if err := validateDiscountAmount(req.Items, req.DiscountAmount); err != nil {
		return data, err
	}

	if req.PaymentMethod != "" {
		valid, err := s.repo.IsValidPaymentMethod(req.PaymentMethod)
		if err != nil {
			return data, err
		}
		if !valid {
			return data, &errors.BadRequestError{Message: "Metode pembayaran tidak valid"}
		}
	}

	dataDB, err := s.repo.Update(req)
	if err != nil {
		return data, err
	}

	itemsDB := make([]dto.PurchaseItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		itemsDB = append(itemsDB, dto.PurchaseItemResponse{
			ID:             v.ID,
			ItemType:       v.ItemType,
			ProductID:      v.ProductID,
			ProductName:    v.ProductName,
			IngredientID:   v.IngredientID,
			IngredientName: v.IngredientName,
			Quantity:       v.Quantity,
			Unit:           v.Unit,
			ConversionQty:  v.ConversionQty,
			PurchasePrice:  v.PurchasePrice,
			Subtotal:       v.Subtotal,
		})
	}

	data = dto.PurchaseResponse{
		ID:              dataDB.ID,
		PurchaseCode:    dataDB.PurchaseCode,
		InvoiceNumber:   dataDB.InvoiceNumber,
		SupplierID:      dataDB.SupplierID,
		SupplierName:    dataDB.SupplierName,
		PurchaseDate:    dataDB.PurchaseDate,
		DiscountAmount:  dataDB.DiscountAmount,
		TotalAmount:     dataDB.TotalAmount,
		PaidAmount:      dataDB.PaidAmount,
		RemainingAmount: dataDB.RemainingAmount,
		PaymentStatus:   dataDB.PaymentStatus,
		Status:          dataDB.Status,
		UserName:        dataDB.UserName,
		Notes:           dataDB.Notes,
		Items:           itemsDB,
	}

	return data, nil
}

func (s *purchaseService) Delete(id int) error {
	exists, err := s.repo.GetRawByID(id)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}
	if exists.Status != "void" {
		return &errors.BadRequestError{Message: "PO harus di-void terlebih dahulu sebelum bisa dihapus"}
	}

	return s.repo.Delete(id)
}

func (s *purchaseService) Pay(req *dto.PayRequest) error {
	exists, err := s.repo.GetRawByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}
	if exists.Status == "void" {
		return &errors.BadRequestError{Message: "PO sudah di-void, tidak bisa dibayar"}
	}
	if exists.PaymentStatus == "paid" {
		return &errors.BadRequestError{Message: "PO sudah lunas"}
	}
	if req.Amount > exists.RemainingAmount {
		return &errors.BadRequestError{Message: "Jumlah pembayaran melebihi sisa tagihan"}
	}

	valid, err := s.repo.IsValidPaymentMethod(req.PaymentMethod)
	if err != nil {
		return err
	}
	if !valid {
		return &errors.BadRequestError{Message: "Metode pembayaran tidak valid"}
	}

	return s.repo.Pay(req)
}

func (s *purchaseService) Void(id int, userID int) error {
	exists, err := s.repo.GetRawByID(id)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}
	if exists.Status == "void" {
		return &errors.BadRequestError{Message: "PO sudah di-void"}
	}

	returnCount, err := s.repo.CountReturnsByPurchaseID(id)
	if err != nil {
		return err
	}
	if returnCount > 0 {
		return &errors.BadRequestError{Message: "PO tidak bisa di-void karena sudah ada retur supplier terkait"}
	}

	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		purchaseTx := s.repo.WithTx(tx)
		return purchaseTx.Void(id, userID)
	})
	if txErr != nil {
		var badReq *errors.BadRequestError
		if stderrors.As(txErr, &badReq) {
			return badReq
		}
		return &errors.InternalServerError{Message: txErr.Error()}
	}
	return nil
}

func (s *purchaseService) AddItems(req *dto.AddItemsRequest) (data dto.PurchaseResponse, err error) {
	existing, err := s.repo.GetRawByID(req.ID)
	if err != nil {
		return data, err
	}
	if existing == nil {
		return data, &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}
	if existing.Status == "void" {
		return data, &errors.BadRequestError{Message: "PO sudah di-void"}
	}
	if existing.PaymentStatus == "unpaid" {
		return data, &errors.BadRequestError{Message: "PO belum ada pembayaran, gunakan Edit untuk menambah item"}
	}

	if err := s.normalizeItems(req.Items); err != nil {
		return data, err
	}

	dataDB, err := s.repo.AddItems(req)
	if err != nil {
		return data, err
	}

	itemsDB := make([]dto.PurchaseItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		itemsDB = append(itemsDB, dto.PurchaseItemResponse{
			ID:             v.ID,
			ItemType:       v.ItemType,
			ProductID:      v.ProductID,
			ProductName:    v.ProductName,
			IngredientID:   v.IngredientID,
			IngredientName: v.IngredientName,
			Quantity:       v.Quantity,
			Unit:           v.Unit,
			ConversionQty:  v.ConversionQty,
			PurchasePrice:  v.PurchasePrice,
			Subtotal:       v.Subtotal,
		})
	}

	data = dto.PurchaseResponse{
		ID:              dataDB.ID,
		PurchaseCode:    dataDB.PurchaseCode,
		InvoiceNumber:   dataDB.InvoiceNumber,
		SupplierID:      dataDB.SupplierID,
		SupplierName:    dataDB.SupplierName,
		PurchaseDate:    dataDB.PurchaseDate,
		DiscountAmount:  dataDB.DiscountAmount,
		TotalAmount:     dataDB.TotalAmount,
		PaidAmount:      dataDB.PaidAmount,
		RemainingAmount: dataDB.RemainingAmount,
		PaymentStatus:   dataDB.PaymentStatus,
		Status:          dataDB.Status,
		UserName:        dataDB.UserName,
		Notes:           dataDB.Notes,
		Items:           itemsDB,
	}

	return data, nil
}

func (s *purchaseService) GetPayments(purchaseID int) (data []*dto.PaymentResponse, err error) {
	dataDB, err := s.repo.GetPayments(purchaseID)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.PaymentResponse{
			ID:            v.ID,
			PaymentDate:   v.PaymentDate,
			Amount:        v.Amount,
			PaymentMethod: v.PaymentMethod,
			Notes:         v.Notes,
			UserName:      v.UserName,
			CreatedAt:     v.CreatedAt,
		})
	}

	return data, nil
}
