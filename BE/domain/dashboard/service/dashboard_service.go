package service

import (
	dto "pos_api/domain/dashboard/dto"
	"time"
)

// periodToDateRange menerjemahkan periode dashboard ("today"/"week"/"month") menjadi
// rentang tanggal (format YYYY-MM-DD, tanpa komponen jam — dibandingkan via DATE(...) di repo).
func periodToDateRange(period string) (startDate, endDate string) {
	now := time.Now()
	end := now.Format("2006-01-02")
	switch period {
	case "week":
		start := now.AddDate(0, 0, -6).Format("2006-01-02")
		return start, end
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		return start, end
	default: // "today"
		return end, end
	}
}

func (s *dashboardService) GetStats(period string) (*dto.StatsResponse, error) {
	startDate, endDate := periodToDateRange(period)

	todayStats, err := s.repo.GetStatsByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	todayExpenses, err := s.repo.GetExpensesByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	todayStats.TotalExpenses = todayExpenses
	todayStats.GrossProfit = todayStats.TotalSales - todayExpenses

	monthStats, err := s.repo.GetMonthStats()
	if err != nil {
		return nil, err
	}
	monthExpenses, err := s.repo.GetMonthExpenses()
	if err != nil {
		return nil, err
	}
	monthStats.TotalExpenses = monthExpenses
	monthStats.GrossProfit = monthStats.TotalSales - monthExpenses

	lowStock, err := s.repo.GetLowStockCount()
	if err != nil {
		return nil, err
	}
	openReceivables, err := s.repo.GetOpenReceivablesCount()
	if err != nil {
		return nil, err
	}

	return &dto.StatsResponse{
		Today:           *todayStats,
		ThisMonth:       *monthStats,
		LowStockCount:   lowStock,
		OpenReceivables: openReceivables,
	}, nil
}

func (s *dashboardService) GetSalesTrend(period string) ([]dto.SalesTrendItem, error) {
	days := 7
	switch period {
	case "30days":
		days = 30
	case "12months":
		days = 365
	}
	items, err := s.repo.GetSalesTrend(days)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SalesTrendItem{}
	}
	return items, nil
}

func (s *dashboardService) GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error) {
	items, err := s.repo.GetTopProducts(filter)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.TopProductItem{}
	}
	return items, nil
}

func (s *dashboardService) GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error) {
	items, err := s.repo.GetTopCategories(filter)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.TopCategoryItem{}
	}

	// hitung percentage
	var grandTotal float64
	for _, item := range items {
		grandTotal += item.TotalSales
	}
	if grandTotal > 0 {
		for i := range items {
			items[i].Percentage = (items[i].TotalSales / grandTotal) * 100
		}
	}
	return items, nil
}

func (s *dashboardService) GetSummaryExtra(period string) (*dto.SummaryExtraResponse, error) {
	now := time.Now()
	var startDate, endDate string

	switch period {
	case "today":
		today := now.Format("2006-01-02")
		startDate = today + " 00:00:00"
		endDate = today + " 23:59:59"
	case "7days":
		startDate = now.AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
		endDate = now.Format("2006-01-02") + " 23:59:59"
	case "30days":
		startDate = now.AddDate(0, 0, -29).Format("2006-01-02") + " 00:00:00"
		endDate = now.Format("2006-01-02") + " 23:59:59"
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02") + " 00:00:00"
		endDate = now.Format("2006-01-02") + " 23:59:59"
	default:
		today := now.Format("2006-01-02")
		startDate = today + " 00:00:00"
		endDate = today + " 23:59:59"
	}

	filter := dto.DateRangeFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	highest, err := s.repo.GetHighestTransaction(filter)
	if err != nil {
		return nil, err
	}

	peakHour, err := s.repo.GetPeakHour(filter)
	if err != nil {
		return nil, err
	}

	avg, err := s.repo.GetAvgTransaction(filter)
	if err != nil {
		return nil, err
	}

	return &dto.SummaryExtraResponse{
		Highest:  highest,
		PeakHour: peakHour,
		Avg:      avg,
	}, nil
}

func (s *dashboardService) GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error) {
	items, err := s.repo.GetPaymentMethods(filter)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.PaymentMethodItem{}
	}

	var grandTotal float64
	for _, item := range items {
		grandTotal += item.Total
	}
	if grandTotal > 0 {
		for i := range items {
			items[i].Percentage = (items[i].Total / grandTotal) * 100
		}
	}
	return items, nil
}
