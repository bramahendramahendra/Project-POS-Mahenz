package service

import (
	"time"

	dto "pos_api/domain/supplier_return/dto"
	"pos_api/errors"
)

func (s *supplierReturnService) GetAll(req *dto.SupplierReturnListRequest) (data []dto.SupplierReturnResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.SupplierReturnResponse{
			ID:                v.ID,
			ReturnCode:        v.ReturnCode,
			PurchaseID:        v.PurchaseID,
			SupplierID:        v.SupplierID,
			SupplierName:      v.SupplierName,
			ReturnDate:        v.ReturnDate,
			TotalReturnAmount: v.TotalReturnAmount,
			Reason:            v.Reason,
			Status:            v.Status,
			UserName:          v.UserName,
			Notes:             v.Notes,
		})
	}

	return data, total, nil
}

func (s *supplierReturnService) GetByID(id int) (data dto.SupplierReturnResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Retur supplier tidak ditemukan"}
	}

	items := make([]dto.SupplierReturnItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		items = append(items, dto.SupplierReturnItemResponse{
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
		})
	}

	data = dto.SupplierReturnResponse{
		ID:                dataDB.ID,
		ReturnCode:        dataDB.ReturnCode,
		PurchaseID:        dataDB.PurchaseID,
		SupplierID:        dataDB.SupplierID,
		SupplierName:      dataDB.SupplierName,
		ReturnDate:        dataDB.ReturnDate,
		TotalReturnAmount: dataDB.TotalReturnAmount,
		Reason:            dataDB.Reason,
		Status:            dataDB.Status,
		UserName:          dataDB.UserName,
		Notes:             dataDB.Notes,
		Items:             items,
	}

	return data, nil
}

func (s *supplierReturnService) Create(req *dto.CreateSupplierReturnRequest) (data dto.SupplierReturnResponse, err error) {
	returnDate, err := time.Parse("2006-01-02", req.ReturnDate)
	if err != nil {
		return data, &errors.BadRequestError{Message: "Format tanggal retur tidak valid"}
	}
	if returnDate.After(time.Now().Truncate(24 * time.Hour)) {
		return data, &errors.BadRequestError{Message: "Tanggal retur tidak boleh lebih dari hari ini"}
	}

	purchaseDateStr, err := s.repo.GetPurchaseDate(req.PurchaseID)
	if err != nil {
		return data, err
	}
	if purchaseDateStr == "" {
		return data, &errors.NotFoundError{Message: "Purchase order tidak ditemukan"}
	}

	purchaseDate, err := time.Parse("2006-01-02", purchaseDateStr[:10])
	if err != nil {
		return data, err
	}
	if returnDate.Before(purchaseDate) {
		return data, &errors.BadRequestError{Message: "Tanggal retur tidak boleh lebih awal dari tanggal pembelian"}
	}

	dataDB, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	itemsDB := make([]dto.SupplierReturnItemResponse, 0, len(dataDB.Items))
	for _, v := range dataDB.Items {
		itemsDB = append(itemsDB, dto.SupplierReturnItemResponse{
			ID:            v.ID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			Quantity:      v.Quantity,
			Unit:          v.Unit,
			PurchasePrice: v.PurchasePrice,
			Subtotal:      v.Subtotal,
		})
	}

	data = dto.SupplierReturnResponse{
		ID:                dataDB.ID,
		ReturnCode:        dataDB.ReturnCode,
		PurchaseID:        dataDB.PurchaseID,
		SupplierID:        dataDB.SupplierID,
		SupplierName:      dataDB.SupplierName,
		ReturnDate:        dataDB.ReturnDate,
		TotalReturnAmount: dataDB.TotalReturnAmount,
		Reason:            dataDB.Reason,
		Status:            dataDB.Status,
		UserName:          dataDB.UserName,
		Notes:             dataDB.Notes,
		Items:             itemsDB,
	}

	return data, nil
}

func (s *supplierReturnService) UpdateStatus(req *dto.UpdateStatusRequest) error {
	current, err := s.repo.GetStatus(req.ID)
	if err != nil {
		return err
	}
	if current == "" {
		return &errors.NotFoundError{Message: "Retur supplier tidak ditemukan"}
	}
	if current == "approved" {
		return &errors.BadRequestError{Message: "Retur yang sudah approved tidak bisa diubah"}
	}

	if req.Status == "approved" {
		if err := s.repo.ApproveWithStockReduction(req.ID, req.UserID); err != nil {
			if badReq, ok := err.(*errors.BadRequestError); ok {
				return badReq
			}
			return err
		}
		return nil
	}

	if req.Status == "rejected" && req.Notes == "" {
		return &errors.BadRequestError{Message: "Catatan penolakan wajib diisi"}
	}

	return s.repo.UpdateStatus(req.ID, req.Status, req.Notes)
}

func (s *supplierReturnService) Delete(req *dto.GetSupplierReturnByIDRequest) error {
	exists, err := s.repo.GetStatus(req.ID)
	if err != nil {
		return err
	}
	if exists == "" {
		return &errors.NotFoundError{Message: "Retur supplier tidak ditemukan"}
	}
	if exists == "approved" {
		return &errors.BadRequestError{Message: "Retur yang sudah approved tidak bisa dihapus"}
	}

	return s.repo.Delete(req)
}
