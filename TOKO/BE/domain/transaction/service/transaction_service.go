package service

import (
	stderrors "errors"

	cash_drawer_model "pos_api/domain/cash_drawer/model"
	"pos_api/domain/transaction/dto"
	time_helper "pos_api/helper/time"
	"pos_api/pkg/pricing"

	"pos_api/errors"

	"gorm.io/gorm"
)

func (s *transactionService) GetAll(req *dto.GetAllRequest) ([]*dto.TransactionResponse, int64, error) {
	transactions, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &errors.InternalServerError{Message: err.Error()}
	}
	return transactions, total, nil
}

func (s *transactionService) GetByID(id int, requestingUserID int, role string) (*dto.TransactionResponse, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if t == nil {
		return nil, &errors.NotFoundError{Message: "Transaksi tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && t.UserID != requestingUserID {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke transaksi ini"}
	}
	return t, nil
}

func (s *transactionService) Create(req *dto.CreateTransactionRequest, userID int) (*dto.CreateTransactionResponse, error) {
	openCashDrawer, err := s.cashDrawerRepo.GetOpenCashDrawer(userID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	if req.ShiftID != nil {
		if openCashDrawer == nil || openCashDrawer.ShiftID == nil || *openCashDrawer.ShiftID != *req.ShiftID {
			return nil, &errors.BadRequestError{Message: "Shift tidak sesuai dengan sesi kasir yang sedang aktif"}
		}
	}

	if err := s.recalculateTotals(req); err != nil {
		return nil, err
	}

	// Validasi balance_used
	if req.BalanceUsed > 0 {
		if req.CustomerID == nil {
			return nil, &errors.BadRequestError{Message: "Pilih pelanggan untuk menggunakan saldo"}
		}
		if req.BalanceUsed > req.TotalAmount {
			return nil, &errors.BadRequestError{Message: "Saldo yang digunakan melebihi total belanja"}
		}
	}

	// Hitung effective total (setelah saldo)
	effectiveTotal := req.TotalAmount - req.BalanceUsed

	// Validasi pembayaran (hanya cek sisa setelah saldo)
	if !req.IsCredit && effectiveTotal > 0 && req.PaymentAmount < effectiveTotal {
		return nil, &errors.BadRequestError{Message: "Jumlah pembayaran kurang dari total transaksi"}
	}
	if effectiveTotal > 0 && req.PaymentAmount > effectiveTotal {
		req.ChangeAmount = req.PaymentAmount - effectiveTotal
	} else {
		req.ChangeAmount = 0
	}

	var resp *dto.CreateTransactionResponse
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txnRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		result, err := txnRepo.Create(req, userID)
		if err != nil {
			return err
		}
		resp = result

		// Deduct customer balance SETELAH transaksi berhasil dibuat
		if req.BalanceUsed > 0 {
			cbRepo := s.customerBalanceRepo.WithTx(tx)
			refID := resp.ID
			_, err := cbRepo.Deduct(*req.CustomerID, req.BalanceUsed, "transaction", &refID, "Bayar belanja dari saldo", userID)
			if err != nil {
				return &errors.BadRequestError{Message: "Saldo pelanggan tidak mencukupi"}
			}
		}

		// Update kas harian (hanya cash, hanya effective_total)
		if req.PaymentMethod == "cash" {
			netCash := req.TotalAmount - req.BalanceUsed
			if netCash > 0 {
				drawer, err := cashDrawerRepo.GetOpenCashDrawer(userID)
				if err != nil {
					return err
				}
				if drawer != nil {
					if err := cashDrawerRepo.UpdateSales(drawer.ID, netCash, netCash, time_helper.GetTimeNow()); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if txErr != nil {
		var badReq *errors.BadRequestError
		if stderrors.As(txErr, &badReq) {
			return nil, badReq
		}
		return nil, &errors.InternalServerError{Message: txErr.Error()}
	}

	return resp, nil
}

func (s *transactionService) recalculateTotals(req *dto.CreateTransactionRequest) error {
	items := make([]pricing.Item, len(req.Items))
	for i, item := range req.Items {
		items[i] = pricing.Item{
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			UnitID:       item.UnitID,
			Quantity:     item.Quantity,
			DiscountItem: item.DiscountItem,
		}
	}

	totals, err := pricing.Recalculate(s.productRepo, items, req.Discount, req.Tax)
	if err != nil {
		return err
	}

	for i := range req.Items {
		req.Items[i].Price = totals.ItemPrices[i]
		req.Items[i].Subtotal = totals.ItemSubtotals[i]
	}
	req.Subtotal = totals.Subtotal
	req.TotalAmount = totals.TotalAmount
	return nil
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

	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txnRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		if err := txnRepo.Void(req.ID, userID); err != nil {
			return err
		}

		// Rollback kas harian
		if t.PaymentMethod == "cash" {
			netCash := t.TotalAmount - t.BalanceUsed
			if netCash > 0 {
				var drawer *cash_drawer_model.CashDrawer
				if t.CashDrawerID != nil {
					// Transaksi punya referensi kas langsung → pakai itu
					d, dErr := cashDrawerRepo.GetByID(*t.CashDrawerID)
					if dErr != nil {
						return dErr
					}
					// Hanya rollback jika kas masih open
					if d != nil && d.Status == "open" {
						drawer = d
					}
				} else {
					// Transaksi lama (sebelum fitur backdate) → pakai kas open biasa
					d, dErr := cashDrawerRepo.GetOpenCashDrawer(t.UserID)
					if dErr != nil {
						return dErr
					}
					drawer = d
				}
				if drawer != nil {
					if err := cashDrawerRepo.UpdateSales(drawer.ID, -netCash, -netCash, time_helper.GetTimeNow()); err != nil {
						return err
					}
				}
			}
		}

		// Rollback saldo pelanggan (jika transaksi pakai saldo)
		if t.BalanceUsed > 0 && t.CustomerID != nil {
			cbRepo := s.customerBalanceRepo.WithTx(tx)
			refID := req.ID
			_, err := cbRepo.Topup(*t.CustomerID, t.BalanceUsed, "void", &refID, "Void transaksi — saldo dikembalikan", userID)
			if err != nil {
				return err
			}
		}

		// Rollback top-up saldo (jika kembalian pernah disimpan ke saldo)
		if t.CustomerID != nil {
			cbRepo := s.customerBalanceRepo.WithTx(tx)
			topupMutation, err := cbRepo.GetMutationByRef("transaction", req.ID, "topup")
			if err != nil {
				return err
			}
			if topupMutation != nil {
				refID := req.ID
				_, err := cbRepo.Deduct(*t.CustomerID, topupMutation.Amount, "void", &refID, "Void transaksi — top-up saldo dibatalkan", userID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if txErr != nil {
		var badReq *errors.BadRequestError
		if stderrors.As(txErr, &badReq) {
			return badReq
		}
		return &errors.InternalServerError{Message: txErr.Error()}
	}
	return nil
}
