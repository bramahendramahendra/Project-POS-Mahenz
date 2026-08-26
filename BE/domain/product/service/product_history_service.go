package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetPurchaseHistory(productID, page, limit int) (*dto.PurchaseHistoryResponse, error) {
	// Validate product exists
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if product == nil {
		return nil, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	data, err := s.repo.GetPurchaseHistory(productID, page, limit)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	return data, nil
}

func (s *productService) GetSaleHistory(productID, page, limit int, status string) (*dto.SaleHistoryResponse, error) {
	// Validate product exists
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if product == nil {
		return nil, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	// Validate status filter
	if status != "" && status != "completed" && status != "void" {
		return nil, &errors.BadRequestError{Message: "Status filter harus 'completed', 'void', atau kosong"}
	}

	data, err := s.repo.GetSaleHistory(productID, page, limit, status)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	return data, nil
}
