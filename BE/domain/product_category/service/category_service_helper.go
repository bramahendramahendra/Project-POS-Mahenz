package service

import (
	"fmt"
	"strings"
	"unicode"
)

func buildCategoryCode(name string) string {
	letters := strings.Map(func(ru rune) rune {
		if unicode.IsLetter(ru) {
			return unicode.ToUpper(ru)
		}
		return -1
	}, name)

	base := letters
	if len(base) > 3 {
		base = base[:3]
	}
	for len(base) < 3 {
		base += "X"
	}
	return base
}

func (s *categoryService) generateUniqueCode(name string) (string, error) {
	base := buildCategoryCode(name)
	candidate := base

	for i := 2; i <= 99; i++ {
		exists, err := s.repo.CheckCodeExists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s%d", base, i)
	}
	return "", fmt.Errorf("tidak bisa generate kode kategori yang unik")
}
