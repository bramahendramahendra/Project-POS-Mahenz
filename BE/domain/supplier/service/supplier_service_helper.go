package service

import "fmt"

func (s *supplierService) generateUniqueCode() (string, error) {
	count, err := s.repo.GetCount()
	if err != nil {
		return "", err
	}
	for i := count + 1; ; i++ {
		code := fmt.Sprintf("SUP-%03d", i)
		exists, err := s.repo.CheckCodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
}
