package service

import (
	"bytes"

	"coffee_api/domain/report/dto"
	repo "coffee_api/domain/report/repo"
)

type (
	ReportServiceInterface interface {
		GetSalesReport(params dto.FilterParams) (*dto.SalesReportResponse, error)
		ExportSalesReport(params dto.FilterParams) (*bytes.Buffer, error)
		GetSalesList(req *dto.SalesListRequest) ([]dto.SalesItem, int64, error)
		GetSalesSummaryData(req *dto.SalesListRequest) (*dto.SalesSummary, error)

		GetProfitLoss(params dto.FilterParams) (*dto.ProfitLossResponse, error)
		GetProfitLossData(req *dto.ProfitLossRequest) (*dto.ProfitLossResponse, error)
		ExportProfitLoss(params dto.FilterParams) (*bytes.Buffer, error)

		GetStockReport() (*dto.StockReportResponse, error)
		GetStockList(req *dto.StockListRequest) ([]dto.StockItem, int64, error)
		GetStockSummaryData(req *dto.StockSummaryRequest) (*dto.StockSummary, error)
		ExportStockReport() (*bytes.Buffer, error)

		GetCashierReport(params dto.FilterParams) ([]dto.CashierItem, error)
		// Laporan Coffee (tahap 4.10).
		GetSalesBreakdown(req *dto.CoffeeReportRequest) (*dto.SalesBreakdownResponse, error)
		GetIngredientReport(req *dto.CoffeeReportRequest) (*dto.IngredientReportResponse, error)
		GetWasteReport(req *dto.CoffeeReportRequest) (*dto.WasteReportResponse, error)
		GetCashierList(req *dto.CashierReportRequest) ([]dto.CashierItem, int64, error)
		ExportCashierReport(params dto.FilterParams) (*bytes.Buffer, error)
	}

	reportService struct {
		repo repo.ReportRepoInterface
	}
)

func NewReportService(r repo.ReportRepoInterface) *reportService {
	return &reportService{repo: r}
}
