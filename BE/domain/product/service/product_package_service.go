package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productService) GetPackagesByProduct(productID int) (data []*dto.PackageResponse, err error) {
	exists, err := s.repo.GetByID(productID)
	if err != nil {
		return data, err
	}
	if exists == nil {
		return data, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	dataDB, err := s.repo.GetPackagesByProduct(productID)
	if err != nil {
		return data, err
	}

	data = make([]*dto.PackageResponse, 0, len(dataDB))
	for _, v := range dataDB {
		data = append(data, &dto.PackageResponse{
			ID:            v.ID,
			ProductID:     v.ProductID,
			UnitID:        v.UnitID,
			UnitName:      v.UnitName,
			Abbreviation:  v.Abbreviation,
			PackageName:   v.PackageName,
			ConversionQty: v.ConversionQty,
			PurchasePrice: v.PurchasePrice,
			SellingPrice:  v.SellingPrice,
			IsDefault:     v.IsDefault,
		})
	}

	return data, nil
}

func (s *productService) SavePackages(req *dto.SavePackageRequest) (err error) {
	exists, err := s.repo.GetByID(req.ProductID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	err = s.repo.SavePackages(req.ProductID, req.Packages)
	return err
}

func (s *productService) DeletePackage(req *dto.DeletePackageRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	err = s.repo.DeletePackage(req.PackageID, req.ID)
	return err
}
