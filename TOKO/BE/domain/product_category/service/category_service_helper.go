package service

import (
	codegen "pos_api/helper/codegen"
)

const categoryCodeLength = 3

func (s *categoryService) generateUniqueCode(name string) (string, error) {
	base := codegen.BuildLetterPrefix(name, categoryCodeLength)
	return codegen.UniqueByPrefix(base, s.repo.CheckCodeExists)
}
