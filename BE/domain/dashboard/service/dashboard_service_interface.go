package service

import (
	"pos_api/domain/dashboard/dto"
	repo "pos_api/domain/dashboard/repo"
)

type (
	DashboardServiceInterface interface {
		GetRecentTransactions(userID int, limit int) ([]dto.RecentTransactionItem, error)
		GetTodaySummary(userID int) (*dto.TodaySummaryResponse, error)
	}

	dashboardService struct {
		repo repo.DashboardRepoInterface
	}
)

func NewDashboardService(r repo.DashboardRepoInterface) *dashboardService {
	return &dashboardService{repo: r}
}
