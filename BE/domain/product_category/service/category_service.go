package service

import (
	"strings"

	dto "pos_api/domain/product_category/dto"
	"pos_api/errors"
)

func (s *categoryService) GetAll(req *dto.CategoryListRequest) (data []dto.CategoryResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
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

	return data, total, nil
}

func (s *categoryService) GetOptions() (data []dto.CategoryOptionResponse, err error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.CategoryOptionResponse{
			ID:   v.ID,
			Name: v.Name,
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

	code, err := s.generateUniqueCode(req.Name)
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

func (s *categoryService) Update(req *dto.UpdateCategoryRequest) (data dto.CategoryResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	existsUpdate, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if existsUpdate == nil {
		return data, &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	exists, err := s.repo.CheckNameExists(req.Name, req.ID)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama kategori sudah digunakan"}
	}

	err = s.repo.Update(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
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

func (s *categoryService) Delete(req *dto.DeleteCategoryRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	count, err := s.repo.CountProductsByCategory(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Kategori masih digunakan oleh produk"}
	}

	return s.repo.Delete(req)
}

func (s *categoryService) ToggleStatus(req *dto.ToggleStatusCategoryRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Kategori tidak ditemukan"}
	}

	if exists.IsActive {
		activeCount, err := s.repo.CountActiveProductsByCategory(req.ID)
		if err != nil {
			return &errors.InternalServerError{Message: err.Error()}
		}
		if activeCount > 0 {
			return &errors.BadRequestError{Message: "Kategori tidak bisa dinonaktifkan karena masih memiliki produk aktif"}
		}
	}

	return s.repo.ToggleStatus(req)
}
