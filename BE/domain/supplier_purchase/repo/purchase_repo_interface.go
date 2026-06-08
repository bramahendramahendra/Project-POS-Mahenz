package repo

import (
	dto_purchase "pos_api/domain/supplier_purchase/dto"
	model_purchase "pos_api/domain/supplier_purchase/model"
)

type PurchaseRepo interface {
	GetAll(req *dto_purchase.PurchaseListRequest) ([]*dto_purchase.PurchaseResponse, int, error)
	GetByID(id int) (*dto_purchase.PurchaseResponse, error)
	GetItems(purchaseID int) ([]model_purchase.PurchaseItem, error)
	GetPayments(purchaseID int) ([]dto_purchase.PurchasePaymentResponse, error)
	GenerateCode() (string, error)
	Create(req *dto_purchase.PurchaseRequest) (*dto_purchase.PurchaseResponse, error)
	Update(req *dto_purchase.PurchaseRequest) (*dto_purchase.PurchaseResponse, error)
	Delete(id int) error
	Pay(req *dto_purchase.PayPurchaseRequest) error
	GetRawByID(id int) (*model_purchase.Purchase, error)
	IsValidPaymentMethod(code string) (bool, error)
}
