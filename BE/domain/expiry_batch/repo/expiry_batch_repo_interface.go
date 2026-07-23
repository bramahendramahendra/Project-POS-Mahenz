package repo

import (
	model "pos_api/domain/expiry_batch/model"

	"gorm.io/gorm"
)

type (
	ExpiryBatchRepoInterface interface {
		// GetWarnings: semua batch status='active' yang expired_date-nya sudah mendekati
		// (dalam nearExpiryDays hari ke depan) atau sudah lewat hari ini.
		GetWarnings(search string) ([]*model.ExpiryBatch, error)
		GetByID(id int) (*model.ExpiryBatch, error)
		Confirm(id, userID int, notes string) error
		WriteOff(id, userID int, notes string) error

		GetDB() *gorm.DB
	}

	expiryBatchRepo struct {
		db *gorm.DB
	}
)

func NewExpiryBatchRepo(db *gorm.DB) *expiryBatchRepo {
	return &expiryBatchRepo{db: db}
}

func (r *expiryBatchRepo) GetDB() *gorm.DB {
	return r.db
}
