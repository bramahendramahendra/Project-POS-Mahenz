package service

import (
	dto "pos_api/domain/receivable/dto"
	repo "pos_api/domain/receivable/repo"
)

type (
	ReceivableServiceInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.ReceivableResponse, int64, error)
		GetByID(id int) (*dto.ReceivableDetailResponse, error)
		GetSummary() ([]*dto.ReceivableSummaryItem, error)
		GetPayments(id int) ([]*dto.PaymentResponse, error)
		Pay(req *dto.PayRequest) (*dto.PayResponse, error)
	}

	receivableService struct {
		repo repo.ReceivableRepoInterface
	}
)

func NewReceivableService(repo repo.ReceivableRepoInterface) *receivableService {
	return &receivableService{repo: repo}
}
