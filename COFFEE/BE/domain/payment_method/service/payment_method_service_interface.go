package service

import (
	dto "coffee_api/domain/payment_method/dto"
	repo "coffee_api/domain/payment_method/repo"
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
