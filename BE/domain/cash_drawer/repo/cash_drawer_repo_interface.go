package repo

import (
	dto "pos_api/domain/cash_drawer/dto"
	model "pos_api/domain/cash_drawer/model"

	"gorm.io/gorm"
)

type (
	CashDrawerRepoInterface interface {
		GetCurrent(userID int) (*dto.CurrentCashDrawerResponse, error)
		GetOpenCashDrawer(userID int) (*model.CashDrawer, error)
		GetByID(id int) (*model.CashDrawer, error)
		GetDetailByID(id int) (*model.CashDrawerDetail, error)
		GetHistory(req *dto.GetHistoryRequest) ([]*dto.CashDrawerHistoryResponse, int64, error)
		Open(userID int, shiftID *int, openingBalance float64, notes string) (int64, error)
		Close(id int, closingBalance, expectedBalance, difference float64, notes string) error
		UpdateSales(id int, totalSales, totalCashSales float64) error
		UpdateExpenses(id int, totalExpenses float64) error

		GetDB() *gorm.DB
	}

	cashDrawerRepo struct {
		db *gorm.DB
	}
)

func NewCashDrawerRepo(db *gorm.DB) *cashDrawerRepo {
	return &cashDrawerRepo{db: db}
}

func (r *cashDrawerRepo) GetDB() *gorm.DB {
	return r.db
}
