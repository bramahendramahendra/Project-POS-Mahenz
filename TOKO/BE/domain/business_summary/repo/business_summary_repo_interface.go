package repo

import (
	"pos_api/domain/business_summary/dto"
	"time"

	"gorm.io/gorm"
)

type (
	BusinessSummaryRepoInterface interface {
		GetStatsByRange(startDate, endDate string) (*dto.TodayStats, error)
		GetExpensesByRange(startDate, endDate string) (float64, error)
		GetMonthStats(month int, year int) (*dto.MonthStats, error)
		GetMonthExpenses(month int, year int) (float64, error)
		GetCOGSByRange(startDate, endDate string) (float64, error)
		GetMonthCOGS(month int, year int) (float64, error)
		GetStockStatusCounts() (outCount int64, lowCount int64, err error)
		GetOpenReceivablesCount() (int64, error)
		GetSalesTrend(days int, now time.Time) ([]dto.SalesTrendItem, error)
		GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error)
	}

	businessSummaryRepo struct {
		db *gorm.DB
	}
)

func NewBusinessSummaryRepo(db *gorm.DB) *businessSummaryRepo {
	return &businessSummaryRepo{db: db}
}
