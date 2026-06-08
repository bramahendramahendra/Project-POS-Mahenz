package service

import (
	dto "pos_api/domain/product/dto"
	"pos_api/errors"
)

func (s *productPackageService) GetByProduct(productID int) (data []*dto.ProductPackageResponse, err error) {
	exists, err := s.prodRepo.GetByID(productID)
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

	for _, v := range dataDB {
		data = append(data, &dto.ProductPackageResponse{
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

func (s *productPackageService) Save(req *dto.SaveProductPackagesRequest) (err error) {
	exists, err := s.prodRepo.GetByID(req.ProductID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	return s.repo.SavePackages(req.ProductID, req.Packages)
}

func (s *productPackageService) DeleteOne(req *dto.PackageIDUriRequest) (err error) {
	exists, err := s.prodRepo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	return s.repo.DeletePackage(req.PackageID, req.ID)
}
