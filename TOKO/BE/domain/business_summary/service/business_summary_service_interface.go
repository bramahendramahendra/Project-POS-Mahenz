package service

import (
	"pos_api/domain/business_summary/dto"
	repo "pos_api/domain/business_summary/repo"
)

type (
	BusinessSummaryServiceInterface interface {
		GetStats(period string) (*dto.StatsResponse, error)
		GetSalesTrend(period string) ([]dto.SalesTrendItem, error)
		GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error)
		GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error)
		GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error)
		GetSummaryExtra(period string) (*dto.SummaryExtraResponse, error)
	}

	businessSummaryService struct {
		repo repo.BusinessSummaryRepoInterface
	}
)

func NewBusinessSummaryService(r repo.BusinessSummaryRepoInterface) *businessSummaryService {
	return &businessSummaryService{repo: r}
}
