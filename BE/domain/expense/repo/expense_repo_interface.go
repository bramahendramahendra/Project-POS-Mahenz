package repo

import (
	dto "pos_api/domain/expense/dto"
	model "pos_api/domain/expense/model"

	"gorm.io/gorm"
)

type (
	ExpenseRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Expense, int64, error)
		GetByID(id int) (*model.Expense, error)
		Create(req *dto.CreateRequest, userID int) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		UpdateFromSync(id int, data map[string]interface{}) error

		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) ExpenseRepoInterface
	}

	expenseRepo struct {
		db *gorm.DB
	}
)

func NewExpenseRepo(db *gorm.DB) *expenseRepo {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) GetDB() *gorm.DB {
	return r.db
}

// WithTx mengembalikan repo instance baru yang terikat ke transaksi tx, supaya operasi
// expense bisa digabung dalam satu DB transaction bersama repo lain (mis. cash_drawer).
func (r *expenseRepo) WithTx(tx *gorm.DB) ExpenseRepoInterface {
	return &expenseRepo{db: tx}
}
