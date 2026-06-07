package service_product

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	dto_product "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetCategoryNames() ([]string, error) {
	cats, err := s.catRepo.GetOptions()
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	names := make([]string, 0, len(cats))
	for _, c := range cats {
		names = append(names, c.Name)
	}
	return names, nil
}

func (s *productService) GetUnitNames() ([]string, error) {
	units, err := s.masterUnitRepo.GetOptions()
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	names := make([]string, 0, len(units))
	for _, u := range units {
		names = append(names, u.Name)
	}
	return names, nil
}

func (s *productService) GetUnitInfos() ([]*dto_product.UnitInfo, error) {
	units, err := s.masterUnitRepo.GetOptions()
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	infos := make([]*dto_product.UnitInfo, 0, len(units))
	for _, u := range units {
		infos = append(infos, &dto_product.UnitInfo{Name: u.Name, Abbreviation: u.Abbreviation})
	}
	return infos, nil
}

func (s *productService) GenerateBarcode() (*dto_product.GenerateBarcodeResponse, error) {
	// EAN-13 dengan prefix 899 (Indonesia)
	digits := make([]int, 12)
	digits[0], digits[1], digits[2] = 8, 9, 9
	for i := 3; i < 12; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, &errors.InternalServerError{Message: "Gagal generate barcode"}
		}
		digits[i] = int(n.Int64())
	}
	sum := 0
	for i, d := range digits {
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	checksum := (10 - (sum % 10)) % 10

	barcode := ""
	for _, d := range digits {
		barcode += strconv.Itoa(d)
	}
	barcode += strconv.Itoa(checksum)

	return &dto_product.GenerateBarcodeResponse{Barcode: barcode}, nil
}

func (s *productService) GenerateSku(categoryID int) (*dto_product.GenerateSkuResponse, error) {
	cat, err := s.catRepo.GetByID(categoryID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if cat == nil {
		return nil, &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	count, err := s.repo.CountSkuByCategory(categoryID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	sku := fmt.Sprintf("%s-%04d", cat.Code, count+1)
	return &dto_product.GenerateSkuResponse{SKU: sku}, nil
}
