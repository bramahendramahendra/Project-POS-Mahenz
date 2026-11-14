package service

import (
	"fmt"
	"strings"

	dto "coffee_api/domain/menu_item/dto"
	model "coffee_api/domain/menu_item/model"
	repo "coffee_api/domain/menu_item/repo"
	"coffee_api/errors"

	"gorm.io/gorm"
)

func (s *menuItemService) GetAll(req *dto.GetAllRequest) ([]dto.MenuItemResponse, int64, error) {
	rows, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, err
	}
	data := make([]dto.MenuItemResponse, 0, len(rows))
	for _, v := range rows {
		data = append(data, toItemResponse(v))
	}
	return data, total, nil
}

func (s *menuItemService) GetByID(id int) (*dto.MenuItemDetailResponse, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}
	variants, err := s.repo.GetVariants(id)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(variants))
	for _, v := range variants {
		ids = append(ids, v.ID)
	}
	lines, err := s.repo.GetRecipeLines(ids)
	if err != nil {
		return nil, err
	}
	byVariant := map[int][]dto.RecipeLineResponse{}
	for _, l := range lines {
		byVariant[l.VariantID] = append(byVariant[l.VariantID], dto.RecipeLineResponse{
			IngredientID: l.IngredientID, IngredientName: l.IngredientName, Amount: l.Amount,
		})
	}
	groups, err := s.repo.GetModifierGroups(id)
	if err != nil {
		return nil, err
	}

	res := &dto.MenuItemDetailResponse{
		MenuItemResponse: toItemResponse(item),
		Variants:         make([]dto.VariantResponse, 0, len(variants)),
		ModifierGroups:   make([]dto.ModifierGroupRef, 0, len(groups)),
	}
	for _, v := range variants {
		recipe := byVariant[v.ID]
		if recipe == nil {
			recipe = []dto.RecipeLineResponse{}
		}
		res.Variants = append(res.Variants, dto.VariantResponse{
			ID: v.ID, SizeID: v.SizeID, SizeName: v.SizeName, Temperature: v.Temperature,
			Price: v.Price, IsActive: v.IsActive, Recipe: recipe,
		})
	}
	for _, g := range groups {
		res.ModifierGroups = append(res.ModifierGroups, dto.ModifierGroupRef{ID: g.GroupID, Name: g.GroupName})
	}
	return res, nil
}

func (s *menuItemService) Create(req *dto.SaveRequest, userID int) (*dto.MenuItemDetailResponse, error) {
	variants, err := s.validate(req)
	if err != nil {
		return nil, err
	}
	var newID int64
	err = s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		id, err := r.Create(req, userID)
		if err != nil {
			return err
		}
		newID = id
		return saveChildren(r, int(id), variants, nil, req.ModifierGroupIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(int(newID))
}

// Update: kombinasi yang sama (ukuran + suhu) mempertahankan id-nya (pesanan kelak merujuk id
// varian), kombinasi baru ditambah, kombinasi yang dihapus dari form dibuang beserta resepnya.
func (s *menuItemService) Update(req *dto.SaveRequest) error {
	current, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if current == nil {
		return &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}
	variants, err := s.validate(req)
	if err != nil {
		return err
	}
	existing, err := s.repo.GetVariants(req.ID)
	if err != nil {
		return err
	}
	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		if err := r.Update(req); err != nil {
			return err
		}
		return saveChildren(r, req.ID, variants, existing, req.ModifierGroupIDs)
	})
}

func (s *menuItemService) ToggleStatus(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	return s.repo.ToggleStatus(id)
}

// SetSoldOut: tandai habis / tersedia manual (keputusan 88: kasir boleh).
func (s *menuItemService) SetSoldOut(req *dto.SoldOutRequest) error {
	if _, err := s.mustGet(req.ID); err != nil {
		return err
	}
	return s.repo.SetSoldOut(req.ID, *req.IsSoldOut)
}

