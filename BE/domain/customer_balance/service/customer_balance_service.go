package service

import (
	"fmt"

	"pos_api/domain/customer_balance/dto"
	"pos_api/errors"

	"gorm.io/gorm"
)

func (s *customerBalanceService) Topup(req *dto.TopupRequest) (*dto.BalanceResponse, error) {
	// Validasi customer exists
	customer, err := s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	var newBalance float64
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		bal, err := txRepo.Topup(req.CustomerID, req.Amount, "manual", nil, req.Notes, req.UserID)
		if err != nil {
			return err
		}
		newBalance = bal
		return nil
	})
	if txErr != nil {
		return nil, &errors.InternalServerError{Message: txErr.Error()}
	}

	return &dto.BalanceResponse{
		CustomerID:   req.CustomerID,
		Balance:      newBalance,
		BalanceAfter: newBalance,
	}, nil
}

func (s *customerBalanceService) Refund(req *dto.RefundRequest) (*dto.BalanceResponse, error) {
	// Validasi customer exists
	customer, err := s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	// Validasi saldo cukup
	if customer.Balance < req.Amount {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("Saldo tidak cukup (saldo: %.0f, refund: %.0f)", customer.Balance, req.Amount),
		}
	}

	var newBalance float64
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		bal, err := txRepo.Refund(req.CustomerID, req.Amount, "manual", nil, req.Notes, req.UserID)
		if err != nil {
			return err
		}
		newBalance = bal
		return nil
	})
	if txErr != nil {
		return nil, &errors.InternalServerError{Message: txErr.Error()}
	}

	return &dto.BalanceResponse{
		CustomerID:   req.CustomerID,
		Balance:      newBalance,
		BalanceAfter: newBalance,
	}, nil
}

func (s *customerBalanceService) GetHistory(req *dto.HistoryRequest) ([]*dto.MutationResponse, int64, error) {
	// Validasi customer exists
	customer, err := s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		return nil, 0, err
	}
	if customer == nil {
		return nil, 0, &errors.NotFoundError{Message: "Pelanggan tidak ditemukan"}
	}

	return s.repo.GetHistory(req)
}

func (s *customerBalanceService) SaveToBalance(req *dto.SaveToBalanceRequest) (*dto.BalanceResponse, error) {
	// Cek idempotency: apakah sudah pernah top-up untuk transaksi ini?
	existing, err := s.repo.GetMutationByRef("transaction", req.TransactionID, "topup")
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Kembalian untuk transaksi ini sudah disimpan ke saldo sebelumnya"}
	}

	// TODO: validasi transaction_id exists dan payment_method=cash dan change_amount >= req.Amount
	// Untuk saat ini, trust FE sudah mengirim data yang benar

	var newBalance float64
	var customerID int

	// Ambil customer_id dari transaksi
	var txData struct {
		CustomerID *int
		ChangeAmt  float64
	}
	if err := s.repo.GetDB().Raw(
		`SELECT customer_id, change_amount FROM transactions WHERE id = ? AND status = 'completed' LIMIT 1`,
		req.TransactionID,
	).Scan(&txData).Error; err != nil {
		return nil, err
	}
	if txData.CustomerID == nil {
		return nil, &errors.BadRequestError{Message: "Transaksi tidak memiliki pelanggan"}
	}
	if req.Amount > txData.ChangeAmt {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("Jumlah simpan (%.0f) melebihi kembalian (%.0f)", req.Amount, txData.ChangeAmt),
		}
	}
	customerID = *txData.CustomerID

	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		refID := req.TransactionID
		bal, err := txRepo.Topup(customerID, req.Amount, "transaction", &refID, "Simpan kembalian ke saldo", req.UserID)
		if err != nil {
			return err
		}
		newBalance = bal

		// Update cash drawer: kembalian tidak jadi keluar laci
		// Cari cash drawer aktif user ini dan tambah total_cash_sales
		err = tx.Exec(
			`UPDATE cash_drawer SET total_cash_sales = total_cash_sales + ?, expected_balance = expected_balance + ?, updated_at = NOW() WHERE user_id = (SELECT user_id FROM transactions WHERE id = ?) AND status = 'open'`,
			req.Amount, req.Amount, req.TransactionID,
		).Error
		return err
	})
	if txErr != nil {
		return nil, &errors.InternalServerError{Message: txErr.Error()}
	}

	return &dto.BalanceResponse{
		CustomerID:   customerID,
		Balance:      newBalance,
		BalanceAfter: newBalance,
	}, nil
}
