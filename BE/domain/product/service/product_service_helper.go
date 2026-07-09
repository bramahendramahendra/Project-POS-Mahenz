package service

import (
	dto "pos_api/domain/product/dto"
	dto_category "pos_api/domain/product_category/dto"
	codegen "pos_api/helper/codegen"
)

// categoryCodeLength harus sama dengan panjang kode yang dipakai domain product_category
// (lihat domain/product_category/service/category_service_helper.go) supaya kategori yang
// dibuat lewat import produk konsisten dengan yang dibuat lewat endpoint kategori langsung.
const categoryCodeLength = 3

func (s *productService) GetLowStock() (data []*dto.GetLowStockResponse, err error) {
	dataDB, err := s.repo.GetLowStock()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetLowStockResponse{
			ID:       v.ID,
			Name:     v.Name,
			Stock:    v.Stock,
			MinStock: v.MinStock,
			UnitName: v.UnitName,
		})
	}

	return data, nil
}

func (s *productService) GetCategoryNames() (data []string, err error) {
	cats, err := s.repoCategory.GetOptions()
	if err != nil {
		return
	}
	for _, c := range cats {
		data = append(data, c.Name)
	}
	return
}

func (s *productService) GetUnitInfos() (data []*dto.GetUnitInfoResponse, err error) {
	units, err := s.repoUnit.GetOptions()
	if err != nil {
		return
	}
	for _, u := range units {
		data = append(data, &dto.GetUnitInfoResponse{Name: u.Name, Abbreviation: u.Abbreviation})
	}
	return
}

func (s *productService) createCategoryWithCode(name, description string) (int64, error) {
	base := codegen.BuildLetterPrefix(name, categoryCodeLength)
	code, err := codegen.UniqueByPrefix(base, s.repoCategory.CheckCodeExists)
	if err != nil {
		return 0, err
	}

	return s.repoCategory.Create(&dto_category.CreateRequest{
		Name:        name,
		Code:        code,
		Description: description,
	})
}
