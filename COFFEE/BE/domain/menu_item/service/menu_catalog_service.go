package service

import (
	dto "coffee_api/domain/menu_item/dto"
	"coffee_api/errors"
)

// Catalog: data halaman kasir (tahap 4.7b). Menu tanpa kombinasi yang dijual tidak ditampilkan;
// kategori tanpa menu tidak ditampilkan. Menu habis tetap tampil (ditandai di kasir).
func (s *menuItemService) Catalog() (*dto.CatalogResponse, error) {
	cats, items, variants, groups, options, err := s.repo.CatalogData()
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	optionsByGroup := map[int][]dto.CatalogModifierOption{}
	for _, o := range options {
		optionsByGroup[o.GroupID] = append(optionsByGroup[o.GroupID], dto.CatalogModifierOption{ID: o.ID, Name: o.Name, Price: o.Price, IsDefault: o.IsDefault})
	}
	groupsByItem := map[int][]dto.CatalogModifierGroup{}
	for _, g := range groups {
		opts := optionsByGroup[g.ID]
		if len(opts) == 0 {
			continue
		}
		groupsByItem[g.MenuItemID] = append(groupsByItem[g.MenuItemID], dto.CatalogModifierGroup{
			ID: g.ID, Name: g.Name, SelectionType: g.SelectionType, IsRequired: g.IsRequired, Options: opts,
		})
	}
	variantsByItem := map[int][]dto.CatalogVariant{}
	for _, v := range variants {
		variantsByItem[v.MenuItemID] = append(variantsByItem[v.MenuItemID], dto.CatalogVariant{
			ID: v.ID, SizeID: v.SizeID, SizeName: v.SizeName, Temperature: v.Temperature, Price: v.Price,
		})
	}

	var productIDs []int
	for _, it := range items {
		if it.ProductID != nil {
			productIDs = append(productIDs, *it.ProductID)
		}
	}
	stock, err := s.repo.AvailableStock(productIDs)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	resp := &dto.CatalogResponse{Categories: []dto.CatalogCategory{}, Items: []dto.CatalogItem{}}
	usedCategory := map[int]bool{}
	for _, it := range items {
		vs := variantsByItem[it.ID]
		if len(vs) == 0 {
			continue
		}
		item := dto.CatalogItem{
			ID: it.ID, Name: it.Name, CategoryID: it.CategoryID, StockMode: it.StockMode,
			HasVariants: it.HasVariants, IsSoldOut: it.IsSoldOut, Variants: vs,
			ModifierGroups: groupsByItem[it.ID],
		}
		if item.ModifierGroups == nil {
			item.ModifierGroups = []dto.CatalogModifierGroup{}
		}
		if it.ProductID != nil {
			qty := stock[*it.ProductID]
			item.Stock = &qty
		}
		resp.Items = append(resp.Items, item)
		usedCategory[it.CategoryID] = true
	}
	for _, c := range cats {
		if usedCategory[c.ID] {
			resp.Categories = append(resp.Categories, dto.CatalogCategory{ID: c.ID, Name: c.Name})
		}
	}
	return resp, nil
}
