package service

import (
	dto "pos_api/domain/supplier_purchase/dto"
	"pos_api/errors"
)

func (s *purchaseService) GetAll(req *dto.GetAllRequest) (data []dto.PurchaseResponse, total int64, err error) {
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
	if existing.PaidAmount > 0 {
		return data, &errors.BadRequestError{Message: "PO tidak bisa diedit karena sudah ada pembayaran"}
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
	if exists.PaidAmount > 0 {
		return &errors.BadRequestError{Message: "PO tidak bisa dihapus karena sudah ada pembayaran"}
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
