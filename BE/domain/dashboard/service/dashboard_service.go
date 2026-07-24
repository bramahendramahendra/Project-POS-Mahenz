package service

import (
	dto "pos_api/domain/dashboard/dto"
	"time"
)

func (s *dashboardService) GetRecentTransactions(userID int, limit int) ([]dto.RecentTransactionItem, error) {
	items, err := s.repo.GetRecentTransactions(userID, limit)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.RecentTransactionItem{}
	}
	return items, nil
}

func (s *dashboardService) GetTodaySummary(userID int) (*dto.TodaySummaryResponse, error) {
	today := time.Now().Format("2006-01-02")
	return s.repo.GetTodaySummary(userID, today)
}
