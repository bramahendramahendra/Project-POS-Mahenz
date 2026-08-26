package service

import (
	stderrors "errors"
	"fmt"
	"time"

	"pos_api/domain/backdate/dto"
	"pos_api/domain/backdate/repo"
	cb_repo "pos_api/domain/customer_balance/repo"
	"pos_api/errors"
	time_helper "pos_api/helper/time"

	"gorm.io/gorm"
)

type BackdateServiceInterface interface {
	OpenCashDrawer(adminUserID int, req *dto.OpenCashDrawerRequest) (*dto.OpenCashDrawerResponse, error)
	GetCurrentCashDrawer(adminUserID int) (*dto.CurrentCashDrawerResponse, error)
	CloseCashDrawer(adminUserID int, req *dto.CloseCashDrawerRequest) (*dto.CloseCashDrawerResponse, error)
	CreateTransaction(adminUserID int, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}

type backdateService struct {
	repo                repo.BackdateRepoInterface
	customerBalanceRepo cb_repo.CustomerBalanceRepoInterface
}

func NewBackdateService(
	r repo.BackdateRepoInterface,
	customerBalanceRepo cb_repo.CustomerBalanceRepoInterface,
) *backdateService {
	return &backdateService{
		repo:                r,
		customerBalanceRepo: customerBalanceRepo,
	}
}

func (s *backdateService) OpenCashDrawer(adminUserID int, req *dto.OpenCashDrawerRequest) (*dto.OpenCashDrawerResponse, error) {
	if req.OpeningBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo awal tidak boleh negatif"}
	}

	// Parse date
	targetDate, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		return nil, &errors.BadRequestError{Message: "Format tanggal tidak valid (gunakan YYYY-MM-DD)"}
	}

	// Validasi: tanggal harus sebelum hari ini
	now := time_helper.GetTimeNow()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if !targetDate.Before(today) {
		return nil, &errors.BadRequestError{Message: "Tanggal harus sebelum hari ini"}
	}

	// Cek apakah admin sudah punya kas backdate yang open
	existing, err := s.repo.GetCurrentCashDrawer(adminUserID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Sudah ada kas historis yang terbuka. Tutup terlebih dahulu sebelum membuka kas baru."}
	}

	// Set open_time = tanggal target pagi (08:00:00)
	openTime := targetDate.Add(8 * time.Hour)

	id, err := s.repo.OpenCashDrawer(req.UserID, req.ShiftID, openTime, req.OpeningBalance, req.Notes, adminUserID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	return &dto.OpenCashDrawerResponse{ID: int(id)}, nil
}

