package service

import (
	"fmt"

	dto "pos_api/domain/product/dto"
	model "pos_api/domain/product/model"
	"pos_api/errors"
)

// validatePackageDrafts memastikan setiap draft paket di form Tambah Produk cuma
// merujuk (ref_temp_id) ke anchor (0) atau draft LAIN yang sudah muncul sebelumnya
// di array — supaya repo.Create bisa memprosesnya berurutan tanpa referensi ke depan
// atau melingkar. Nomor temp_id sendiri (bebas dari FE, biasanya counter) harus unik.
func validatePackageDrafts(drafts []dto.PackageDraftRequest) error {
	seenTempIDs := make(map[int]bool, len(drafts))
	seenTempIDs[0] = true // 0 selalu berarti anchor, dianggap "sudah ada" dari awal

	for i, p := range drafts {
		if seenTempIDs[p.TempID] {
			return &errors.BadRequestError{Message: fmt.Sprintf("temp_id %d dipakai lebih dari sekali", p.TempID)}
		}
		if !seenTempIDs[p.RefTempID] {
			return &errors.BadRequestError{Message: fmt.Sprintf("Paket ke-%d merujuk paket yang belum dibuat", i+1)}
		}
		if p.RefTempID == p.TempID {
			return &errors.BadRequestError{Message: fmt.Sprintf("Paket ke-%d tidak bisa merujuk dirinya sendiri", i+1)}
		}
		seenTempIDs[p.TempID] = true
	}
	return nil
}

func toPackageResponse(v *model.ProductPackage, factor float64) *dto.PackageResponse {
	return &dto.PackageResponse{
		ID:             v.ID,
		ProductID:      v.ProductID,
		UnitID:         v.UnitID,
		UnitName:       v.UnitName,
		Abbreviation:   v.Abbreviation,
		PackageName:    v.PackageName,
		RefPackageID:   v.RefPackageID,
		Qty:            v.Qty,
		RefQty:         v.RefQty,
		ResolvedFactor: factor,
		PurchasePrice:  v.PurchasePrice,
		SellingPrice:   v.SellingPrice,
		IsDefault:      v.IsDefault,
	}
}

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
		factor, factorErr := model.ResolvePackageFactor(dataDB, v.ID)
		if factorErr != nil {
			return data, &errors.InternalServerError{Message: factorErr.Error()}
		}
		data = append(data, toPackageResponse(v, factor))
	}

	return data, nil
}

func (s *productService) CreatePackage(productID int, req *dto.CreatePackageRequest) (data *dto.PackageResponse, err error) {
	exists, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, err
	}
	if exists == nil {
		return nil, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	existingPackages, err := s.repo.GetPackagesByProduct(productID)
	if err != nil {
		return nil, err
	}
	if _, factorErr := model.ResolvePackageFactor(existingPackages, req.RefPackageID); factorErr != nil {
		return nil, &errors.BadRequestError{Message: "Paket referensi tidak valid: " + factorErr.Error()}
	}

	newID, err := s.repo.CreatePackage(productID, req)
	if err != nil {
		return nil, err
	}

	allPackages, err := s.repo.GetPackagesByProduct(productID)
	if err != nil {
		return nil, err
	}
	created := findPackageByID(allPackages, int(newID))
	if created == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil paket baru"}
	}
	factor, factorErr := model.ResolvePackageFactor(allPackages, created.ID)
	if factorErr != nil {
		return nil, &errors.InternalServerError{Message: factorErr.Error()}
	}

	return toPackageResponse(created, factor), nil
}

func (s *productService) UpdatePackage(id, productID int, req *dto.UpdatePackageRequest) (data *dto.PackageResponse, err error) {
	exists, err := s.repo.GetByID(productID)
	if err != nil {
		return nil, err
	}
	if exists == nil {
		return nil, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	existingPackages, err := s.repo.GetPackagesByProduct(productID)
	if err != nil {
		return nil, err
	}
	target := findPackageByID(existingPackages, id)
	if target == nil {
		return nil, &errors.NotFoundError{Message: "Paket tidak ditemukan"}
	}
	if target.IsDefault {
		return nil, &errors.BadRequestError{Message: "Paket anchor (satuan pencatatan stok) tidak bisa diubah"}
	}
	if req.RefPackageID == id {
		return nil, &errors.BadRequestError{Message: "Paket tidak bisa mereferensikan dirinya sendiri"}
	}

	// Simulasikan perubahan sebelum disimpan supaya rantai referensi baru dicek dulu
	// tidak melingkar (mis. A jadi acuan B, lalu B diubah jadi acuan A).
	simulated := make([]*model.ProductPackage, len(existingPackages))
	for i, p := range existingPackages {
		copyP := *p
		if copyP.ID == id {
			refID := req.RefPackageID
			refQty := req.RefQty
			copyP.RefPackageID = &refID
			copyP.Qty = req.Qty
			copyP.RefQty = &refQty
		}
		simulated[i] = &copyP
	}
	if _, factorErr := model.ResolvePackageFactor(simulated, id); factorErr != nil {
		return nil, &errors.BadRequestError{Message: "Referensi tidak valid: " + factorErr.Error()}
	}

	if err = s.repo.UpdatePackage(id, productID, req); err != nil {
		return nil, err
	}

	allPackages, err := s.repo.GetPackagesByProduct(productID)
	if err != nil {
		return nil, err
	}
	updated := findPackageByID(allPackages, id)
	if updated == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil paket setelah diubah"}
	}
	factor, factorErr := model.ResolvePackageFactor(allPackages, updated.ID)
	if factorErr != nil {
		return nil, &errors.InternalServerError{Message: factorErr.Error()}
	}

	return toPackageResponse(updated, factor), nil
}

func (s *productService) DeletePackage(req *dto.DeletePackageRequest) (err error) {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	packages, err := s.repo.GetPackagesByProduct(req.ID)
	if err != nil {
		return err
	}
	target := findPackageByID(packages, req.PackageID)
	if target == nil {
		return &errors.NotFoundError{Message: "Paket tidak ditemukan"}
	}
	if target.IsDefault {
		return &errors.BadRequestError{Message: "Paket anchor (satuan pencatatan stok) tidak bisa dihapus"}
	}

	refCount, err := s.repo.CountPackagesReferencing(req.PackageID)
	if err != nil {
		return err
	}
	if refCount > 0 {
		return &errors.BadRequestError{Message: "Paket masih dipakai sebagai acuan paket lain, hapus/pindahkan itu dulu"}
	}

	return s.repo.DeletePackage(req.PackageID, req.ID)
}

func findPackageByID(packages []*model.ProductPackage, id int) *model.ProductPackage {
	for _, p := range packages {
		if p.ID == id {
			return p
		}
	}
	return nil
}
