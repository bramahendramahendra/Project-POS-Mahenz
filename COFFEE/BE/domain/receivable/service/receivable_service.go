package service

import (
	"fmt"

	dto "coffee_api/domain/receivable/dto"
	"coffee_api/errors"
)

func (s *receivableService) GetAll(req *dto.GetAllRequest) ([]*dto.ReceivableResponse, int64, error) {
	return s.repo.GetAll(req)
}

func (s *receivableService) GetByID(id int) (*dto.ReceivableDetailResponse, error) {
	detail, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, &errors.NotFoundError{Message: "Piutang tidak ditemukan"}
	}
	return detail, nil
}

func (s *receivableService) GetStats() (*dto.ReceivableStats, error) {
	return s.repo.GetStats()
}

func (s *receivableService) Pay(req *dto.PayRequest) (*dto.PayResponse, error) {
	rec, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, &errors.NotFoundError{Message: "Piutang tidak ditemukan"}
	}
	if rec.Status == "paid" {
		return nil, &errors.BadRequestError{Message: "Piutang sudah lunas"}
	}
	if rec.Status == "void" {
		return nil, &errors.BadRequestError{Message: "Piutang ini sudah dibatalkan (transaksi induk di-void)"}
	}
	if req.Amount > rec.RemainingAmount {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("Jumlah bayar (%.0f) melebihi sisa piutang (%.0f)", req.Amount, rec.RemainingAmount),
		}
	}

	if err := s.repo.CreatePayment(req.ID, req, req.UserID); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAfterPayment(req.ID, req.Amount); err != nil {
		return nil, err
	}

	newPaid := rec.PaidAmount + req.Amount
	newRemaining := rec.RemainingAmount - req.Amount
	status := "partial"
	if newRemaining <= 0 {
		status = "paid"
	}

	return &dto.PayResponse{
		ReceivableID:    req.ID,
		PaidAmount:      newPaid,
		RemainingAmount: newRemaining,
		Status:          status,
	}, nil
}
