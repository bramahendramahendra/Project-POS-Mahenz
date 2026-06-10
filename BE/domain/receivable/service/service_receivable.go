package service

import (
	"fmt"

	dto "pos_api/domain/receivable/dto"
	"pos_api/errors"
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

func (s *receivableService) GetSummary() ([]*dto.ReceivableSummaryItem, error) {
	return s.repo.GetSummary()
}

func (s *receivableService) GetPayments(id int) ([]*dto.PaymentResponse, error) {
	rec, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, &errors.NotFoundError{Message: "Piutang tidak ditemukan"}
	}
	return s.repo.GetPayments(id)
}

func (s *receivableService) Pay(id int, req *dto.PayRequest, userID int) (*dto.PayResponse, error) {
	rec, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, &errors.NotFoundError{Message: "Piutang tidak ditemukan"}
	}
	if rec.Status == "paid" {
		return nil, &errors.BadRequestError{Message: "Piutang sudah lunas"}
	}
	if req.Amount > rec.RemainingAmount {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("Jumlah bayar (%.0f) melebihi sisa piutang (%.0f)", req.Amount, rec.RemainingAmount),
		}
	}

	if err := s.repo.CreatePayment(id, req, userID); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAfterPayment(id, req.Amount); err != nil {
		return nil, err
	}

	newPaid := rec.PaidAmount + req.Amount
	newRemaining := rec.RemainingAmount - req.Amount
	status := "partial"
	if newRemaining <= 0 {
		status = "paid"
	}

	return &dto.PayResponse{
		ReceivableID:    id,
		PaidAmount:      newPaid,
		RemainingAmount: newRemaining,
		Status:          status,
	}, nil
}
