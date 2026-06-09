package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetAll(req *dto.GetAllRequest) (data []dto.ProductResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.ProductResponse{
			ID:               v.ID,
			Barcode:          v.Barcode,
			SKU:              v.SKU,
			Name:             v.Name,
			CategoryID:       v.CategoryID,
			CategoryName:     v.CategoryName,
			PurchasePrice:    v.PurchasePrice,
			SellingPrice:     v.SellingPrice,
			Stock:            v.Stock,
			MinStock:         v.MinStock,
			UnitID:           v.UnitID,
			UnitName:         v.UnitName,
			UnitAbbreviation: v.UnitAbbreviation,
			IsActive:         v.IsActive,
			ExtraPackages:    v.ExtraPackages,
			PriceTiersCount:  v.PriceTiersCount,
		})
	}

	return data, total, nil
}

func (s *productService) GetOptions() (data []*dto.GetOptionResponse, err error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetOptionResponse{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return data, nil
}

func (s *productService) Search(req *dto.SearchRequest) (data []*dto.SearchResponse, err error) {
	dataDB, err := s.repo.Search(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.SearchResponse{
			ID:           v.ID,
			Barcode:      v.Barcode,
			Name:         v.Name,
			SellingPrice: v.SellingPrice,
			Stock:        v.Stock,
			UnitID:       v.UnitID,
			UnitName:     v.UnitName,
		})
	}

	return data, nil
}

func (s *productService) GetByID(id int) (data dto.ProductResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	data = dto.ProductResponse{
		ID:               dataDB.ID,
		Barcode:          dataDB.Barcode,
		SKU:              dataDB.SKU,
		Name:             dataDB.Name,
		CategoryID:       dataDB.CategoryID,
		CategoryName:     dataDB.CategoryName,
		PurchasePrice:    dataDB.PurchasePrice,
		SellingPrice:     dataDB.SellingPrice,
		Stock:            dataDB.Stock,
		MinStock:         dataDB.MinStock,
		UnitID:           dataDB.UnitID,
		UnitName:         dataDB.UnitName,
		UnitAbbreviation: dataDB.UnitAbbreviation,
		IsActive:         dataDB.IsActive,
	}

	return data, nil
}

func (s *productService) GetByBarcode(barcode string) (data dto.ProductResponse, err error) {
	dataDB, err := s.repo.GetByBarcode(barcode)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	data = dto.ProductResponse{
		ID:               dataDB.ID,
		Barcode:          dataDB.Barcode,
		SKU:              dataDB.SKU,
		Name:             dataDB.Name,
		CategoryID:       dataDB.CategoryID,
		CategoryName:     dataDB.CategoryName,
		PurchasePrice:    dataDB.PurchasePrice,
		SellingPrice:     dataDB.SellingPrice,
		Stock:            dataDB.Stock,
		MinStock:         dataDB.MinStock,
		UnitID:           dataDB.UnitID,
		UnitName:         dataDB.UnitName,
		UnitAbbreviation: dataDB.UnitAbbreviation,
		IsActive:         dataDB.IsActive,
	}

	return data, nil
}

func (s *productService) Create(req *dto.CreateRequest) (data dto.ProductResponse, err error) {
	exists, err := s.repo.CheckBarcodeExists(req.Barcode, 0)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Barcode sudah digunakan"}
	}

	skuExists, err := s.repo.CheckSkuExists(req.SKU, 0)
	if err != nil {
		return data, err
	}
	if skuExists {
		return data, &errors.BadRequestError{Message: "SKU sudah digunakan"}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data produk baru"}
	}

	data = dto.ProductResponse{
		ID:               dataDB.ID,
		Barcode:          dataDB.Barcode,
		SKU:              dataDB.SKU,
		Name:             dataDB.Name,
		CategoryID:       dataDB.CategoryID,
		CategoryName:     dataDB.CategoryName,
		PurchasePrice:    dataDB.PurchasePrice,
		SellingPrice:     dataDB.SellingPrice,
		Stock:            dataDB.Stock,
		MinStock:         dataDB.MinStock,
		UnitID:           dataDB.UnitID,
		UnitName:         dataDB.UnitName,
		UnitAbbreviation: dataDB.UnitAbbreviation,
		IsActive:         dataDB.IsActive,
	}

	return data, nil
}

func (s *productService) Update(req *dto.UpdateRequest) (data dto.ProductResponse, err error) {
	existsUpdate, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if existsUpdate == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	exists, err := s.repo.CheckBarcodeExists(req.Barcode, req.ID)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Barcode sudah digunakan"}
	}

	skuExists, err := s.repo.CheckSkuExists(req.SKU, req.ID)
	if err != nil {
		return data, err
	}
	if skuExists {
		return data, &errors.BadRequestError{Message: "SKU sudah digunakan"}
	}

	if err = s.repo.Update(req); err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data produk"}
	}

	data = dto.ProductResponse{
		ID:               dataDB.ID,
		Barcode:          dataDB.Barcode,
		SKU:              dataDB.SKU,
		Name:             dataDB.Name,
		CategoryID:       dataDB.CategoryID,
		CategoryName:     dataDB.CategoryName,
		PurchasePrice:    dataDB.PurchasePrice,
		SellingPrice:     dataDB.SellingPrice,
		Stock:            dataDB.Stock,
		MinStock:         dataDB.MinStock,
		UnitID:           dataDB.UnitID,
		UnitName:         dataDB.UnitName,
		UnitAbbreviation: dataDB.UnitAbbreviation,
		IsActive:         dataDB.IsActive,
	}

	return data, nil
}

func (s *productService) Delete(req *dto.DeleteRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	count, err := s.repo.CountTransactionItems(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa penggunaan produk"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Produk tidak bisa dihapus karena sudah ada di transaksi"}
	}

	return s.repo.Delete(req)
}

func (s *productService) ToggleStatus(req *dto.ToggleStatusRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	return s.repo.ToggleStatus(req)
}