// Delete: varian, resep, dan kelompok tambahan yang terpasang ikut terhapus. (Mulai tahap 4.7,
// menu yang sudah pernah dipesan sebaiknya dinonaktifkan, bukan dihapus.)
func (s *menuItemService) Delete(id int) error {
	if _, err := s.mustGet(id); err != nil {
		return err
	}
	// Sejak 4.7a: menu yang sudah pernah dipesan tidak bisa dihapus (riwayat penjualan), cukup dinonaktifkan.
	sold, err := s.repo.CountOrderItems(id)
	if err != nil {
		return err
	}
	if sold > 0 {
		return &errors.BadRequestError{Message: "Menu sudah pernah dipesan dan tidak bisa dihapus. Nonaktifkan saja"}
	}
	return s.repo.Delete(id)
}

func (s *menuItemService) mustGet(id int) (*model.MenuItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, &errors.NotFoundError{Message: "Menu tidak ditemukan"}
	}
	return item, nil
}

// validate memeriksa aturan menu dan mengembalikan daftar varian yang akan disimpan
// (menu tanpa varian = satu varian tanpa ukuran/suhu).
func (s *menuItemService) validate(req *dto.SaveRequest) ([]dto.VariantInput, error) {
	req.Name = strings.TrimSpace(req.Name)
	bad := func(msg string) ([]dto.VariantInput, error) { return nil, &errors.BadRequestError{Message: msg} }

	exists, err := s.repo.NameExists(req.Name, req.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return bad("Nama menu sudah digunakan")
	}
	active, err := s.repo.CategoryIsActive(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if !active {
		return bad("Kategori menu tidak ditemukan atau tidak aktif")
	}

	// Barang jadi: satu menu = satu stok, tanpa varian, tanpa resep (keputusan 43, 89).
	if req.StockMode == dto.StockModeDirect {
		if req.HasVariants {
			return bad("Menu barang jadi tidak bisa memakai varian. Buat menu terpisah per rasa")
		}
		if len(req.Recipe) > 0 {
			return bad("Menu barang jadi tidak memakai resep")
		}
		// Tepat satu produk, dan produk itu belum dipakai menu lain (keputusan 91).
		if req.ProductID == nil {
			return bad("Menu barang jadi wajib memilih produk barang jadi")
		}
		active, err := s.repo.ProductIsActive(*req.ProductID)
		if err != nil {
			return nil, err
		}
		if !active {
			return bad("Produk barang jadi tidak ditemukan atau tidak aktif")
		}
		linked, err := s.repo.ProductLinked(*req.ProductID, req.ID)
		if err != nil {
			return nil, err
		}
		if linked {
			return bad("Produk ini sudah dipakai menu lain")
		}
	} else {
		// Menu racikan tidak memakai produk.
		req.ProductID = nil
	}

	var variants []dto.VariantInput
	if !req.HasVariants {
		if req.Price == nil {
			return bad("Harga wajib diisi")
		}
		if len(req.Variants) > 0 {
			return bad("Menu tanpa varian tidak boleh berisi daftar varian")
		}
		variants = []dto.VariantInput{{Price: *req.Price, IsActive: true, Recipe: req.Recipe}}
	} else {
		if len(req.Recipe) > 0 || req.Price != nil {
			return bad("Menu bervarian: harga dan resep diisi per kombinasi")
		}
		if len(req.Variants) == 0 {
			return bad("Minimal satu kombinasi varian")
		}
		seen := map[string]bool{}
		sizeIDs := map[int]bool{}
		activeCount := 0
		for _, v := range req.Variants {
			if v.SizeID == nil && v.Temperature == nil {
				return bad("Setiap kombinasi varian harus punya ukuran atau suhu")
			}
			key := variantKey(v.SizeID, v.Temperature)
			if seen[key] {
				return bad("Kombinasi varian ganda")
			}
			seen[key] = true
			if v.SizeID != nil {
				sizeIDs[*v.SizeID] = true
			}
			if v.IsActive {
				activeCount++
			}
		}
		if activeCount == 0 {
			return bad("Minimal satu kombinasi varian yang dijual (aktif)")
		}
		if err := s.ensureCount(keys(sizeIDs), s.repo.CountActiveSizes, "Ukuran tidak ditemukan atau tidak aktif"); err != nil {
			return nil, err
		}
		variants = req.Variants
	}

	// Resep: bahan harus ada dan tidak ganda dalam satu kombinasi.
	ingredientIDs := map[int]bool{}
	for i := range variants {
		perVariant := map[int]bool{}
		for j := range variants[i].Recipe {
			line := &variants[i].Recipe[j]
			line.Amount = strings.TrimSpace(line.Amount)
			if perVariant[line.IngredientID] {
				return bad("Bahan yang sama muncul dua kali dalam satu resep")
			}
			perVariant[line.IngredientID] = true
			ingredientIDs[line.IngredientID] = true
		}
	}
	if err := s.ensureCount(keys(ingredientIDs), s.repo.CountIngredients, "Bahan resep tidak ditemukan"); err != nil {
		return nil, err
	}

	groupIDs := map[int]bool{}
	for _, gid := range req.ModifierGroupIDs {
		if groupIDs[gid] {
			return bad("Kelompok tambahan dipilih dua kali")
		}
		groupIDs[gid] = true
	}
	if err := s.ensureCount(keys(groupIDs), s.repo.CountActiveModifierGroups, "Kelompok tambahan tidak ditemukan atau tidak aktif"); err != nil {
		return nil, err
	}
	return variants, nil
}

func (s *menuItemService) ensureCount(ids []int, counter func([]int) (int64, error), msg string) error {
	if len(ids) == 0 {
		return nil
	}
	n, err := counter(ids)
	if err != nil {
		return err
	}
	if int(n) != len(ids) {
		return &errors.BadRequestError{Message: msg}
	}
	return nil
}

// saveChildren menyimpan varian (mempertahankan id kombinasi yang sama), resep per varian,
// dan kelompok tambahan yang dipasang.
func saveChildren(r repo.MenuItemRepoInterface, itemID int, variants []dto.VariantInput, existing []*model.Variant, groupIDs []int) error {
	existingByKey := map[string]int{}
	for _, v := range existing {
		existingByKey[variantKey(v.SizeID, v.Temperature)] = v.ID
	}

	keep := []int{}
	variantIDs := make([]int, len(variants))
	for i := range variants {
		if id, found := existingByKey[variantKey(variants[i].SizeID, variants[i].Temperature)]; found {
			keep = append(keep, id)
			variantIDs[i] = id
		}
	}
	if err := r.DeleteVariantsNotIn(itemID, keep); err != nil {
		return err
	}
	for i := range variants {
		v := &variants[i]
		if variantIDs[i] > 0 {
			if err := r.UpdateVariant(variantIDs[i], v.Price, v.IsActive); err != nil {
				return err
			}
		} else {
			id, err := r.InsertVariant(itemID, v)
			if err != nil {
				return err
			}
			variantIDs[i] = int(id)
		}
		if err := r.DeleteRecipeLines(variantIDs[i]); err != nil {
			return err
		}
		for j := range v.Recipe {
			if err := r.InsertRecipeLine(variantIDs[i], &v.Recipe[j], j); err != nil {
				return err
			}
		}
	}
	if groupIDs == nil {
		groupIDs = []int{}
	}
	return r.ReplaceModifierGroups(itemID, groupIDs)
}

func variantKey(sizeID *int, temperature *string) string {
	size, temp := 0, ""
	if sizeID != nil {
		size = *sizeID
	}
	if temperature != nil {
		temp = *temperature
	}
	return fmt.Sprintf("%d|%s", size, temp)
}

func keys(m map[int]bool) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func toItemResponse(v *model.MenuItem) dto.MenuItemResponse {
	return dto.MenuItemResponse{
		ID: v.ID, CategoryID: v.CategoryID, CategoryName: v.CategoryName, Name: v.Name,
		StockMode: v.StockMode, ProductID: v.ProductID, ProductName: v.ProductName, HasVariants: v.HasVariants, IsActive: v.IsActive, IsSoldOut: v.IsSoldOut,
		PriceMin: v.PriceMin, PriceMax: v.PriceMax, VariantCount: v.VariantCount,
	}
}
