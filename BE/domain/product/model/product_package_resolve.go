package model

import "fmt"

// ResolvePackageFactor menghitung faktor konversi sebuah paket ke paket anchor produk
// (satuan pencatatan stok), dengan menelusuri rantai ref_package_id sampai ketemu anchor
// (ref_package_id NULL). Dipakai di mana pun stok perlu dikonversi dari satuan paket yang
// dipilih user (kasir, pembelian) ke satuan stok produk — lihat transaction_repo.go dan
// purchase_repo.go.
func ResolvePackageFactor(packages []*ProductPackage, packageID int) (float64, error) {
	byID := make(map[int]*ProductPackage, len(packages))
	for _, p := range packages {
		byID[p.ID] = p
	}

	current, ok := byID[packageID]
	if !ok {
		return 0, fmt.Errorf("paket %d tidak ditemukan", packageID)
	}

	factor := 1.0
	visited := make(map[int]bool, len(packages))
	for current.RefPackageID != nil {
		if visited[current.ID] {
			return 0, fmt.Errorf("konversi paket %d melingkar", packageID)
		}
		visited[current.ID] = true

		if current.RefQty == nil {
			return 0, fmt.Errorf("paket %d tidak punya rasio referensi", current.ID)
		}
		if current.Qty <= 0 {
			return 0, fmt.Errorf("paket %d punya qty tidak valid", current.ID)
		}
		factor *= *current.RefQty / current.Qty

		next, ok := byID[*current.RefPackageID]
		if !ok {
			return 0, fmt.Errorf("paket referensi %d tidak ditemukan", *current.RefPackageID)
		}
		current = next
	}

	return factor, nil
}
