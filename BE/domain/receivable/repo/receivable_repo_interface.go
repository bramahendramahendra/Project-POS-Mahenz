package repo

import (
	dto "pos_api/domain/receivable/dto"
	model "pos_api/domain/receivable/model"

	"gorm.io/gorm"
)

type (
	ReceivableRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*dto.ReceivableResponse, int64, error)
		GetByID(id int) (*model.Receivable, error)
		GetDetailByID(id int) (*dto.ReceivableDetailResponse, error)
		GetSummary() ([]*dto.ReceivableSummaryItem, error)
		GetPayments(receivableID int) ([]*dto.PaymentResponse, error)
		CreatePayment(receivableID int, req *dto.PayRequest, userID int) error
		UpdateAfterPayment(receivableID int, amount float64) error
	}

	receivableRepo struct {
		db *gorm.DB
	}
)

func NewReceivableRepo(db *gorm.DB) *receivableRepo {
	return &receivableRepo{db: db}
}
