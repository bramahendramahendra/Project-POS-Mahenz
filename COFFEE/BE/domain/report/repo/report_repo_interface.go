package repo

import (
	"coffee_api/domain/report/dto"

	"gorm.io/gorm"
)

type (
	ReportRepoInterface interface {
		GetSalesItems(params dto.FilterParams) ([]dto.SalesItem, error)
		GetSalesSummary(params dto.FilterParams) (*dto.SalesSummary, error)
		GetSalesItemsPaginated(req *dto.SalesListRequest) ([]dto.SalesItem, int64, error)
		GetSalesSummaryWithFilters(req *dto.SalesListRequest) (*dto.SalesSummary, error)
		GetProfitLossItems(params dto.FilterParams) ([]dto.ProfitLossItem, error)
		GetExpenseSummary(params dto.FilterParams) ([]dto.ExpenseSummaryItem, error)
		GetStockItems() ([]dto.StockItem, error)
		GetStockItemsPaginated(req *dto.StockListRequest) ([]dto.StockItem, int64, error)
		GetStockSummaryWithFilters(req *dto.StockSummaryRequest) (*dto.StockSummary, error)
		GetCashierItems(params dto.FilterParams) ([]dto.CashierItem, error)
		// Laporan Coffee (tahap 4.10).
		GetSalesBreakdownRows(params dto.FilterParams) ([]SalesBreakdownRow, error)
		GetSalesByPayment(params dto.FilterParams) ([]dto.SalesGroupRow, error)
		GetSalesByOrderType(params dto.FilterParams) ([]dto.SalesGroupRow, error)
		GetIngredientReport(params dto.FilterParams) ([]dto.IngredientReportRow, error)
		GetWasteList(params dto.FilterParams) ([]dto.WasteRow, error)
		GetIngredientCost(params dto.FilterParams) (float64, error)
		GetCashierItemsPaginated(req *dto.CashierReportRequest) ([]dto.CashierItem, int64, error)
	}

	reportRepo struct {
		db *gorm.DB
	}
)

func NewReportRepo(db *gorm.DB) *reportRepo {
	return &reportRepo{db: db}
}