func (s *backdateService) GetCurrentCashDrawer(adminUserID int) (*dto.CurrentCashDrawerResponse, error) {
	cd, err := s.repo.GetCurrentCashDrawer(adminUserID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if cd == nil {
		return nil, nil
	}

	// Build response with extra info
	var shiftNamePtr *string

	// Re-query with full detail
	type detailRow struct {
		UserName      string  `gorm:"column:user_name"`
		ShiftName     *string `gorm:"column:shift_name"`
		CreatedByName string  `gorm:"column:created_by_name"`
	}
	var detail detailRow
	s.repo.GetDB().Raw(
		`SELECT COALESCE(u.full_name, '') as user_name, s.name as shift_name, COALESCE(cb.full_name, '') as created_by_name
		 FROM cash_drawer cd
		 LEFT JOIN users u ON cd.user_id = u.id
		 LEFT JOIN shifts s ON cd.shift_id = s.id
		 LEFT JOIN users cb ON cd.created_by = cb.id
		 WHERE cd.id = ?`, cd.ID,
	).Scan(&detail)

	shiftNamePtr = detail.ShiftName

	return &dto.CurrentCashDrawerResponse{
		ID:              cd.ID,
		UserID:          cd.UserID,
		UserName:        detail.UserName,
		ShiftID:         cd.ShiftID,
		ShiftName:       shiftNamePtr,
		OpenTime:        cd.OpenTime,
		OpeningBalance:  cd.OpeningBalance,
		TotalSales:      cd.TotalSales,
		TotalCashSales:  cd.TotalCashSales,
		TotalExpenses:   cd.TotalExpenses,
		ExpectedBalance: cd.ExpectedBalance,
		Status:          cd.Status,
		OpenNotes:       cd.Notes,
		CreatedBy:       cd.CreatedBy,
		CreatedByName:   detail.CreatedByName,
	}, nil
}

func (s *backdateService) CloseCashDrawer(adminUserID int, req *dto.CloseCashDrawerRequest) (*dto.CloseCashDrawerResponse, error) {
	if req.ClosingBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo akhir tidak boleh negatif"}
	}

	cd, err := s.repo.GetCashDrawerByID(req.ID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if cd == nil {
		return nil, &errors.NotFoundError{Message: "Kas historis tidak ditemukan"}
	}
	if cd.CreatedBy == nil || *cd.CreatedBy != adminUserID {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas historis ini"}
	}
	if cd.Status != "open" {
		return nil, &errors.BadRequestError{Message: "Kas sudah ditutup"}
	}

	expected := cd.ExpectedBalance
	difference := req.ClosingBalance - expected

	// Close time = tanggal kas + 23:59
	closeTime := cd.OpenTime.Truncate(24 * time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	if err := s.repo.CloseCashDrawer(cd.ID, req.ClosingBalance, expected, difference, req.Notes, closeTime); err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}

	return &dto.CloseCashDrawerResponse{
		ExpectedBalance: expected,
		ClosingBalance:  req.ClosingBalance,
		Difference:      difference,
	}, nil
}

func (s *backdateService) CreateTransaction(adminUserID int, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	// Get current backdate cash drawer
	cd, err := s.repo.GetCurrentCashDrawer(adminUserID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if cd == nil {
		return nil, &errors.BadRequestError{Message: "Belum ada kas historis yang terbuka. Buka kas terlebih dahulu."}
	}

	// Parse transaction time (HH:mm)
	timeParts, err := time.Parse("15:04", req.TransactionTime)
	if err != nil {
		return nil, &errors.BadRequestError{Message: "Format jam tidak valid (gunakan HH:mm)"}
	}

	// Combine date from cash drawer + time from request
	cashDate := cd.OpenTime.Truncate(24 * time.Hour)
	transactionDate := cashDate.Add(time.Duration(timeParts.Hour())*time.Hour + time.Duration(timeParts.Minute())*time.Minute)

	// For backdate, we allow custom prices (edit harga live)
	// So we DON'T recalculate prices — use what FE sends
	// But still validate subtotal/total consistency
	var calculatedSubtotal float64
	for _, item := range req.Items {
		calculatedSubtotal += item.Subtotal
	}

	totalAmount := calculatedSubtotal - req.Discount + req.Tax
	if totalAmount < 0 {
		totalAmount = 0
	}
	req.TotalAmount = totalAmount
	req.Subtotal = calculatedSubtotal

	// Validate balance_used
	if req.BalanceUsed > 0 {
		if req.CustomerID == nil {
			return nil, &errors.BadRequestError{Message: "Pilih pelanggan untuk menggunakan saldo"}
		}
		if req.BalanceUsed > req.TotalAmount {
			return nil, &errors.BadRequestError{Message: "Saldo yang digunakan melebihi total belanja"}
		}
	}

	// Calculate effective total
	effectiveTotal := req.TotalAmount - req.BalanceUsed
	if !req.IsCredit && effectiveTotal > 0 && req.PaymentAmount < effectiveTotal {
		return nil, &errors.BadRequestError{Message: "Jumlah pembayaran kurang dari total transaksi"}
	}
	if effectiveTotal > 0 && req.PaymentAmount > effectiveTotal {
		req.ChangeAmount = req.PaymentAmount - effectiveTotal
	} else {
		req.ChangeAmount = 0
	}

	// Build params
	itemParams := make([]repo.CreateTransactionItemParams, len(req.Items))
	for i, item := range req.Items {
		itemParams[i] = repo.CreateTransactionItemParams{
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			Quantity:     item.Quantity,
			Unit:         item.Unit,
			Price:        item.Price,
			Subtotal:     item.Subtotal,
			DiscountItem: item.DiscountItem,
			UnitID:       item.UnitID,
		}
	}

	params := &repo.CreateTransactionParams{
		Subtotal:      req.Subtotal,
		Discount:      req.Discount,
		Tax:           req.Tax,
		TotalAmount:   req.TotalAmount,
		PaymentMethod: req.PaymentMethod,
		PaymentAmount: req.PaymentAmount,
		ChangeAmount:  req.ChangeAmount,
		BalanceUsed:   req.BalanceUsed,
		CustomerID:    req.CustomerID,
		IsCredit:      req.IsCredit,
		DeviceSource:  req.DeviceSource,
		Items:         itemParams,
	}

	// Wrap in DB transaction for kas update + balance deduction
	var result *repo.CreateTransactionResult
	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// Create transaction (includes stock reduction)
		var createErr error
		result, createErr = s.repo.CreateTransaction(cd.UserID, cd.ShiftID, transactionDate, cd.ID, adminUserID, params)
		if createErr != nil {
			return createErr
		}

		// Deduct customer balance if used
		if req.BalanceUsed > 0 && req.CustomerID != nil {
			cbRepo := s.customerBalanceRepo.WithTx(tx)
			refID := result.ID
			_, err := cbRepo.Deduct(*req.CustomerID, req.BalanceUsed, "transaction", &refID, fmt.Sprintf("Bayar belanja backdate dari saldo (%s)", result.TransactionCode), adminUserID)
			if err != nil {
				return &errors.BadRequestError{Message: "Saldo pelanggan tidak mencukupi"}
			}
		}

		// Update kas harian (hanya cash, hanya effective_total)
		if req.PaymentMethod == "cash" {
			netCash := req.TotalAmount - req.BalanceUsed
			if netCash > 0 {
				if err := tx.Exec(
					`UPDATE cash_drawer SET total_sales = total_sales + ?, total_cash_sales = total_cash_sales + ?, updated_at = ? WHERE id = ?`,
					netCash, netCash, time_helper.GetTimeNow(), cd.ID,
				).Error; err != nil {
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

	return &dto.CreateTransactionResponse{
		ID:              result.ID,
		TransactionCode: result.TransactionCode,
		TransactionDate: result.TransactionDate,
		TotalAmount:     result.TotalAmount,
	}, nil
}
