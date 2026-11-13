package repo

import (
	dto "coffee_api/domain/payment_method/dto"

	"gorm.io/gorm"
)

type (
	PaymentMethodRepoInterface interface {
		GetAll() ([]*dto.PaymentMethodResponse, error)
	}

	paymentMethodRepo struct {
		db *gorm.DB
	}
)

func NewPaymentMethodRepo(db *gorm.DB) *paymentMethodRepo {
	return &paymentMethodRepo{db: db}
}
