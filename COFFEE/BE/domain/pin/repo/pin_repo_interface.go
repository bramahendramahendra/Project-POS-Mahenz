package repo

import "gorm.io/gorm"

type (
	PinRepoInterface interface {
		GetPinHash(userID int) (string, error)
		SetPinHash(userID int, pinHash string) error
		ClearPinHash(userID int) error
	}

	pinRepo struct {
		db *gorm.DB
	}
)

func NewPinRepo(db *gorm.DB) *pinRepo {
	return &pinRepo{db: db}
}
