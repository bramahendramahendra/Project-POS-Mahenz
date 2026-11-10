package repo

import (
	dto "pos_api/domain/payment_status/dto"

	"gorm.io/gorm"
)

type (
	PaymentStatusRepoInterface interface {
		GetAll() ([]*dto.PaymentStatusResponse, error)
	}

	paymentStatusRepo struct {
		db *gorm.DB
	}
)

func NewPaymentStatusRepo(db *gorm.DB) *paymentStatusRepo {
	return &paymentStatusRepo{db: db}
}
