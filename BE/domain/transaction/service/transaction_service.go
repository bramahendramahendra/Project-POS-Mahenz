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

	resp, err := s.repo.Create(req, userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "stok_insufficient:") {
			name := strings.TrimPrefix(err.Error(), "stok_insufficient:")
			return nil, &errors.BadRequestError{Message: "Stok tidak mencukupi untuk " + name}
		}
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	if req.PaymentMethod == "cash" && openCashDrawer != nil {
		_ = s.cashDrawerRepo.UpdateSales(openCashDrawer.ID, resp.TotalAmount, resp.TotalAmount)
	}

	return resp, nil
}

// recalculateTotals menghitung ulang subtotal/diskon/pajak/total transaksi di server
// berdasarkan harga produk asli, mengabaikan nilai price/subtotal/total_amount dari client.
func (s *transactionService) recalculateTotals(req *dto.CreateTransactionRequest) error {
	var computedSubtotal float64

	for i, item := range req.Items {
		unitPrice, err := s.resolveUnitPrice(item.ProductID, item.UnitID, item.Quantity)
		if err != nil {
			return err
		}

		lineGross := unitPrice * item.Quantity
		if item.DiscountItem < 0 || item.DiscountItem > lineGross {
			return &errors.BadRequestError{Message: "Diskon item tidak valid untuk " + item.ProductName}
		}

		req.Items[i].Price = unitPrice
		req.Items[i].Subtotal = lineGross - item.DiscountItem
		computedSubtotal += req.Items[i].Subtotal
	}

	if req.Discount < 0 || req.Discount > computedSubtotal {
		return &errors.BadRequestError{Message: "Diskon transaksi tidak valid"}
	}

	computedTotal := computedSubtotal - req.Discount + req.Tax
	if computedTotal < 0 {
		computedTotal = 0
	}

	req.Subtotal = computedSubtotal
	req.TotalAmount = computedTotal
	return nil
}

// resolveUnitPrice mengambil harga jual asli produk dari master data, bukan dari payload client.
func (s *transactionService) resolveUnitPrice(productID int, unitID *int, quantity float64) (float64, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return 0, &errors.InternalServerError{Message: err.Error()}
	}
	if product == nil {
		return 0, &errors.NotFoundError{Message: "Produk tidak ditemukan"}
	}

	if unitID != nil && *unitID > 0 {
		packages, err := s.productRepo.GetPackagesByProduct(productID)
		if err != nil {
			return 0, &errors.InternalServerError{Message: err.Error()}
		}
		for _, pkg := range packages {
			if pkg.ID == *unitID {
				return pkg.SellingPrice, nil
			}
		}
		return 0, &errors.BadRequestError{Message: "Kemasan produk tidak ditemukan untuk " + product.Name}
	}

	prices, err := s.productRepo.GetPricesByProduct(productID)
	if err != nil {
		return 0, &errors.InternalServerError{Message: err.Error()}
	}

	price := product.SellingPrice
	bestMinQty := -1.0
	for _, tier := range prices {
		if quantity >= tier.MinQty && tier.MinQty > bestMinQty {
			price = tier.Price
			bestMinQty = tier.MinQty
		}
	}
	return price, nil
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
