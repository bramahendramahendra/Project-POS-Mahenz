package service

import (
	"fmt"
	"strings"
	"unicode"

	"pos_api/errors"
)

const categoryCodeLength = 3

func buildCategoryCode(name string) string {
	letters := strings.Map(func(ru rune) rune {
		if unicode.IsLetter(ru) {
			return unicode.ToUpper(ru)
		}
		return -1
	}, name)

	base := letters
	if len(base) > categoryCodeLength {
		base = base[:categoryCodeLength]
	}
	for len(base) < categoryCodeLength {
		base += "X"
	}
	return base
}

func (s *categoryService) generateUniqueCode(name string) (string, error) {
	base := buildCategoryCode(name)

	exists, err := s.repo.CheckCodeExists(base)
	if err != nil {
		return "", err
	}
	if !exists {
		return base, nil
	}

	for i := 2; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%d", base, i)
		exists, err := s.repo.CheckCodeExists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", &errors.InternalServerError{Message: "tidak bisa generate kode kategori yang unik"}
}
