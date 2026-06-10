package service

import (
	dto "pos_api/domain/cash_drawer/dto"
	"pos_api/errors"
)

func (s *cashDrawerService) GetCurrent(userID int) (*dto.CurrentCashDrawerResponse, error) {
	res, err := s.repo.GetCurrent(userID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	return res, nil
}

func (s *cashDrawerService) GetByID(id int, requestingUserID int, role string) (*dto.CashDrawerDetailResponse, error) {
	res, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if res == nil {
		return nil, &errors.NotFoundError{Message: "Kas tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && res.UserID != requestingUserID {
		return nil, &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke kas ini"}
	}
	return res, nil
}

func (s *cashDrawerService) GetHistory(req *dto.GetHistoryRequest) (data []*dto.CashDrawerHistoryResponse, total int64, err error) {
	data, total, err = s.repo.GetHistory(req)
	if err != nil {
		return nil, 0, &errors.InternalServerError{Message: err.Error()}
	}
	return data, total, nil
}

func (s *cashDrawerService) Open(userID int, req *dto.OpenRequest) (*dto.OpenResponse, error) {
	if req.OpeningBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo awal tidak boleh negatif"}
	}

	existing, err := s.repo.GetOpenCashDrawer(userID)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	if existing != nil {
		return nil, &errors.BadRequestError{Message: "Sudah ada kas yang terbuka"}
	}

	id, err := s.repo.Open(userID, req.ShiftID, req.OpeningBalance, req.Notes)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
	}
	return &dto.OpenResponse{ID: int(id)}, nil
}

func (s *cashDrawerService) Close(id int, req *dto.CloseRequest, requestingUserID int, role string) (*dto.CloseResponse, error) {
	if req.ClosingBalance < 0 {
		return nil, &errors.BadRequestError{Message: "Saldo akhir tidak boleh negatif"}
	}

	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &errors.InternalServerError{Message: err.Error()}
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
		return nil, &errors.InternalServerError{Message: err.Error()}
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
		return &errors.InternalServerError{Message: err.Error()}
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

func (s *cashDrawerService) UpdateExpenses(id int, req *dto.UpdateExpensesRequest, requestingUserID int, role string) error {
	cd, err := s.repo.GetByID(id)
	if err != nil {
		return &errors.InternalServerError{Message: err.Error()}
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
