package service

import (
	dto "pos_api/domain/cash_drawer/dto"
	"pos_api/errors"
)

func (s *cashDrawerService) GetCurrent(userID int) (*dto.CurrentCashDrawerResponse, error) {
	return s.repo.GetCurrent(userID)
}

func (s *cashDrawerService) GetMyCash(userID int) (*dto.MyCashResponse, error) {
	dataDB, transactions, expenses, err := s.repo.GetMyCash(userID)
	if err != nil {
		return nil, err
	}

	if dataDB == nil {
		return &dto.MyCashResponse{
			Status:       "closed",
			Transactions: []dto.CashDrawerTransaction{},
			Expenses:     []dto.CashDrawerExpenseItem{},
		}, nil
	}

	trxList := make([]dto.CashDrawerTransaction, len(transactions))
	for i, t := range transactions {
		trxList[i] = dto.CashDrawerTransaction{
			TransactionDate: t.TransactionDate,
			TransactionCode: t.TransactionCode,
			CustomerName:    t.CustomerName,
			TotalAmount:     t.TotalAmount,
		}
	}

	expList := make([]dto.CashDrawerExpenseItem, len(expenses))
	for i, e := range expenses {
		expList[i] = dto.CashDrawerExpenseItem{
			Category:    e.Category,
			Description: e.Description,
			Amount:      e.Amount,
		}
	}

	return &dto.MyCashResponse{
		ID:              &dataDB.ID,
		Status:          dataDB.Status,
		ShiftName:       dataDB.ShiftName,
		ShiftStart:      dataDB.ShiftStart,
		ShiftEnd:        dataDB.ShiftEnd,
		OpenTime:        &dataDB.OpenTime,
		OpeningBalance:  dataDB.OpeningBalance,
		TotalCashSales:  dataDB.TotalCashSales,
		TotalExpenses:   dataDB.TotalExpenses,
		ExpectedBalance: dataDB.ExpectedBalance,
		OpenNotes:       dataDB.OpenNotes,
		Transactions:    trxList,
		Expenses:        expList,
	}, nil
}

func (s *cashDrawerService) GetByID(id int, requestingUserID int, role string) (*dto.CashDrawerDetailResponse, error) {
	dataDB, transactions, expenses, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, err
	}
	if dataDB == nil {
		return nil, &errors.NotFoundError{Message: "Kas tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && dataDB.UserID != requestingUserID {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas ini"}
	}

	trxList := make([]dto.CashDrawerTransaction, len(transactions))
	for i, t := range transactions {
		trxList[i] = dto.CashDrawerTransaction{
			TransactionDate: t.TransactionDate,
			TransactionCode: t.TransactionCode,
			CustomerName:    t.CustomerName,
			TotalAmount:     t.TotalAmount,
		}
	}

	expList := make([]dto.CashDrawerExpenseItem, len(expenses))
	for i, e := range expenses {
		expList[i] = dto.CashDrawerExpenseItem{
			Category:    e.Category,
			Description: e.Description,
			Amount:      e.Amount,
		}
	}

	data := &dto.CashDrawerDetailResponse{
		ID:              dataDB.ID,
		UserID:          dataDB.UserID,
		CashierName:     dataDB.CashierName,
		ShiftName:       dataDB.ShiftName,
		ShiftStart:      dataDB.ShiftStart,
		ShiftEnd:        dataDB.ShiftEnd,
		OpenTime:        dataDB.OpenTime,
		CloseTime:       dataDB.CloseTime,
		OpeningBalance:  dataDB.OpeningBalance,
		ClosingBalance:  dataDB.ClosingBalance,
		ExpectedBalance: dataDB.ExpectedBalance,
		TotalCashSales:  dataDB.TotalCashSales,
		TotalExpenses:   dataDB.TotalExpenses,
		Difference:      dataDB.Difference,
		Status:          dataDB.Status,
		Notes:           dataDB.Notes,
		OpenNotes:       dataDB.OpenNotes,
		Transactions:    trxList,
		Expenses:        expList,
	}

	return data, nil
}

func (s *cashDrawerService) GetHistory(req *dto.GetHistoryRequest) ([]*dto.CashDrawerHistoryResponse, int64, error) {
	return s.repo.GetHistory(req)
}

func (s *cashDrawerService) Open(userID int, req *dto.OpenRequest) (*dto.OpenResponse, error) {
	if req.OpeningBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo awal tidak boleh negatif"}
	}

	existing, err := s.repo.GetOpenCashDrawer(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Sudah ada kas yang terbuka"}
	}

	id, err := s.repo.Open(userID, req.ShiftID, req.OpeningBalance, req.Notes)
	if err != nil {
		return nil, err
	}
	return &dto.OpenResponse{ID: int(id)}, nil
}

func (s *cashDrawerService) Close(id int, req *dto.CloseRequest, requestingUserID int, role string) (*dto.CloseResponse, error) {
	if req.ClosingBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo akhir tidak boleh negatif"}
	}

	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, &errors.NotFoundError{Message: "Kas tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && current.UserID != requestingUserID {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas ini"}
	}
	if current.Status != "open" {
		return nil, &errors.BadRequestError{Message: "Kas sudah ditutup"}
	}

	expected := current.ExpectedBalance
	difference := req.ClosingBalance - expected

	if err := s.repo.Close(id, req.ClosingBalance, expected, difference, req.Notes); err != nil {
		return nil, err
	}

	return &dto.CloseResponse{
		ExpectedBalance: expected,
		ClosingBalance:  req.ClosingBalance,
		Difference:      difference,
	}, nil
}

func (s *cashDrawerService) UpdateSales(id int, req *dto.UpdateSalesRequest, requestingUserID int, role string) error {
	cd, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if cd == nil {
		return &errors.NotFoundError{Message: "Kas tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && cd.UserID != requestingUserID {
		return &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas ini"}
	}
	if cd.Status != "open" {
		return &errors.BadRequestError{Message: "Kas sudah ditutup"}
	}

	return s.repo.UpdateSales(id, req.TotalSales, req.TotalCashSales)
}

func (s *cashDrawerService) AutoCloseYesterday() (int, error) {
	return s.repo.AutoCloseYesterday()
}

func (s *cashDrawerService) UpdateExpenses(id int, req *dto.UpdateExpensesRequest, requestingUserID int, role string) error {
	cd, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if cd == nil {
		return &errors.NotFoundError{Message: "Kas tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && cd.UserID != requestingUserID {
		return &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas ini"}
	}
	if cd.Status != "open" {
		return &errors.BadRequestError{Message: "Kas sudah ditutup"}
	}

	return s.repo.UpdateExpenses(id, req.TotalExpenses)
}
