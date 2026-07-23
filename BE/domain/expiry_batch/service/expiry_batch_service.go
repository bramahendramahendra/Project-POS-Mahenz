package service

import (
	"time"

	dto "pos_api/domain/expiry_batch/dto"
	"pos_api/errors"
	time_helper "pos_api/helper/time"
)

func severityOf(expiredDateStr string, today string) (severity string, daysLeft int) {
	expired, err1 := time.Parse("2006-01-02", expiredDateStr)
	now, err2 := time.Parse("2006-01-02", today)
	if err1 != nil || err2 != nil {
		return "near", 0
	}
	daysLeft = int(expired.Sub(now).Hours() / 24)
	if daysLeft < 0 {
		return "expired", daysLeft
	}
	return "near", daysLeft
}

func (s *expiryBatchService) GetWarnings(req *dto.GetWarningsRequest) (data []dto.WarningResponse, productSeverity []dto.ProductSeverityResponse, err error) {
	dataDB, err := s.repo.GetWarnings(req.Search)
	if err != nil {
		return data, productSeverity, err
	}

	today := time_helper.GetTimeNow().Format("2006-01-02")
	severityByProduct := make(map[int]string)
	countByProduct := make(map[int]int)
	productOrder := make([]int, 0)

	for _, v := range dataDB {
		expiredDateStr := v.ExpiredDate.Format("2006-01-02")
		severity, daysLeft := severityOf(expiredDateStr, today)

		data = append(data, dto.WarningResponse{
			ID:          v.ID,
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			Qty:         v.Qty,
			ExpiredDate: expiredDateStr,
			Severity:    severity,
			DaysLeft:    daysLeft,
		})

		if _, seen := severityByProduct[v.ProductID]; !seen {
			productOrder = append(productOrder, v.ProductID)
		}
		// "expired" selalu menang atas "near" kalau produk itu punya lebih dari 1 batch warning.
		if severityByProduct[v.ProductID] != "expired" {
			severityByProduct[v.ProductID] = severity
		}
		countByProduct[v.ProductID]++
	}

	for _, productID := range productOrder {
		productSeverity = append(productSeverity, dto.ProductSeverityResponse{
			ProductID:    productID,
			Severity:     severityByProduct[productID],
			WarningCount: countByProduct[productID],
		})
	}

	return data, productSeverity, nil
}

func (s *expiryBatchService) GetByProduct(productID int) (data []dto.BatchHistoryResponse, err error) {
	dataDB, err := s.repo.GetByProduct(productID)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.BatchHistoryResponse{
			ID:          v.ID,
			Qty:         v.Qty,
			ExpiredDate: v.ExpiredDate.Format("2006-01-02"),
			Status:      v.Status,
		})
	}

	return data, nil
}

func (s *expiryBatchService) Confirm(req *dto.ConfirmRequest) (err error) {
	batch, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if batch == nil {
		return &errors.NotFoundError{Message: "Batch expired tidak ditemukan"}
	}
	if batch.Status != "active" {
		return &errors.BadRequestError{Message: "Batch ini sudah diproses sebelumnya"}
	}

	return s.repo.Confirm(req.ID, req.UserID, req.Notes)
}

func (s *expiryBatchService) WriteOff(req *dto.WriteOffRequest) (err error) {
	batch, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if batch == nil {
		return &errors.NotFoundError{Message: "Batch expired tidak ditemukan"}
	}
	if batch.Status != "active" {
		return &errors.BadRequestError{Message: "Batch ini sudah diproses sebelumnya"}
	}

	return s.repo.WriteOff(req.ID, req.UserID, req.Notes)
}
