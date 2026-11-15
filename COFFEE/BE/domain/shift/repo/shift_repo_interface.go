package repo

import (
	dto "coffee_api/domain/shift/dto"
	model "coffee_api/domain/shift/model"

	"gorm.io/gorm"
)

type (
	ShiftRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Shift, int64, error)
		GetOptions() ([]*model.Shift, error)
		GetByID(id int) (*model.Shift, error)
		CountOpenCashDrawer(shiftID int) (int, error)
		Create(req *dto.CreateRequest) (int64, error)
		Update(req *dto.UpdateRequest) error
		Delete(req *dto.DeleteRequest) error
		ToggleStatus(req *dto.ToggleStatusRequest) error
		GetSummary() ([]*dto.ShiftSummaryResponse, error)
	}

	shiftRepo struct {
		db *gorm.DB
	}
)

func NewShiftRepo(db *gorm.DB) *shiftRepo {
	return &shiftRepo{db: db}
}
