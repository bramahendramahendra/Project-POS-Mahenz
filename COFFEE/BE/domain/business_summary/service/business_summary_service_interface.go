package service

import (
	"coffee_api/domain/business_summary/dto"
	repo "coffee_api/domain/business_summary/repo"
)

type (
	BusinessSummaryServiceInterface interface {
		GetStats(period string) (*dto.StatsResponse, error)
		GetSalesTrend(period string) ([]dto.SalesTrendItem, error)
		GetTopProducts(filter dto.DateRangeFilter) ([]dto.TopProductItem, error)
		GetTopCategories(filter dto.DateRangeFilter) ([]dto.TopCategoryItem, error)
		GetPaymentMethods(filter dto.DateRangeFilter) ([]dto.PaymentMethodItem, error)
		GetOrderTypes(filter dto.DateRangeFilter) ([]dto.OrderTypeItem, error)
	}

	businessSummaryService struct {
		repo repo.BusinessSummaryRepoInterface
	}
)

func NewBusinessSummaryService(r repo.BusinessSummaryRepoInterface) *businessSummaryService {
	return &businessSummaryService{repo: r}
}
