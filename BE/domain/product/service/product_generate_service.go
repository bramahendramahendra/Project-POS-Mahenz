package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GenerateBarcode() (data dto.GenerateBarcodeResponse, err error) {
	// EAN-13 dengan prefix 899 (Indonesia)
	digits := make([]int, 12)
	digits[0], digits[1], digits[2] = 8, 9, 9
	for i := 3; i < 12; i++ {
		n, randErr := rand.Int(rand.Reader, big.NewInt(10))
		if randErr != nil {
			return data, &errors.InternalServerError{Message: "Gagal generate barcode"}
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

	return dto.GenerateBarcodeResponse{Barcode: barcode}, nil
}

func (s *productService) GenerateSku(categoryID int) (data dto.GenerateSkuResponse, err error) {
	cat, err := s.catRepo.GetByID(categoryID)
	if err != nil {
		return
	}
	if cat == nil {
		return data, &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	count, err := s.repo.CountSkuByCategory(categoryID)
	if err != nil {
		return
	}

	return dto.GenerateSkuResponse{SKU: fmt.Sprintf("%s-%04d", cat.Code, count+1)}, nil
}

func (s *productService) GetCategoryNames() (data []string, err error) {
	cats, err := s.catRepo.GetOptions()
	if err != nil {
		return
	}
	for _, c := range cats {
		data = append(data, c.Name)
	}
	return
}

func (s *productService) GetUnitNames() (data []string, err error) {
	units, err := s.masterUnitRepo.GetOptions()
	if err != nil {
		return
	}
	for _, u := range units {
		data = append(data, u.Name)
	}
	return
}

func (s *productService) GetUnitInfos() (data []*dto.UnitInfo, err error) {
	units, err := s.masterUnitRepo.GetOptions()
	if err != nil {
		return
	}
	for _, u := range units {
		data = append(data, &dto.UnitInfo{Name: u.Name, Abbreviation: u.Abbreviation})
	}
	return
}
