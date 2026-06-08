package service

import dto_purchase "pos_api/domain/supplier_purchase/dto"

type PurchaseService interface {
	GetAll(req *dto_purchase.PurchaseListRequest) ([]*dto_purchase.PurchaseResponse, int, error)
	GetByID(id int) (*dto_purchase.PurchaseResponse, error)
	GetItems(purchaseID int) ([]dto_purchase.PurchaseItemResponse, error)
	GenerateCode() (*dto_purchase.GeneratePurchaseCodeResponse, error)
	Create(req *dto_purchase.PurchaseRequest) (*dto_purchase.PurchaseResponse, error)
	Update(req *dto_purchase.PurchaseRequest) (*dto_purchase.PurchaseResponse, error)
	Delete(id int) error
	Pay(req *dto_purchase.PayPurchaseRequest) error
	GetPayments(purchaseID int) ([]dto_purchase.PurchasePaymentResponse, error)
}
