package service

import (
	"strings"

	dto "pos_api/domain/product_category/dto"
	model "pos_api/domain/product_category/model"
	"pos_api/errors"
)

func (s *categoryService) GetAll() (data []dto.CategoryResponse, err error) {
	dataDB, err := s.repo.GetAll()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.CategoryResponse{
			ID:                 v.ID,
			Name:               v.Name,
			Code:               v.Code,
			Description:        v.Description,
			IsActive:           v.IsActive,
			ProductCount:       v.ProductCount,
			ActiveProductCount: v.ActiveProductCount,
			CreatedAt:          v.CreatedAt,
		})
	}

	return data, nil
}

func (s *categoryService) GetByID(id int) (data dto.CategoryResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}

	data = dto.CategoryResponse{
		ID:                 dataDB.ID,
		Name:               dataDB.Name,
		Code:               dataDB.Code,
		Description:        dataDB.Description,
		IsActive:           dataDB.IsActive,
		ProductCount:       dataDB.ProductCount,
		ActiveProductCount: dataDB.ActiveProductCount,
		CreatedAt:          dataDB.CreatedAt,
	}

	return data, nil
}

func (s *categoryService) Create(req *dto.CreateCategoryRequest) (data dto.CategoryResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	exists, err := s.repo.CheckNameExists(req.Name, 0)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama kategori sudah digunakan"}
	}

	code, err := s.repo.GenerateUniqueCode(req.Name)
	if err != nil {
		return data, err
	}
	req.Code = code

	newID, err := s.repo.Create(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}

	data = dto.CategoryResponse{
		ID:                 dataDB.ID,
		Name:               dataDB.Name,
		Code:               dataDB.Code,
		Description:        dataDB.Description,
		IsActive:           dataDB.IsActive,
		ProductCount:       dataDB.ProductCount,
		ActiveProductCount: dataDB.ActiveProductCount,
		CreatedAt:          dataDB.CreatedAt,
	}

	return data, nil
}

func (s *categoryService) Update(id int, req *dto.UpdateCategoryRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	c, err := s.repo.GetByID(id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if c == nil {
		return &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	exists, err := s.repo.CheckNameExists(req.Name, id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama kategori sudah digunakan"}
	}

	return s.repo.Update(id, req.Name, req.Description)
}

func (s *categoryService) Delete(id int) error {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if c == nil {
		return &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	count, err := s.repo.CountProductsByCategory(id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Kategori masih digunakan oleh produk"}
	}

	return s.repo.Delete(id)
}

func (s *categoryService) ToggleStatus(id int) error {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if c == nil {
		return &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	if c.IsActive {
		activeCount, err := s.repo.CountActiveProductsByCategory(id)
		if err != nil {
			return &errors.InternalServerError{Message: err.Error()}
		}
		if activeCount > 0 {
			return &errors.BadRequestError{Message: "Kategori tidak bisa dinonaktifkan karena masih memiliki produk aktif"}
		}
	}

	return s.repo.ToggleStatus(id)
}

func toCategoryResponse(c *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Code:               c.Code,
		Description:        c.Description,
		IsActive:           c.IsActive,
		ProductCount:       c.ProductCount,
		ActiveProductCount: c.ActiveProductCount,
		CreatedAt:          c.CreatedAt,
	}
}
