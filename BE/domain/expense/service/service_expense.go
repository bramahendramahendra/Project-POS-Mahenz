package service

import (
	dto "pos_api/domain/expense/dto"
	"pos_api/errors"
)

func (s *expenseService) GetAll(req *dto.GetAllRequest) (data []dto.ExpenseResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.ExpenseResponse{
			ID:            v.ID,
			ExpenseDate:   v.ExpenseDate,
			Category:      v.Category,
			Description:   v.Description,
			Amount:        v.Amount,
			PaymentMethod: v.PaymentMethod,
			UserID:        v.UserID,
			UserName:      v.UserName,
			Notes:         v.Notes,
		})
	}

	return data, total, nil
}

func (s *expenseService) GetByID(id int) (data dto.ExpenseResponse, err error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.NotFoundError{Message: "Pengeluaran tidak ditemukan"}
	}

	data = dto.ExpenseResponse{
		ID:            dataDB.ID,
		ExpenseDate:   dataDB.ExpenseDate,
		Category:      dataDB.Category,
		Description:   dataDB.Description,
		Amount:        dataDB.Amount,
		PaymentMethod: dataDB.PaymentMethod,
		UserID:        dataDB.UserID,
		UserName:      dataDB.UserName,
		Notes:         dataDB.Notes,
	}

	return data, nil
}

func (s *expenseService) Create(req *dto.CreateRequest, userID int) (data dto.ExpenseResponse, err error) {
	newID, err := s.repo.Create(req, userID)
	if err != nil {
		return data, err
	}

	if req.CashDrawerID != nil {
		cd, _ := s.cashDrawerRepo.GetByID(*req.CashDrawerID)
		if cd != nil && cd.Status == "open" {
			_ = s.cashDrawerRepo.UpdateExpenses(cd.ID, req.Amount)
		}
	} else {
		openCashDrawer, _ := s.cashDrawerRepo.GetOpenCashDrawer(userID)
		if openCashDrawer != nil {
			_ = s.cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, req.Amount)
		}
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}
	if dataDB == nil {
		return data, &errors.InternalServerError{Message: "Gagal mengambil data pengeluaran"}
	}

	data = dto.ExpenseResponse{
		ID:            dataDB.ID,
		ExpenseDate:   dataDB.ExpenseDate,
		Category:      dataDB.Category,
		Description:   dataDB.Description,
		Amount:        dataDB.Amount,
		PaymentMethod: dataDB.PaymentMethod,
		UserID:        dataDB.UserID,
		UserName:      dataDB.UserName,
		Notes:         dataDB.Notes,
	}

	return data, nil
}

func (s *expenseService) Update(req *dto.UpdateRequest) (err error) {
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return &errors.NotFoundError{Message: "Pengeluaran tidak ditemukan"}
	}

	if err = s.repo.Update(req); err != nil {
		return err
	}

	delta := req.Amount - existing.Amount
	if delta != 0 {
		openCashDrawer, _ := s.cashDrawerRepo.GetOpenCashDrawer(existing.UserID)
		if openCashDrawer != nil && openCashDrawer.OpenTime.Local().Format("2006-01-02") <= existing.ExpenseDate {
			_ = s.cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, delta)
		}
	}

	return nil
}

func (s *expenseService) Delete(req *dto.DeleteRequest) (err error) {
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return &errors.NotFoundError{Message: "Pengeluaran tidak ditemukan"}
	}

	if err = s.repo.Delete(req); err != nil {
		return err
	}

	openCashDrawer, _ := s.cashDrawerRepo.GetOpenCashDrawer(existing.UserID)
	if openCashDrawer != nil && openCashDrawer.OpenTime.Local().Format("2006-01-02") <= existing.ExpenseDate {
		_ = s.cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, -existing.Amount)
	}

	return nil
}
