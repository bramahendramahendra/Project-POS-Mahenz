package service

import (
	stderrors "errors"
	"fmt"
	"math"

	dto "pos_api/domain/supplier_purchase/dto"
	model "pos_api/domain/supplier_purchase/model"
	"pos_api/errors"

	"gorm.io/gorm"
)

func mapExpiryBatches(batches []model.PurchaseItemExpiryBatch) []dto.PurchaseItemExpiryBatchResponse {
	if len(batches) == 0 {
		return nil
	}
	result := make([]dto.PurchaseItemExpiryBatchResponse, 0, len(batches))
	for _, b := range batches {
		result = append(result, dto.PurchaseItemExpiryBatchResponse{
			Qty:         b.Qty,
			ExpiredDate: b.ExpiredDate.Format("2006-01-02"),
		})
	}
	return result
}

func validateExpiryBatches(items []dto.PurchaseRequest) error {
	const epsilon = 0.0001
	for i, item := range items {
		if len(item.ExpiryBatches) == 0 {
			continue
		}
		var total float64
		for _, b := range item.ExpiryBatches {
			total += b.Qty
		}
		if math.Abs(total-item.Quantity) > epsilon {
			return &errors.BadRequestError{Message: fmt.Sprintf(
				"Item ke-%d: total qty batch expired (%.3f) harus sama dengan qty item (%.3f)",
				i+1, total, item.Quantity,
			)}
		}
	}
	return nil
}

func validateDuplicateProducts(items []dto.PurchaseRequest) error {
	type key struct {
		productID int
		packageID int
	}
	seen := make(map[key]bool, len(items))
	for i, item := range items {
		pkgID := 0
		if item.PackageID != nil {
			pkgID = *item.PackageID
		}
		k := key{productID: item.ProductID, packageID: pkgID}
		if seen[k] {
			return &errors.BadRequestError{Message: fmt.Sprintf(
				"Item ke-%d: produk dengan satuan ini sudah dipilih di baris lain", i+1,
			)}
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
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			ConversionQty: v.ConversionQty,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
			ExpiryBatches: mapExpiryBatches(v.ExpiryBatches),
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
	if err := validateDuplicateProducts(req.Items); err != nil {
		return data, err
	}

	if err := validateDiscountAmount(req.Items, req.DiscountAmount); err != nil {
		return data, err
	}

	if err := validateExpiryBatches(req.Items); err != nil {
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
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			ConversionQty: v.ConversionQty,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
			ExpiryBatches: mapExpiryBatches(v.ExpiryBatches),
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

	if err := validateDuplicateProducts(req.Items); err != nil {
		return data, err
	}

	if err := validateDiscountAmount(req.Items, req.DiscountAmount); err != nil {
		return data, err
	}

	if err := validateExpiryBatches(req.Items); err != nil {
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
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			ConversionQty: v.ConversionQty,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
			ExpiryBatches: mapExpiryBatches(v.ExpiryBatches),
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

	if err := validateDuplicateProducts(req.Items); err != nil {
		return data, err
	}

	if err := validateExpiryBatches(req.Items); err != nil {
		return data, err
	}

	dataDB, err := s.repo.AddItems(req)
	if err != nil {
		return data, err
	}

	itemsDB := make([]dto.PurchaseItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		itemsDB = append(itemsDB, dto.PurchaseItemResponse{
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			ConversionQty: v.ConversionQty,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
			ExpiryBatches: mapExpiryBatches(v.ExpiryBatches),
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
