package repo

import (
	"pos_api/domain/business_summary/dto"
	"time"

	"gorm.io/gorm"
)

type (
	BusinessSummaryRepoInterface interface {
		GetTodayStats(date string) (*dto.TodayStats, error)
		GetTodayExpenses(date string) (float64, error)
		GetStatsByRange(startDate, endDate string) (*dto.TodayStats, error)
		GetExpensesByRange(startDate, endDate string) (float64, error)
		GetMonthStats(month int, year int) (*dto.MonthStats, error)
		GetMonthExpenses(month int, year int) (float64, error)
		GetLowStockCount() (int64, error)
		GetOpenReceivablesCount() (int64, error)
		GetSalesTrend(days int, now time.Time) ([]dto.SalesTrendItem, error)
		GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error)
		GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error)
		GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error)
		GetHighestTransaction(filter dto.DateRangeFilter) (*dto.HighestTransactionItem, error)
		GetPeakHour(filter dto.DateRangeFilter) (*dto.PeakHourItem, error)
		GetAvgTransaction(filter dto.DateRangeFilter) (*dto.AvgTransactionItem, error)
	}

	businessSummaryRepo struct {
		db *gorm.DB
	}
)

func NewBusinessSummaryRepo(db *gorm.DB) *businessSummaryRepo {
	return &businessSummaryRepo{db: db}
}
