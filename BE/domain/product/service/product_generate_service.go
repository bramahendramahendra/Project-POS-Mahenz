package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

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

	var sb strings.Builder
	for _, d := range digits {
		sb.WriteString(strconv.Itoa(d))
	}
	sb.WriteString(strconv.Itoa(checksum))
	barcode := sb.String()

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

