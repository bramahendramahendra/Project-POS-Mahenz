package service

import (
	stderrors "errors"
	"fmt"
	"math"
	"strings"

	cash_drawer_model "coffee_api/domain/cash_drawer/model"
	"coffee_api/domain/transaction/dto"
	"coffee_api/domain/transaction/model"
	repo "coffee_api/domain/transaction/repo"
	time_helper "coffee_api/helper/time"

	"coffee_api/errors"

	"gorm.io/gorm"
)

const (
	orderTypeDineIn   = "dine_in"
	orderTypeTakeAway = "take_away"
	stockModeDirect   = "direct"
)

var temperatureLabel = map[string]string{"hot": "Panas", "iced": "Dingin"}

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

// Create: pesanan bayar sekarang (tahap 4.7a). Harga, subtotal, total, dan kembalian dihitung server.
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

	if err := s.resolveOrder(req); err != nil {
		return nil, err
	}
	if err := s.resolveItems(req); err != nil {
		return nil, err
	}

	if req.Discount > req.Subtotal {
		return nil, &errors.BadRequestError{Message: "Diskon transaksi tidak boleh lebih besar dari subtotal"}
	}
	req.TotalAmount = math.Max(0, req.Subtotal-req.Discount+req.Tax)

	if req.IsCredit && req.CustomerID == nil {
		return nil, &errors.BadRequestError{Message: "Pilih pelanggan untuk transaksi utang"}
	}
	if !req.IsCredit && req.TotalAmount > 0 && req.PaymentAmount < req.TotalAmount {
		return nil, &errors.BadRequestError{Message: "Jumlah pembayaran kurang dari total transaksi"}
	}
	if req.TotalAmount > 0 && req.PaymentAmount > req.TotalAmount {
		req.ChangeAmount = req.PaymentAmount - req.TotalAmount
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

		// Kas harian bertambah hanya untuk pembayaran tunai (seperti Toko).
		if req.PaymentMethod == "cash" && req.TotalAmount > 0 {
			drawer, err := cashDrawerRepo.GetOpenCashDrawer(userID)
			if err != nil {
				return err
			}
			if drawer != nil {
				if err := cashDrawerRepo.UpdateSales(drawer.ID, req.TotalAmount, req.TotalAmount, time_helper.GetTimeNow()); err != nil {
					return err
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

	// Lengkapi item jawaban dari data tersimpan (id item, tambahan) agar sama dengan detail.
	if full, err := s.repo.GetByID(resp.ID); err == nil && full != nil {
		resp.TransactionResponse = *full
	}
	return resp, nil
}

// resolveOrder: jenis pesanan, meja, dan nama pemesan (keputusan 25, 97, 99).
func (s *transactionService) resolveOrder(req *dto.CreateTransactionRequest) error {
	req.OrderName = strings.TrimSpace(req.OrderName)
	switch req.OrderType {
	case orderTypeDineIn:
		if req.TableID == nil {
			return &errors.BadRequestError{Message: "Pilih meja untuk pesanan makan di tempat"}
		}
		name, active, found, err := s.repo.GetTableForOrder(*req.TableID)
		if err != nil {
			return &errors.InternalServerError{Message: err.Error()}
		}
		if !found || !active {
			return &errors.BadRequestError{Message: "Meja tidak ditemukan atau tidak aktif"}
		}
		req.TableName = name
	case orderTypeTakeAway:
		if req.TableID != nil {
			return &errors.BadRequestError{Message: "Pesanan bawa pulang tidak memakai meja"}
		}
	}
	return nil
}

// resolveItems memeriksa setiap baris dan menghitung harganya: harga kombinasi varian + harga
// tambahan terpilih, dikali jumlah, dikurangi diskon item (keputusan 53-55, 57, 87, 88, 91).
func (s *transactionService) resolveItems(req *dto.CreateTransactionRequest) error {
	bad := func(i int, format string, args ...any) error {
		return &errors.BadRequestError{Message: fmt.Sprintf("Item ke-%d: ", i+1) + fmt.Sprintf(format, args...)}
	}
	optionsByMenu := map[int][]model.SaleModifierOption{}
	req.Subtotal = 0

	for i := range req.Items {
		item := &req.Items[i]
		if item.Quantity != math.Trunc(item.Quantity) {
			return bad(i, "jumlah harus bilangan bulat")
		}
		item.Notes = strings.TrimSpace(item.Notes)

		v, err := s.repo.GetSaleVariant(item.VariantID)
		if err != nil {
			return &errors.InternalServerError{Message: err.Error()}
		}
		if v == nil {
			return bad(i, "menu tidak ditemukan")
		}
		switch {
		case !v.MenuActive:
			return bad(i, "menu %s tidak aktif", v.MenuName)
		case v.SoldOut:
			return bad(i, "menu %s sedang habis", v.MenuName)
		case !v.VariantActive:
			return bad(i, "pilihan %s %s tidak dijual", v.MenuName, variantName(v))
		}
		item.MenuItemID = v.MenuItemID
		item.MenuName = v.MenuName
		item.VariantName = variantName(v)
		item.ProductID = nil
		if v.StockMode == stockModeDirect {
			if v.ProductID == nil {
				return bad(i, "menu %s belum ditautkan ke produk barang jadi", v.MenuName)
			}
			item.ProductID = v.ProductID
		}

		opts, ok := optionsByMenu[v.MenuItemID]
		if !ok {
			opts, err = s.repo.GetMenuModifierOptions(v.MenuItemID)
			if err != nil {
				return &errors.InternalServerError{Message: err.Error()}
			}
			optionsByMenu[v.MenuItemID] = opts
		}
		mods, err := resolveModifiers(v.MenuName, item.ModifierOptionIDs, opts)
		if err != nil {
			return bad(i, "%s", err.Error())
		}
		item.Modifiers = mods

		unit := v.Price
		for _, m := range mods {
			unit += m.Price
		}
		gross := unit * item.Quantity
		if item.DiscountItem > gross {
			return bad(i, "diskon item tidak boleh lebih besar dari harga %s", v.MenuName)
		}
		item.Price = unit
		item.Subtotal = gross - item.DiscountItem
		req.Subtotal += item.Subtotal
	}
	return nil
}

// resolveModifiers: pilihan harus milik kelompok yang terpasang di menu, kelompok dan pilihan aktif,
// tidak ganda, kelompok "pilih satu" paling banyak satu, kelompok wajib minimal satu (keputusan 87).
func resolveModifiers(menuName string, selected []int, opts []model.SaleModifierOption) ([]dto.LineModifier, error) {
	byID := make(map[int]model.SaleModifierOption, len(opts))
	for _, o := range opts {
		byID[o.OptionID] = o
	}
	seen := map[int]bool{}
	perGroup := map[int]int{}
	var result []dto.LineModifier
	for _, id := range selected {
		o, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("tambahan tidak tersedia untuk menu %s", menuName)
		}
		if !o.OptionActive || !o.GroupActive {
			return nil, fmt.Errorf("tambahan %s tidak tersedia", o.OptionName)
		}
		if seen[id] {
			return nil, fmt.Errorf("tambahan %s dipilih lebih dari sekali", o.OptionName)
		}
		seen[id] = true
		perGroup[o.GroupID]++
		if o.SelectionType == "single" && perGroup[o.GroupID] > 1 {
			return nil, fmt.Errorf("kelompok %s hanya boleh satu pilihan", o.GroupName)
		}
		result = append(result, dto.LineModifier{OptionID: o.OptionID, GroupName: o.GroupName, OptionName: o.OptionName, Price: o.Price})
	}
	checked := map[int]bool{}
	for _, o := range opts {
		if checked[o.GroupID] || !o.IsRequired || !o.GroupActive {
			continue
		}
		checked[o.GroupID] = true
		if perGroup[o.GroupID] == 0 && groupHasActiveOption(opts, o.GroupID) {
			return nil, fmt.Errorf("kelompok %s wajib dipilih", o.GroupName)
		}
	}
	return result, nil
}

func groupHasActiveOption(opts []model.SaleModifierOption, groupID int) bool {
	for _, o := range opts {
		if o.GroupID == groupID && o.OptionActive {
			return true
		}
	}
	return false
}

// variantName: "Regular · Dingin", "Large", "Panas", atau kosong untuk menu tanpa varian.
func variantName(v *model.SaleVariant) string {
	var parts []string
	if v.SizeName != "" {
		parts = append(parts, v.SizeName)
	}
	if v.Temperature != nil {
		if label, ok := temperatureLabel[*v.Temperature]; ok {
			parts = append(parts, label)
		}
	}
	return strings.Join(parts, " · ")
}

// Void transaksi lunas (admin/owner). Barang jadi kembali ke stok kecuali item yang dipilih
// "dibuang" (keputusan 100: per item, alasan wajib). Kas tunai dikembalikan seperti Toko.
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

	waste := map[int]bool{}
	seen := map[int]bool{}
	for _, a := range req.Body.Items {
		if seen[a.ItemID] {
			return &errors.BadRequestError{Message: "Pilihan item ganda"}
		}
		seen[a.ItemID] = true
		var item *dto.TransactionItemResponse
		for i := range t.Items {
			if t.Items[i].ID == a.ItemID {
				item = &t.Items[i]
			}
		}
		if item == nil {
			return &errors.BadRequestError{Message: "Item tidak ada di transaksi ini"}
		}
		if a.Action == "waste" {
			if item.ProductID == nil {
				return &errors.BadRequestError{Message: fmt.Sprintf("%s bukan barang jadi, tidak bisa dipilih dibuang", item.ProductName)}
			}
			waste[a.ItemID] = true
		}
	}
	reason := strings.TrimSpace(req.Body.Reason)
	if len(waste) > 0 && reason == "" {
		return &errors.BadRequestError{Message: "Alasan wajib diisi untuk barang yang dibuang"}
	}

	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		txnRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		if err := txnRepo.Void(req.ID, userID, waste, reason); err != nil {
			return err
		}

		// Kas harian dikurangi kembali untuk transaksi tunai (seperti Toko).
		if t.PaymentMethod == "cash" && t.TotalAmount > 0 {
			var drawer *cash_drawer_model.CashDrawer
			if t.CashDrawerID != nil {
				d, dErr := cashDrawerRepo.GetByID(*t.CashDrawerID)
				if dErr != nil {
					return dErr
				}
				if d != nil && d.Status == "open" {
					drawer = d
				}
			} else {
				d, dErr := cashDrawerRepo.GetOpenCashDrawer(t.UserID)
				if dErr != nil {
					return dErr
				}
				drawer = d
			}
			if drawer != nil {
				if err := cashDrawerRepo.UpdateSales(drawer.ID, -t.TotalAmount, -t.TotalAmount, time_helper.GetTimeNow()); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if txErr != nil {
		if stderrors.Is(txErr, repo.ErrAlreadyVoid) {
			return &errors.BadRequestError{Message: "Transaksi sudah di-void"}
		}
		var badReq *errors.BadRequestError
		if stderrors.As(txErr, &badReq) {
			return badReq
		}
		return &errors.InternalServerError{Message: txErr.Error()}
	}
	return nil
}
