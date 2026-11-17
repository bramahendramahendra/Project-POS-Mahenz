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
	}

	businessSummaryService struct {
		repo repo.BusinessSummaryRepoInterface
	}
)

func NewBusinessSummaryService(r repo.BusinessSummaryRepoInterface) *businessSummaryService {
	return &businessSummaryService{repo: r}
}
