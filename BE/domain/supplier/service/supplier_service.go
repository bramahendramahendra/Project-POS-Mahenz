package service

import (
	"strings"

	dto "pos_api/domain/supplier/dto"
	"pos_api/errors"
)

func (s *supplierService) GetAll(req *dto.SupplierListRequest) (data []dto.SupplierResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, dto.SupplierResponse{
			ID:            v.ID,
			SupplierCode:  v.SupplierCode,
			Name:          v.Name,
			Address:       v.Address,
			Phone:         v.Phone,
			Email:         v.Email,
			ContactPerson: v.ContactPerson,
			Notes:         v.Notes,
			IsActive:      v.IsActive,
			CreatedAt:     v.CreatedAt,
		})
	}

	return data, total, nil
}

func (s *supplierService) GetOptions() (data []dto.SupplierOptionResponse, err error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.SupplierOptionResponse{
			ID:           v.ID,
			SupplierCode: v.SupplierCode,
			Name:         v.Name,
		})
	}

	return data, nil
}

func (s *supplierService) GetDetail(id int) (data dto.SupplierDetailResponse, err error) {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return data, err
	}
	if supplier == nil {
		return data, &errors.NotFoundError{Message: "Supplier tidak ditemukan"}
	}

	purchasesDB, _ := s.repo.GetPurchaseHistory(id)
	purchases := make([]dto.SupplierPurchaseItem, 0, len(purchasesDB))
	for _, v := range purchasesDB {
		purchases = append(purchases, dto.SupplierPurchaseItem{
			ID:              v.ID,
			PurchaseCode:    v.PurchaseCode,
			PurchaseDate:    v.PurchaseDate,
			TotalAmount:     v.TotalAmount,
			PaymentStatus:   v.PaymentStatus,
			RemainingAmount: v.RemainingAmount,
		})
	}

	returnsDB, _ := s.repo.GetReturnHistory(id)
	returns := make([]dto.SupplierReturnHistoryItem, 0, len(returnsDB))
	for _, v := range returnsDB {
		returns = append(returns, dto.SupplierReturnHistoryItem{
			ID:         v.ID,
			ReturnCode: v.ReturnCode,
			ReturnDate: v.ReturnDate,
			TotalReturn: v.TotalReturn,
			Reason:     v.Reason,
			Status:     v.Status,
		})
	}

	var totalAmount, totalDebt, totalReturnAmount float64
	for _, p := range purchases {
		totalAmount += p.TotalAmount
		totalDebt += p.RemainingAmount
	}
	for _, r := range returns {
		totalReturnAmount += r.TotalReturn
	}

	data = dto.SupplierDetailResponse{
		ID:              supplier.ID,
		SupplierCode:    supplier.SupplierCode,
		Name:            supplier.Name,
		Address:         supplier.Address,
		Phone:           supplier.Phone,
		Email:           supplier.Email,
		ContactPerson:   supplier.ContactPerson,
		Notes:           supplier.Notes,
		IsActive:        supplier.IsActive,
		TotalPurchases:  len(purchases),
		TotalAmount:     totalAmount,
		TotalDebt:       totalDebt,
		TotalReturn:     totalReturnAmount,
		PurchaseHistory: purchases,
		ReturnHistory:   returns,
	}

	return data, nil
}

func (s *supplierService) Create(req *dto.CreateSupplierRequest) (data dto.SupplierResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)

	exists, err := s.repo.CheckNameExists(req.Name, 0)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama supplier sudah digunakan"}
	}

	code, err := s.generateUniqueCode()
	if err != nil {
		return data, err
	}

	newID, err := s.repo.Create(req, code)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(int(newID))
	if err != nil {
		return data, err
	}

	data = dto.SupplierResponse{
		ID:            dataDB.ID,
		SupplierCode:  dataDB.SupplierCode,
		Name:          dataDB.Name,
		Address:       dataDB.Address,
		Phone:         dataDB.Phone,
		Email:         dataDB.Email,
		ContactPerson: dataDB.ContactPerson,
		Notes:         dataDB.Notes,
		IsActive:      dataDB.IsActive,
		CreatedAt:     dataDB.CreatedAt,
	}

	return data, nil
}

func (s *supplierService) Update(req *dto.UpdateSupplierRequest) (data dto.SupplierResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)

	existsUpdate, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}
	if existsUpdate == nil {
		return data, &errors.NotFoundError{Message: "Supplier tidak ditemukan"}
	}

	exists, err := s.repo.CheckNameExists(req.Name, req.ID)
	if err != nil {
		return data, err
	}
	if exists {
		return data, &errors.BadRequestError{Message: "Nama supplier sudah digunakan"}
	}

	err = s.repo.Update(req)
	if err != nil {
		return data, err
	}

	dataDB, err := s.repo.GetByID(req.ID)
	if err != nil {
		return data, err
	}

	data = dto.SupplierResponse{
		ID:            dataDB.ID,
		SupplierCode:  dataDB.SupplierCode,
		Name:          dataDB.Name,
		Address:       dataDB.Address,
		Phone:         dataDB.Phone,
		Email:         dataDB.Email,
		ContactPerson: dataDB.ContactPerson,
		Notes:         dataDB.Notes,
		IsActive:      dataDB.IsActive,
		CreatedAt:     dataDB.CreatedAt,
	}

	return data, nil
}

func (s *supplierService) Delete(req *dto.DeleteSupplierRequest) error {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Supplier tidak ditemukan"}
	}

	count, err := s.repo.CountPurchasesBySupplier(req.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Gagal memeriksa penggunaan supplier"}
	}
	if count > 0 {
		return &errors.BadRequestError{Message: "Supplier tidak bisa dihapus karena sudah ada Purchase Order"}
	}

	return s.repo.Delete(req)
}

func (s *supplierService) ToggleStatus(req *dto.ToggleStatusSupplierRequest) error {
	exists, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if exists == nil {
		return &errors.NotFoundError{Message: "Supplier tidak ditemukan"}
	}

	if exists.IsActive {
		debtCount, err := s.repo.CountActiveDebtBySupplier(req.ID)
		if err != nil {
			return &errors.InternalServerError{Message: "Gagal memeriksa hutang aktif supplier"}
		}
		if debtCount > 0 {
			return &errors.BadRequestError{Message: "Supplier tidak bisa dinonaktifkan karena masih memiliki hutang yang belum lunas"}
		}
	}

	return s.repo.ToggleStatus(req)
}
