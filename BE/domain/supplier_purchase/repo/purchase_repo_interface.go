package repo

import (
	dto "pos_api/domain/supplier_purchase/dto"
	model "pos_api/domain/supplier_purchase/model"

	"gorm.io/gorm"
)

type (
	PurchaseRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.PurchaseRow, int64, error)
		GetByID(id int) (*model.PurchaseRow, error)
		GetPayments(purchaseID int) ([]model.PurchasePayment, error)
		GenerateCode() (string, error)
		Create(req *dto.CreateRequest) (*model.PurchaseRow, error)
		Update(req *dto.UpdateRequest) (*model.PurchaseRow, error)
		Delete(id int) error
		Pay(req *dto.PayRequest) error
		GetRawByID(id int) (*model.Purchase, error)
		IsValidPaymentMethod(code string) (bool, error)

		GetDB() *gorm.DB
	}

	purchaseRepo struct {
		db *gorm.DB
	}
)

func NewPurchaseRepo(db *gorm.DB) *purchaseRepo {
	return &purchaseRepo{db: db}
}

func (r *purchaseRepo) GetDB() *gorm.DB {
	return r.db
}
