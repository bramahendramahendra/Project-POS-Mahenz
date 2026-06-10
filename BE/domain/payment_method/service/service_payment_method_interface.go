package service

import (
	dto "pos_api/domain/payment_method/dto"
	repo "pos_api/domain/payment_method/repo"
)

type (
	PaymentMethodServiceInterface interface {
		GetAll() ([]*dto.PaymentMethodResponse, error)
	}

	paymentMethodService struct {
		repo repo.PaymentMethodRepoInterface
	}
)

func NewPaymentMethodService(repo repo.PaymentMethodRepoInterface) *paymentMethodService {
	return &paymentMethodService{repo: repo}
}
