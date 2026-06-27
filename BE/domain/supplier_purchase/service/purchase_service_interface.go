package service

import (
	dto "pos_api/domain/supplier_purchase/dto"
	repo "pos_api/domain/supplier_purchase/repo"
)

type (
	PurchaseServiceInterface interface {
		GetAll(req *dto.GetAllRequest) (data []dto.PurchaseResponse, total int64, err error)
		GetByID(id int) (data dto.PurchaseResponse, err error)
		GenerateCode() (data dto.GenerateCodeResponse, err error)
		Create(req *dto.CreateRequest) (data dto.PurchaseResponse, err error)
		Update(req *dto.UpdateRequest) (data dto.PurchaseResponse, err error)
		Delete(id int) error
		Pay(req *dto.PayRequest) error
		GetPayments(purchaseID int) (data []*dto.PaymentResponse, err error)
	}

	purchaseService struct {
		repo repo.PurchaseRepoInterface
	}
)

func NewPurchaseService(r repo.PurchaseRepoInterface) *purchaseService {
	return &purchaseService{repo: r}
}
