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
		GetDetailByID(id int) (*model.CashDrawerDetail, []model.CashDrawerTransactionItem, []model.CashDrawerExpenseItem, error)
		GetHistory(req *dto.GetHistoryRequest) ([]*dto.CashDrawerHistoryResponse, int64, error)
		Open(userID int, shiftID *int, openingBalance float64, notes string) (int64, error)
		Close(id int, closingBalance, expectedBalance, difference float64, notes string) error
		UpdateSales(id int, totalSales, totalCashSales float64) error
		UpdateExpenses(id int, totalExpenses float64) error
		GetMyCash(userID int) (*model.CashDrawerDetail, []model.CashDrawerTransactionItem, []model.CashDrawerExpenseItem, error)
		GetNonCashSales(userID int, openTime string, closeTime *string) ([]dto.NonCashSaleItem, error)
		GetNonCashTransactions(userID int, openTime string, closeTime *string, nextOpenTime *string) ([]model.CashDrawerNonCashTransactionItem, error)
		AutoCloseYesterday() (int, error)
		GetSummary(req *dto.GetHistoryRequest) (*dto.CashDrawerSummaryResponse, error)
		GetKasirOptions() ([]dto.KasirOptionResponse, error)

		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) CashDrawerRepoInterface
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

// WithTx mengembalikan repo instance baru yang terikat ke transaksi tx, supaya operasi
// cash_drawer bisa digabung dalam satu DB transaction bersama repo lain (mis. expense).
func (r *cashDrawerRepo) WithTx(tx *gorm.DB) CashDrawerRepoInterface {
	return &cashDrawerRepo{db: tx}
}
