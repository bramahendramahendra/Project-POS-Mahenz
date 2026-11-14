package service

import (
	"sort"

	dto "coffee_api/domain/report/dto"
	time_helper "coffee_api/helper/time"
)

// Laporan Coffee (tahap 4.10, keputusan 114-117).

func coffeeParams(req *dto.CoffeeReportRequest) dto.FilterParams {
	from, to := time_helper.NormalizeDateRange(req.DateFrom, req.DateTo)
	return dto.FilterParams{DateFrom: from, DateTo: to}
}

func share(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return part / total * 100
}

// GetSalesBreakdown: per menu (bisa dibuka per kombinasi), per kategori menu, per metode bayar, per jenis pesanan.
func (s *reportService) GetSalesBreakdown(req *dto.CoffeeReportRequest) (*dto.SalesBreakdownResponse, error) {
	params := coffeeParams(req)
	rows, err := s.repo.GetSalesBreakdownRows(params)
	if err != nil {
		return nil, err
	}
	resp := &dto.SalesBreakdownResponse{Menus: []dto.SalesMenuRow{}, Categories: []dto.SalesCategoryRow{}}

	menuIdx := map[int]int{}
	catIdx := map[string]int{}
	for _, r := range rows {
		resp.TotalRevenue += r.Revenue
		resp.TotalQty += r.Qty

		i, ok := menuIdx[r.MenuItemID]
		if !ok {
			resp.Menus = append(resp.Menus, dto.SalesMenuRow{MenuItemID: r.MenuItemID, MenuName: r.MenuName, CategoryName: r.CategoryName, Variants: []dto.SalesVariantRow{}})
			i = len(resp.Menus) - 1
			menuIdx[r.MenuItemID] = i
		}
		m := &resp.Menus[i]
		m.Qty += r.Qty
		m.Revenue += r.Revenue
		if r.VariantName != "" {
			m.Variants = append(m.Variants, dto.SalesVariantRow{VariantName: r.VariantName, Qty: r.Qty, Revenue: r.Revenue})
		}

		cat := r.CategoryName
		if cat == "" {
			cat = "Tanpa kategori"
		}
		j, ok := catIdx[cat]
		if !ok {
			resp.Categories = append(resp.Categories, dto.SalesCategoryRow{CategoryName: cat})
			j = len(resp.Categories) - 1
			catIdx[cat] = j
		}
		resp.Categories[j].Qty += r.Qty
		resp.Categories[j].Revenue += r.Revenue
	}
	for i := range resp.Menus {
		resp.Menus[i].Share = share(resp.Menus[i].Revenue, resp.TotalRevenue)
		sort.Slice(resp.Menus[i].Variants, func(a, b int) bool { return resp.Menus[i].Variants[a].Revenue > resp.Menus[i].Variants[b].Revenue })
	}
	for i := range resp.Categories {
		resp.Categories[i].Share = share(resp.Categories[i].Revenue, resp.TotalRevenue)
	}
	sort.SliceStable(resp.Menus, func(a, b int) bool { return resp.Menus[a].Revenue > resp.Menus[b].Revenue })
	sort.SliceStable(resp.Categories, func(a, b int) bool { return resp.Categories[a].Revenue > resp.Categories[b].Revenue })

	if resp.PaymentMethod, err = s.repo.GetSalesByPayment(params); err != nil {
		return nil, err
	}
	if resp.OrderTypes, err = s.repo.GetSalesByOrderType(params); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetIngredientReport: per bahan di periode (keputusan 116); nilai pakai harga beli rata-rata (keputusan 114).
func (s *reportService) GetIngredientReport(req *dto.CoffeeReportRequest) (*dto.IngredientReportResponse, error) {
	rows, err := s.repo.GetIngredientReport(coffeeParams(req))
	if err != nil {
		return nil, err
	}
	resp := &dto.IngredientReportResponse{Items: rows}
	for i := range resp.Items {
		it := &resp.Items[i]
		it.UsedValue = it.UsedQty * it.AvgPrice
		it.WastedValue = it.WastedQty * it.AvgPrice
		resp.TotalPurchaseValue += it.PurchaseValue
		resp.TotalUsedValue += it.UsedValue
		resp.TotalWastedValue += it.WastedValue
	}
	return resp, nil
}

// GetWasteReport: bahan dan barang jadi terbuang beserta alasan, pelaku, waktu, nilai (keputusan 116).
func (s *reportService) GetWasteReport(req *dto.CoffeeReportRequest) (*dto.WasteReportResponse, error) {
	rows, err := s.repo.GetWasteList(coffeeParams(req))
	if err != nil {
		return nil, err
	}
	resp := &dto.WasteReportResponse{Items: rows}
	for _, r := range rows {
		resp.TotalValue += r.Value
		if r.Kind == "ingredient" {
			resp.TotalIngredientValue += r.Value
		} else {
			resp.TotalProductValue += r.Value
		}
	}
	return resp, nil
}
