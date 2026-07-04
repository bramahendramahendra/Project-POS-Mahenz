package service

import (
	dto "pos_api/domain/expense/dto"
	"pos_api/errors"

	"gorm.io/gorm"
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
	var newID int64

	txErr := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		expenseRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		id, err := expenseRepo.Create(req, userID)
		if err != nil {
			return err
		}
		newID = id

		if req.CashDrawerID != nil {
			cd, err := cashDrawerRepo.GetByID(*req.CashDrawerID)
			if err != nil {
				return err
			}
			if cd != nil && cd.Status == "open" {
				if err := cashDrawerRepo.UpdateExpenses(cd.ID, req.Amount); err != nil {
					return err
				}
			}
		} else {
			openCashDrawer, err := cashDrawerRepo.GetOpenCashDrawer(userID)
			if err != nil {
				return err
			}
			if openCashDrawer != nil {
				if err := cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, req.Amount); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if txErr != nil {
		return data, txErr
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

func (s *expenseService) Update(req *dto.UpdateRequest, requestingUserID int, role string) (err error) {
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return &errors.NotFoundError{Message: "Pengeluaran tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && existing.UserID != requestingUserID {
		return &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke pengeluaran ini"}
	}

	delta := req.Amount - existing.Amount

	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		expenseRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		if err := expenseRepo.Update(req); err != nil {
			return err
		}

		if delta != 0 {
			openCashDrawer, err := cashDrawerRepo.GetOpenCashDrawer(existing.UserID)
			if err != nil {
				return err
			}
			if openCashDrawer != nil && openCashDrawer.OpenTime.Local().Format("2006-01-02") <= existing.ExpenseDate {
				if err := cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, delta); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *expenseService) Delete(req *dto.DeleteRequest, requestingUserID int, role string) (err error) {
	existing, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return &errors.NotFoundError{Message: "Pengeluaran tidak ditemukan"}
	}
	if role != "owner" && role != "admin" && existing.UserID != requestingUserID {
		return &errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke pengeluaran ini"}
	}

	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		expenseRepo := s.repo.WithTx(tx)
		cashDrawerRepo := s.cashDrawerRepo.WithTx(tx)

		if err := expenseRepo.Delete(req); err != nil {
			return err
		}

		openCashDrawer, err := cashDrawerRepo.GetOpenCashDrawer(existing.UserID)
		if err != nil {
			return err
		}
		if openCashDrawer != nil && openCashDrawer.OpenTime.Local().Format("2006-01-02") <= existing.ExpenseDate {
			if err := cashDrawerRepo.UpdateExpenses(openCashDrawer.ID, -existing.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}
