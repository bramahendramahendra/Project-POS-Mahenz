package service

import (
	stderrors "errors"

	cash_drawer_model "coffee_api/domain/cash_drawer/model"
	"coffee_api/domain/transaction/dto"
	time_helper "coffee_api/helper/time"
	"coffee_api/pkg/pricing"

	"coffee_api/errors"

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

	// Coffee tidak memakai saldo pelanggan (modul customer_balance dipangkas, keputusan 26).
	if req.BalanceUsed > 0 {
		return nil, &errors.BadRequestError{Message: "Pembayaran dengan saldo pelanggan tidak tersedia"}
	}

	effectiveTotal := req.TotalAmount

	// Validasi pembayaran
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

		// Update kas harian (hanya cash)
		if req.PaymentMethod == "cash" {
			netCash := req.TotalAmount
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
			netCash := t.TotalAmount
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
					// Transaksi tanpa referensi kas → pakai kas open milik kasirnya
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
