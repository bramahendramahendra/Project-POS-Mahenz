package repo

import (
	"pos_api/domain/dashboard/dto"

	"gorm.io/gorm"
)

type (
	DashboardRepoInterface interface {
		GetRecentTransactions(userID int, limit int) ([]dto.RecentTransactionItem, error)
		GetTodaySummary(userID int, date string) (*dto.TodaySummaryResponse, error)
	}

	dashboardRepo struct {
		db *gorm.DB
	}
)

func NewDashboardRepo(db *gorm.DB) *dashboardRepo {
	return &dashboardRepo{db: db}
}
