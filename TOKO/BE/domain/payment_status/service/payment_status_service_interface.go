package service

import (
	dto "pos_api/domain/payment_status/dto"
	repo "pos_api/domain/payment_status/repo"
)

type (
	PaymentStatusServiceInterface interface {
		GetAll() ([]*dto.PaymentStatusResponse, error)
	}

	paymentStatusService struct {
		repo repo.PaymentStatusRepoInterface
	}
)

func NewPaymentStatusService(repo repo.PaymentStatusRepoInterface) *paymentStatusService {
	return &paymentStatusService{repo: repo}
}
