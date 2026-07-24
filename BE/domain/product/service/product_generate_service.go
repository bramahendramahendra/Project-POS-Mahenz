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
	// EAN-13 dengan prefix 899 (Indonesia). Dicoba berulang sampai dapat yang belum
	// dipakai — sama seperti GenerateSku di bawah — karena 9 digit acak bisa saja
	// bentrok dengan barcode produk lain yang sudah ada (khususnya saat import massal
	// memanggil ini berkali-kali dalam satu batch).
	for {
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

		exists, checkErr := s.repo.CheckBarcodeExists(barcode, 0)
		if checkErr != nil {
			return data, checkErr
		}
		if !exists {
			return dto.GenerateBarcodeResponse{Barcode: barcode}, nil
		}
	}
}

func (s *productService) GenerateSku(categoryID int) (data dto.GenerateSkuResponse, err error) {
	dataDB, err := s.repoCategory.GetByID(categoryID)
	if err != nil {
		return
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	count, err := s.repo.CountSkuByCategory(categoryID)
	if err != nil {
		return
	}

	// Count baris tidak selalu = nomor urut SKU tertinggi (bisa ada SKU manual yang
	// melompat/tidak sekuensial), dan dua request generate yang hampir bersamaan bisa
	// membaca count yang sama sebelum salah satu produk tersimpan — jadi tetap perlu
	// dicek keunikannya, bukan cuma count+1 langsung dipakai.
	for next := count + 1; ; next++ {
		candidate := fmt.Sprintf("%s-%04d", dataDB.Code, next)
		exists, err := s.repo.CheckSkuExists(candidate, 0)
		if err != nil {
			return data, err
		}
		if !exists {
			return dto.GenerateSkuResponse{SKU: candidate}, nil
		}
	}
}
