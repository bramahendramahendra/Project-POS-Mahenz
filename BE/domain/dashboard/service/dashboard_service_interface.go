package service

import (
	"pos_api/domain/dashboard/dto"
	repo "pos_api/domain/dashboard/repo"
)

type (
	DashboardServiceInterface interface {
		GetStats(date string) (*dto.StatsResponse, error)
		GetSalesTrend(period string) ([]dto.SalesTrendItem, error)
		GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error)
		GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error)
		GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error)
		GetSummaryExtra(period string) (*dto.SummaryExtraResponse, error)
	}

	dashboardService struct {
		repo repo.DashboardRepoInterface
	}
)

func NewDashboardService(r repo.DashboardRepoInterface) *dashboardService {
	return &dashboardService{repo: r}
}
