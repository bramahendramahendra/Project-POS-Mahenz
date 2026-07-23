package service

import (
	dto "pos_api/domain/expiry_batch/dto"
	repo "pos_api/domain/expiry_batch/repo"
)

type (
	ExpiryBatchServiceInterface interface {
		GetWarnings(req *dto.GetWarningsRequest) (data []dto.WarningResponse, productSeverity []dto.ProductSeverityResponse, err error)
		Confirm(req *dto.ConfirmRequest) (err error)
		WriteOff(req *dto.WriteOffRequest) (err error)
	}

	expiryBatchService struct {
		repo repo.ExpiryBatchRepoInterface
	}
)

func NewExpiryBatchService(repo repo.ExpiryBatchRepoInterface) *expiryBatchService {
	return &expiryBatchService{repo: repo}
}
