package service

import (
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	dto "pos_api/domain/expense/dto"
	repo "pos_api/domain/expense/repo"
)

type (
	ExpenseServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.ExpenseResponse, total int64, err error)
		GetByID(id int) (data dto.ExpenseResponse, err error)
		Create(req *dto.CreateRequest, userID int) (data dto.ExpenseResponse, err error)
		Update(req *dto.UpdateRequest, requestingUserID int, role string) (err error)
		Delete(req *dto.DeleteRequest, requestingUserID int, role string) (err error)
	}

	expenseService struct {
		repo           repo.ExpenseRepoInterface
		cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface
	}
)

func NewExpenseService(repo repo.ExpenseRepoInterface, cashDrawerRepo cash_drawer_repo.CashDrawerRepoInterface) *expenseService {
	return &expenseService{repo: repo, cashDrawerRepo: cashDrawerRepo}
}
