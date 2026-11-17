package service

import (
	dto "pos_api/domain/business_summary/dto"
	time_helper "pos_api/helper/time"
	"time"
)

func periodToDateRange(period string) (startDate, endDate string) {
	now := time_helper.GetTimeNow()
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

// previousPeriodRange mengembalikan rentang periode SEBELUMNYA dengan panjang
// setara periode berjalan, untuk pembanding "naik/turun" di dashboard.
//   - today: kemarin (1 hari sebelum hari ini)
//   - week : 7 hari sebelum jendela minggu berjalan
//   - month: bulan kalender sebelumnya (tanggal 1 s/d akhir bulan lalu)
func previousPeriodRange(period string) (startDate, endDate string) {
	now := time_helper.GetTimeNow()
	switch period {
	case "week":
		// jendela berjalan: (now-6 .. now) -> sebelumnya: (now-13 .. now-7)
		start := now.AddDate(0, 0, -13).Format("2006-01-02")
		end := now.AddDate(0, 0, -7).Format("2006-01-02")
		return start, end
	case "month":
		firstThis := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		lastPrev := firstThis.AddDate(0, 0, -1)                                                  // hari terakhir bulan lalu
		firstPrev := time.Date(lastPrev.Year(), lastPrev.Month(), 1, 0, 0, 0, 0, now.Location()) // tanggal 1 bulan lalu
		return firstPrev.Format("2006-01-02"), lastPrev.Format("2006-01-02")
	default: // "today" -> kemarin
		y := now.AddDate(0, 0, -1).Format("2006-01-02")
		return y, y
	}
}

func (s *businessSummaryService) GetStats(period string) (*dto.StatsResponse, error) {
	startDate, endDate := periodToDateRange(period)

	todayStats, err := s.repo.GetStatsByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	todayExpenses, err := s.repo.GetExpensesByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	todayCOGS, err := s.repo.GetCOGSByRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	todayStats.TotalExpenses = todayExpenses
	todayStats.TotalCOGS = todayCOGS
	// Laba Kotor = Pendapatan - HPP, SELARAS dengan definisi "Laba Kotor" di Laporan Laba
	// Rugi (report_service.go: grossProfit = totalRevenue - totalCOGS). Pengeluaran
	// operasional (expense) SENGAJA tidak dikurangkan di sini -- itu baru mengurangi Laba
	// Bersih, bukan Laba Kotor. TotalExpenses tetap disimpan di response untuk transparansi,
	// tapi tidak ikut dalam GrossProfit.
	todayStats.GrossProfit = todayStats.TotalSales - todayCOGS

	now := time_helper.GetTimeNow()
	monthStats, err := s.repo.GetMonthStats(int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}
	monthExpenses, err := s.repo.GetMonthExpenses(int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}
	monthCOGS, err := s.repo.GetMonthCOGS(int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}
	monthStats.TotalExpenses = monthExpenses
	monthStats.TotalCOGS = monthCOGS
	monthStats.GrossProfit = monthStats.TotalSales - monthCOGS

	// Pembanding periode sebelumnya (panjang setara periode berjalan).
	prevStart, prevEnd := previousPeriodRange(period)
	prevStats, err := s.repo.GetStatsByRange(prevStart, prevEnd)
	if err != nil {
		return nil, err
	}
	prevCOGS, err := s.repo.GetCOGSByRange(prevStart, prevEnd)
	if err != nil {
		return nil, err
	}
	prev := dto.PrevStats{
		TotalTransactions: prevStats.TotalTransactions,
		TotalSales:        prevStats.TotalSales,
		GrossProfit:       prevStats.TotalSales - prevCOGS,
		// "Available" bila periode lalu memang ada transaksi -> pembanding bermakna.
		Available: prevStats.TotalTransactions > 0,
	}

	outStock, lowStock, err := s.repo.GetStockStatusCounts()
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
		Prev:            prev,
		LowStockCount:   outStock + lowStock, // habis + menipis (kompatibel lama)
		OutOfStockCount: outStock,
		LowOnlyCount:    lowStock,
		OpenReceivables: openReceivables,
	}, nil
}

func (s *businessSummaryService) GetSalesTrend(period string) ([]dto.SalesTrendItem, error) {
	days := 7
	switch period {
	case "30days":
		days = 30
	case "12months":
		days = 365
	}
	items, err := s.repo.GetSalesTrend(days, time_helper.GetTimeNow())
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SalesTrendItem{}
	}
	return items, nil
}

func (s *businessSummaryService) GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error) {
	items, err := s.repo.GetTopProducts(filter)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.TopProductItem{}
	}
	return items, nil
}
