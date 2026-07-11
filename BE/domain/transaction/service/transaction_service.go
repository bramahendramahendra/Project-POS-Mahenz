package service

import (
	"strings"

	"pos_api/domain/transaction/dto"
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

	// Transaksi kredit (piutang) boleh dibayar sebagian/tidak sama sekali di muka —
	// sisanya tercatat sebagai piutang. Transaksi non-kredit wajib lunas saat itu juga.
	if !req.IsCredit && req.PaymentAmount < req.TotalAmount {
		return nil, &errors.BadRequestError{Message: "Jumlah pembayaran kurang dari total transaksi"}
	}
	if change := req.PaymentAmount - req.TotalAmount; change > 0 {
		req.ChangeAmount = change
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

		if req.PaymentMethod == "cash" {
			drawer, err := cashDrawerRepo.GetOpenCashDrawer(userID)
			if err != nil {
				return err
			}
			if drawer != nil {
				if err := cashDrawerRepo.UpdateSales(drawer.ID, req.TotalAmount, req.TotalAmount); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if txErr != nil {
		if strings.HasPrefix(txErr.Error(), "stok_insufficient:") {
			name := strings.TrimPrefix(txErr.Error(), "stok_insufficient:")
			return nil, &errors.BadRequestError{Message: "Stok tidak mencukupi untuk " + name}
		}
		return nil, &errors.InternalServerError{Message: txErr.Error()}
	}

	return resp, nil
}

// recalculateTotals menghitung ulang subtotal/diskon/pajak/total transaksi di server
// berdasarkan harga produk asli, mengabaikan nilai price/subtotal/total_amount dari client.
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

		if t.PaymentMethod == "cash" {
			drawer, err := cashDrawerRepo.GetOpenCashDrawer(t.UserID)
			if err != nil {
				return err
			}
			if drawer != nil {
				if err := cashDrawerRepo.UpdateSales(drawer.ID, -t.TotalAmount, -t.TotalAmount); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if txErr != nil {
		return &errors.InternalServerError{Message: txErr.Error()}
	}
	return nil
}
