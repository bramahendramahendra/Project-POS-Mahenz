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
		GetSalesItemsPaginated(req *dto.SalesListRequest) ([]dto.SalesItem, int64, error)
		GetSalesSummaryWithFilters(req *dto.SalesListRequest) (*dto.SalesSummary, error)
		GetProfitLossItems(params dto.FilterParams) ([]dto.ProfitLossItem, error)
		GetExpenseSummary(params dto.FilterParams) ([]dto.ExpenseSummaryItem, error)
		GetStockItems() ([]dto.StockItem, error)
		GetStockItemsPaginated(req *dto.StockListRequest) ([]dto.StockItem, int64, error)
		GetStockSummaryWithFilters(req *dto.StockSummaryRequest) (*dto.StockSummary, error)
		GetCashierItems(params dto.FilterParams) ([]dto.CashierItem, error)
		GetCashierItemsPaginated(req *dto.CashierReportRequest) ([]dto.CashierItem, int64, error)
	}

	reportRepo struct {
		db *gorm.DB
	}
)

func NewReportRepo(db *gorm.DB) *reportRepo {
	return &reportRepo{db: db}
}
