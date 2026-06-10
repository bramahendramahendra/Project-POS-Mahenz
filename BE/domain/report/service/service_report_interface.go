package service

import (
	"bytes"

	"pos_api/domain/report/dto"
	repo "pos_api/domain/report/repo"
)

type (
	ReportServiceInterface interface {
		GetSalesReport(params dto.FilterParams) (*dto.SalesReportResponse, error)
		GetSalesChart(params dto.FilterParams) ([]dto.SalesChartItem, error)
		ExportSalesReport(params dto.FilterParams) (*bytes.Buffer, error)

		GetProfitLoss(params dto.FilterParams) (*dto.ProfitLossResponse, error)
		ExportProfitLoss(params dto.FilterParams) (*bytes.Buffer, error)

		GetStockReport() (*dto.StockReportResponse, error)
		ExportStockReport() (*bytes.Buffer, error)

		GetCashierReport(params dto.FilterParams) ([]dto.CashierItem, error)
		ExportCashierReport(params dto.FilterParams) (*bytes.Buffer, error)
	}

	reportService struct {
		repo repo.ReportRepoInterface
	}
)

func NewReportService(r repo.ReportRepoInterface) *reportService {
	return &reportService{repo: r}
}
