package service

import repo "pos_api/domain/pin/repo"

type (
	PinServiceInterface interface {
		HasPin(userID int) (bool, error)
		SetPin(userID int, pin string) error
		VerifyPin(userID int, pin string) (bool, error)
		ChangePin(userID int, oldPin, newPin string) error
	}

	pinService struct {
		repo repo.PinRepoInterface
	}
)

func NewPinService(r repo.PinRepoInterface) *pinService {
	return &pinService{repo: r}
}
