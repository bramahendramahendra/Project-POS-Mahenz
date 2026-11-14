package service

import (
	stderrors "errors"
	"fmt"
	"math"
	"strings"

	"coffee_api/domain/transaction/dto"
	repo "coffee_api/domain/transaction/repo"
	time_helper "coffee_api/helper/time"

	"coffee_api/errors"

	"gorm.io/gorm"
)

// Open bill Coffee (tahap 4.8, keputusan 27-29, 40, 44, 59, 103-109). Semua operasi mengunci baris
// pesanan (FOR UPDATE) lalu memastikan statusnya masih open.

// ensureTableFree: meja ada, aktif, dan tidak dipakai pesanan terbuka lain (keputusan 103).
func ensureTableFree(r repo.TransactionRepoInterface, tableID, excludeBillID int) error {
	name, active, found, occupiedBy, err := r.LockTableForOrder(tableID, excludeBillID)
	if err != nil {
		return err
	}
	if !found || !active {
		return &errors.BadRequestError{Message: "Meja tidak ditemukan atau tidak aktif"}
	}
	if occupiedBy > 0 {
		return &errors.BadRequestError{Message: fmt.Sprintf("Meja %s masih punya pesanan terbuka. Tambah item ke pesanan itu atau gabung bill", name)}
	}
	return nil
}

// openTx menjalankan fn di dalam transaksi DB dengan pesanan terkunci dan berstatus open.
func (s *transactionService) openTx(id int, fn func(r repo.TransactionRepoInterface, bill *repo.LockedBill, tx *gorm.DB) error) error {
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		bill, err := r.LockBill(id)
		if err != nil {
			return err
		}
		if bill == nil {
			return &errors.NotFoundError{Message: "Pesanan tidak ditemukan"}
		}
		if bill.Status != repo.StatusOpen {
			return repo.ErrNotOpen
		}
		return fn(r, bill, tx)
	})
	return mapOpenErr(txErr)
}

// rejectOverdue: pesanan tertunda (dari hari sebelumnya) hanya bisa dibayar, dijadikan utang, atau dibatalkan.
func rejectOverdue(bill *repo.LockedBill) error {
	if bill.IsOverdue() {
		return &errors.BadRequestError{Message: fmt.Sprintf("Pesanan tertunda dari tanggal %s hanya bisa dibayar, dijadikan utang, atau dibatalkan", bill.BillDate)}
	}
	return nil
}

// ensureDrawerNotPending: kasir yang kasnya tertunda (dari hari sebelumnya) tidak bisa berjualan baru sampai
// kas itu ditutup (keputusan 112).
func (s *transactionService) ensureDrawerNotPending(userID int) error {
	drawer, err := s.cashDrawerRepo.GetOpenCashDrawer(userID)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
	}
	if drawer != nil && drawer.OpenTime.In(time_helper.GetTimeNow().Location()).Format("2006-01-02") < time_helper.GetTimeNow().Format("2006-01-02") {
		return &errors.BadRequestError{Message: fmt.Sprintf("Kas Anda masih tertunda dari tanggal %s. Selesaikan pesanan tertunda lalu tutup kas itu dulu", drawer.OpenTime.Format("02-01-2006"))}
	}
	return nil
}

func mapOpenErr(err error) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, repo.ErrNotOpen) {
		return &errors.BadRequestError{Message: "Pesanan sudah tidak terbuka (sudah dibayar, dibatalkan, atau digabung)"}
	}
	var badReq *errors.BadRequestError
	if stderrors.As(err, &badReq) {
		return badReq
	}
	var notFound *errors.NotFoundError
	if stderrors.As(err, &notFound) {
		return notFound
	}
	return &errors.InternalServerError{Message: err.Error()}
}

func (s *transactionService) ListOpen() ([]dto.OpenBillSummary, error) {
	rows, err := s.repo.ListOpen()
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	return rows, nil
}

