package service

import (
	codegen "pos_api/helper/codegen"
)

func (s *supplierService) generateUniqueCode() (string, error) {
	count, err := s.repo.GetCount()
	if err != nil {
		return "", err
	}
	return codegen.Sequential("SUP-", 3, count+1, s.repo.CheckCodeExists)
}
