package service

import (
	"strings"

	"pos_api/domain/transaction/dto"
	
	"pos_api/errors"
)

func (s *transactionService) GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error) {
	transactions, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &errors.InternalServerError{Message: err.Error()}
	}
	return transactions, total, nil
}

func (s *transactionService) GetByID(id int) (*dto.TransactionResponse, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if t == nil {
		return nil, &errors.NotFoundError{Message: "Transaksi tidak ditemukan"}
	}
	return t, nil
}

func (s *transactionService) Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	resp, err := s.repo.Create(req, userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "stok_insufficient:") {
			name := strings.TrimPrefix(err.Error(), "stok_insufficient:")
			return nil, &errors.BadRequestError{Message: "Stok tidak mencukupi untuk " + name}
		}
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	if req.PaymentMethod == "cash" {
		openCashDrawer, _ := s.cashDrawerRepo.GetOpenCashDrawer(userID)
		if openCashDrawer != nil {
			_ = s.cashDrawerRepo.UpdateSales(openCashDrawer.ID, resp.TotalAmount, resp.TotalAmount)
		}
	}

	return resp, nil
}

func (s *transactionService) Void(req *dto.VoidRequest, userID int) error {
	t, err := s.repo.GetByID(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if t == nil {
		return &errors.NotFoundError{Message: "Transaksi tidak ditemukan"}
	}
	if t.Status == "void" {
		return &errors.BadRequestError{Message: "Transaksi sudah di-void"}
	}

	if err := s.repo.Void(req.ID, userID); err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	return nil
}

