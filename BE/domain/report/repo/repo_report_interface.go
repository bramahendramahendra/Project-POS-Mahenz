package repo

import (
	"pos_api/domain/report/dto"

	"gorm.io/gorm"
)

type (
	ReportRepoInterface interface {
		GetSalesItems(params dto.FilterParams) ([]dto.SalesItem, error)
		GetSalesSummary(params dto.FilterParams) (*dto.SalesSummary, error)
		GetSalesChart(params dto.FilterParams) ([]dto.SalesChartItem, error)
		GetProfitLossItems(params dto.FilterParams) ([]dto.ProfitLossItem, error)
		GetExpenseSummary(params dto.FilterParams) ([]dto.ExpenseSummaryItem, error)
		GetStockItems() ([]dto.StockItem, error)
		GetCashierItems(params dto.FilterParams) ([]dto.CashierItem, error)
	}

	reportRepo struct {
		db *gorm.DB
	}
)

func NewReportRepo(db *gorm.DB) *reportRepo {
	return &reportRepo{db: db}
}
