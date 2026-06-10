package service

import (
	dto "pos_api/domain/cash_drawer/dto"
	repo "pos_api/domain/cash_drawer/repo"
)

type (
	CashDrawerServiceInterface interface {
		GetCurrent(userID int) (*dto.CurrentCashDrawerResponse, error)
		GetByID(id int, requestingUserID int, role string) (*dto.CashDrawerDetailResponse, error)
		GetHistory(req *dto.GetHistoryRequest) (data []*dto.CashDrawerHistoryResponse, total int64, err error)
		Open(userID int, req *dto.OpenRequest) (*dto.OpenResponse, error)
		Close(id int, req *dto.CloseRequest, requestingUserID int, role string) (*dto.CloseResponse, error)
		UpdateSales(id int, req *dto.UpdateSalesRequest, requestingUserID int, role string) error
		UpdateExpenses(id int, req *dto.UpdateExpensesRequest, requestingUserID int, role string) error
	}

	cashDrawerService struct {
		repo repo.CashDrawerRepoInterface
	}
)

func NewCashDrawerService(repo repo.CashDrawerRepoInterface) *cashDrawerService {
	return &cashDrawerService{repo: repo}
}
