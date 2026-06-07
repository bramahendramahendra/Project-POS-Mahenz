package service

import (
	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"
	"pos_api/errors"
)

func (s *productService) GetAll(req *dto.ProductListRequest) (data []*dto.ProductResponse, total int64, err error) {
	items, t, err := s.repo.GetAll(req)
	if err != nil {
		return
	}
	total = int64(t)
	data = items
	return
}

func (s *productService) GetOptions() (data []*dto.ProductOption, err error) {
	return s.repo.GetOptions()
}

func (s *productService) GetByID(id int) (data dto.ProductResponse, err error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return
	}
	if p == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}
	return toProductResponse(p, ""), nil
}

func (s *productService) GetByBarcode(barcode string) (data dto.ProductResponse, err error) {
	p, err := s.repo.GetByBarcode(barcode)
	if err != nil {
		return
	}
	if p == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}
	return toProductResponse(p, ""), nil
}

func (s *productService) Search(keyword string, limit int) (data []*dto.ProductSearchResult, err error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	data, err = s.repo.Search(keyword, limit)
	return
}

func (s *productService) GetLowStock() (data []*dto.LowStockProduct, err error) {
	data, err = s.repo.GetLowStock()
	return
}

func (s *productService) Create(req *dto.ProductRequest) (data dto.ProductResponse, err error) {
	exists, err := s.repo.CheckBarcodeExists(req.Barcode, 0)
	if err != nil {
		return
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Barcode sudah digunakan"}
	}

	skuExists, err := s.repo.CheckSkuExists(req.SKU, 0)
	if err != nil {
		return
	}
	if skuExists {
		return data, &errors.BadRequestError{Message: "SKU sudah digunakan"}
	}

	newID, err := s.repo.Create(req)
	if err != nil {
		return
	}

	created, err := s.repo.GetByID(int(newID))
	if err != nil || created == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data produk baru"}
	}
	return toProductResponse(created, ""), nil
}

func (s *productService) Update(req *dto.UpdateProductRequest) (data dto.ProductResponse, err error) {
	p, err := s.repo.GetByID(req.ID)
	if err != nil {
		return
	}
	if p == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	exists, err := s.repo.CheckBarcodeExists(req.Barcode, req.ID)
	if err != nil {
		return
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Barcode sudah digunakan"}
	}

	skuExists, err := s.repo.CheckSkuExists(req.SKU, req.ID)
	if err != nil {
		return
	}
	if skuExists {
		return data, &errors.BadRequestError{Message: "SKU sudah digunakan"}
	}

	if err = s.repo.Update(req); err != nil {
		return
	}

	updated, err := s.repo.GetByID(req.ID)
	if err != nil || updated == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data produk"}
	}
	return toProductResponse(updated, ""), nil
}

func (s *productService) Delete(id int) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	count, err := s.repo.CountTransactionItems(id)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa penggunaan produk"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Produk tidak bisa dihapus karena sudah ada di transaksi"}
	}

	return s.repo.Delete(id)
}

func (s *productService) ToggleStatus(id int) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}
	return s.repo.ToggleStatus(id)
}

func toProductResponse(p *model.Product, categoryName string) dto.ProductResponse {
	catName := categoryName
	if catName == "" {
		catName = p.CategoryName
	}
	return dto.ProductResponse{
		ID:               p.ID,
		Barcode:          p.Barcode,
		SKU:              p.SKU,
		Name:             p.Name,
		CategoryID:       p.CategoryID,
		CategoryName:     catName,
		PurchasePrice:    p.PurchasePrice,
		SellingPrice:     p.SellingPrice,
		Stock:            p.Stock,
		MinStock:         p.MinStock,
		UnitID:           p.UnitID,
		UnitName:         p.UnitName,
		UnitAbbreviation: p.UnitAbbreviation,
		IsActive:         p.IsActive,
	}
}
