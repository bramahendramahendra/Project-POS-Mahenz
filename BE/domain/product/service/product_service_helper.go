package service

import (
	"fmt"

	dto "pos_api/domain/product/dto"
	dto_category "pos_api/domain/product_category/dto"
)

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

func (s *productService) GetUnitInfos() (data []*dto.UnitInfo, err error) {
	units, err := s.repoUnit.GetOptions()
	if err != nil {
		return
	}
	for _, u := range units {
		data = append(data, &dto.UnitInfo{Name: u.Name, Abbreviation: u.Abbreviation})
	}
	return
}

func (s *productService) createCategoryWithCode(name, description string) (int64, error) {
	base := ""
	for _, r := range name {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			if r >= 'a' {
				r -= 32
			}
			base += string(r)
		}
	}
	if len(base) > 3 {
		base = base[:3]
	}
	for len(base) < 3 {
		base += "X"
	}

	candidate := base
	for i := 2; i <= 99; i++ {
		exists, err := s.repoCategory.CheckCodeExists(candidate)
		if err != nil {
			return 0, err
		}
		if !exists {
			break
		}
		candidate = fmt.Sprintf("%s%d", base, i)
	}

	return s.repoCategory.Create(&dto_category.CreateCategoryRequest{
		Name:        name,
		Code:        candidate,
		Description: description,
	})
}

func (s *productService) GetLowStock() (data []*dto.LowStockProduct, err error) {
	dataDB, err := s.repo.GetLowStock()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.LowStockProduct{
			ID:       v.ID,
			Name:     v.Name,
			Stock:    v.Stock,
			MinStock: v.MinStock,
			UnitName: v.UnitName,
		})
	}

	return data, nil
}