// AddOpenItems: harga dan aturan menu sama dengan pesanan baru (resolveItems).
func (s *transactionService) AddOpenItems(id int, req *dto.OpenAddItemsRequest, userID int) (*dto.TransactionResponse, error) {
	tmp := &dto.CreateTransactionRequest{Items: req.Items}
	if err := s.resolveItems(tmp); err != nil {
		return nil, err
	}
	if err := s.ensureDrawerNotPending(userID); err != nil {
		return nil, err
	}
	if err := s.openTx(id, func(r repo.TransactionRepoInterface, bill *repo.LockedBill, _ *gorm.DB) error {
		if err := rejectOverdue(bill); err != nil {
			return err
		}
		return r.AddOpenItems(bill, tmp.Items, userID)
	}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

// RemoveOpenItem: hapus/kurangi item. Barang jadi wajib memilih kembali ke stok atau dibuang; dibuang
// wajib alasan (keputusan 100, 105).
func (s *transactionService) RemoveOpenItem(id int, req *dto.OpenRemoveItemRequest, userID int) (*dto.TransactionResponse, error) {
	reason := strings.TrimSpace(req.Reason)
	if req.Quantity != math.Trunc(req.Quantity) {
		return nil, &errors.BadRequestError{Message: "Jumlah harus bilangan bulat"}
	}
	if err := s.openTx(id, func(r repo.TransactionRepoInterface, bill *repo.LockedBill, _ *gorm.DB) error {
		if err := rejectOverdue(bill); err != nil {
			return err
		}
		item, err := r.GetItem(id, req.ItemID)
		if err != nil {
			return err
		}
		if item == nil {
			return &errors.BadRequestError{Message: "Item tidak ada di pesanan ini"}
		}
		action := req.Action
		if item.ProductID > 0 {
			if action == "" {
				return &errors.BadRequestError{Message: fmt.Sprintf("Pilih %s kembali ke stok atau dibuang", item.ProductName)}
			}
			if action == repo.StockActionWaste && reason == "" {
				return &errors.BadRequestError{Message: "Alasan wajib diisi untuk barang yang dibuang"}
			}
		} else {
			action = ""
		}
		qty := req.Quantity
		if qty == 0 {
			qty = item.Quantity
		}
		return r.RemoveOpenItem(bill, item, qty, action, reason, userID)
	}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

// MoveOpenTable: hanya pesanan makan di tempat, ke meja aktif yang kosong (keputusan 103).
func (s *transactionService) MoveOpenTable(id int, req *dto.OpenMoveTableRequest) (*dto.TransactionResponse, error) {
	if err := s.openTx(id, func(r repo.TransactionRepoInterface, bill *repo.LockedBill, _ *gorm.DB) error {
		if err := rejectOverdue(bill); err != nil {
			return err
		}
		if bill.OrderType != orderTypeDineIn {
			return &errors.BadRequestError{Message: "Pindah meja hanya untuk pesanan makan di tempat"}
		}
		if bill.TableID != nil && *bill.TableID == req.TableID {
			return &errors.BadRequestError{Message: "Pesanan sudah di meja ini"}
		}
		if err := ensureTableFree(r, req.TableID, id); err != nil {
			return err
		}
		name, _, _, err := r.GetTableForOrder(req.TableID)
		if err != nil {
			return err
		}
		return r.MoveOpenTable(id, req.TableID, name)
	}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

// MergeOpen: pesanan terbuka sumber digabung ke pesanan pada URI (keputusan 107). Kedua baris dikunci
// berurutan menurut id agar tidak terjadi deadlock saat dua kasir menggabung bersamaan.
func (s *transactionService) MergeOpen(targetID int, req *dto.OpenMergeRequest) (*dto.TransactionResponse, error) {
	if req.SourceID == targetID {
		return nil, &errors.BadRequestError{Message: "Pesanan tidak bisa digabung dengan dirinya sendiri"}
	}
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		first, second := targetID, req.SourceID
		if second < first {
			first, second = second, first
		}
		bills := map[int]*repo.LockedBill{}
		for _, id := range []int{first, second} {
			b, err := r.LockBill(id)
			if err != nil {
				return err
			}
			if b == nil {
				return &errors.NotFoundError{Message: "Pesanan tidak ditemukan"}
			}
			if b.Status != repo.StatusOpen {
				return repo.ErrNotOpen
			}
			if err := rejectOverdue(b); err != nil {
				return err
			}
			bills[id] = b
		}
		return r.MergeOpen(targetID, req.SourceID)
	})
	if err := mapOpenErr(txErr); err != nil {
		return nil, err
	}
	return s.repo.GetByID(targetID)
}

// CancelOpen: batal pesanan terbuka, alasan wajib; barang jadi per item kembali ke stok (bawaan) atau dibuang.
func (s *transactionService) CancelOpen(id int, req *dto.OpenCancelRequest, userID int) error {
	reason := strings.TrimSpace(req.Reason)
	return s.openTx(id, func(r repo.TransactionRepoInterface, bill *repo.LockedBill, _ *gorm.DB) error {
		items, err := r.GetItems(id)
		if err != nil {
			return err
		}
		waste := map[int]bool{}
		seen := map[int]bool{}
		for _, a := range req.Items {
			if seen[a.ItemID] {
				return &errors.BadRequestError{Message: "Pilihan item ganda"}
			}
			seen[a.ItemID] = true
			found := false
			for _, it := range items {
				if it.ID != a.ItemID {
					continue
				}
				found = true
				if a.Action == repo.StockActionWaste {
					if it.ProductID == 0 {
						return &errors.BadRequestError{Message: fmt.Sprintf("%s bukan barang jadi, tidak bisa dipilih dibuang", it.ProductName)}
					}
					waste[it.ID] = true
				}
			}
			if !found {
				return &errors.BadRequestError{Message: "Item tidak ada di pesanan ini"}
			}
		}
		return r.CancelOpen(bill, waste, reason, userID)
	})
}

// PayOpen: bayar atau jadikan utang (is_credit + pelanggan). Diskon nota dan pajak diisi sekarang
// (keputusan 104); uang tunai masuk ke kas kasir penerima (keputusan 106).
func (s *transactionService) PayOpen(id int, req *dto.OpenPayRequest, userID int) (*dto.TransactionResponse, error) {
	openDrawer, err := s.cashDrawerRepo.GetOpenCashDrawer(userID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if req.ShiftID != nil && (openDrawer == nil || openDrawer.ShiftID == nil || *openDrawer.ShiftID != *req.ShiftID) {
		return nil, &errors.BadRequestError{Message: "Shift tidak sesuai dengan sesi kasir yang sedang aktif"}
	}
	if req.IsCredit && req.CustomerID == nil {
		return nil, &errors.BadRequestError{Message: "Pilih pelanggan untuk menjadikan utang"}
	}
	txErr := s.openTx(id, func(r repo.TransactionRepoInterface, bill *repo.LockedBill, tx *gorm.DB) error {
		items, err := r.GetItems(id)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return &errors.BadRequestError{Message: "Pesanan belum berisi item"}
		}
		req.Subtotal = 0
		for _, it := range items {
			req.Subtotal += it.Subtotal
		}
		if req.Discount > req.Subtotal {
			return &errors.BadRequestError{Message: "Diskon transaksi tidak boleh lebih besar dari subtotal"}
		}
		req.TotalAmount = math.Max(0, req.Subtotal-req.Discount+req.Tax)
		if !req.IsCredit && req.TotalAmount > 0 && req.PaymentAmount < req.TotalAmount {
			return &errors.BadRequestError{Message: "Jumlah pembayaran kurang dari total transaksi"}
		}
		req.ChangeAmount = 0
		if req.TotalAmount > 0 && req.PaymentAmount > req.TotalAmount {
			req.ChangeAmount = req.PaymentAmount - req.TotalAmount
		}
		now := time_helper.GetTimeNow()
		cd := s.cashDrawerRepo.WithTx(tx)
		ownerID, txDate := userID, now
		var drawerID int
		if bill.IsOverdue() {
			// Keputusan 113: tanggal penjualan = tanggal pesanan; uang ke kas tertunda hari itu (kas pembuat,
			// lalu kas penerima, lalu kas lain hari itu). Pelaku penyelesaian tercatat di settled_by.
			creator := 0
			if bill.CreatedBy != nil {
				creator = *bill.CreatedBy
			}
			target, err := cd.FindDrawerForDate(bill.BillDate, creator, userID)
			if err != nil {
				return err
			}
			txDate = bill.TransactionDate
			if target != nil {
				ownerID, drawerID = target.UserID, target.ID
				// Tetap di dalam rentang kas penampung agar saldo kas yang diharapkan ikut menghitungnya.
				if target.OpenTime.After(txDate) {
					txDate = target.OpenTime
				}
			} else {
				// Tidak ada kas hari itu (pesanan dibuat tanpa kas): catat seperti pembayaran biasa ke kas
				// kasir penerima agar uang tetap masuk ke sebuah kas.
				drawer, err := cd.GetOpenCashDrawer(userID)
				if err != nil {
					return err
				}
				txDate = now
				if drawer != nil {
					drawerID = drawer.ID
				}
			}
		} else {
			if err := s.ensureDrawerNotPending(userID); err != nil {
				return err
			}
			drawer, err := cd.GetOpenCashDrawer(userID)
			if err != nil {
				return err
			}
			if drawer != nil {
				drawerID = drawer.ID
			}
		}
		if err := r.PayOpen(id, req, ownerID, userID, txDate, now); err != nil {
			return err
		}
		if req.PaymentMethod == "cash" && req.TotalAmount > 0 && drawerID > 0 {
			if err := cd.UpdateSales(drawerID, req.TotalAmount, req.TotalAmount, now); err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return s.repo.GetByID(id)
}
