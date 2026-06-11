package service

import (
	"encoding/json"
	"strings"
	"time"

	"pos_api/domain/sync/dto"
	"pos_api/domain/sync/model"
	"pos_api/errors"
)

func (s *syncService) detectConflict(item *dto.SyncItem) (bool, string, error) {
	if item.ServerID == 0 {
		return false, "", nil
	}

	snapshot, err := s.repo.GetEntitySnapshot(item.EntityType, item.ServerID)
	if err != nil || snapshot == nil {
		return false, "", nil
	}

	onlineTime, err := time.Parse(time.RFC3339, snapshot.UpdatedAt)
	if err != nil {
		return false, "", nil
	}

	desktopTime, err := time.Parse(time.RFC3339, item.UpdatedAt)
	if err != nil {
		return false, "", nil
	}

	if onlineTime.After(desktopTime) {
		return true, snapshot.Data, nil
	}

	return false, "", nil
}

func (s *syncService) PushSync(req *dto.PushSyncRequest) (*dto.PushSyncResponse, error) {
	startedAt := time.Now()
	processed, conflicts, failed := 0, 0, 0
	results := make([]dto.SyncItemResult, 0, len(req.Items))

	for i := range req.Items {
		item := &req.Items[i]

		isConflict, onlineData, _ := s.detectConflict(item)
		if isConflict {
			conflictID, err := s.repo.CreateConflict(req.DeviceID, item, onlineData)
			if err != nil {
				failed++
				results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed"})
				continue
			}
			conflicts++
			results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "conflict", ConflictID: conflictID})
			continue
		}

		if item.EntityType == "transaction" && item.ServerID == 0 {
			serverID, err := s.transactionRepo.ApplySyncTransaction(item.Payload, item.LocalID)
			if err != nil {
				if strings.Contains(err.Error(), "stok produk") {
					conflictID, cerr := s.repo.CreateConflict(req.DeviceID, item, err.Error())
					if cerr != nil {
						failed++
						results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed"})
						continue
					}
					conflicts++
					results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "conflict", ConflictID: conflictID, Message: err.Error()})
				} else {
					failed++
					results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed"})
				}
				continue
			}
			processed++
			results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: serverID})
			continue
		}

		queueID, err := s.repo.CreateQueueItem(req.DeviceID, item)
		if err != nil {
			failed++
			results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed"})
			continue
		}

		_ = s.repo.UpdateQueueStatus(queueID, "synced", "")
		processed++
		results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: item.ServerID})
	}

	resp := &dto.PushSyncResponse{
		Processed: processed,
		Conflicts: conflicts,
		Failed:    failed,
		Results:   results,
	}

	s.saveSyncHistory(req.DeviceID, req.DeviceType, results, startedAt)

	return resp, nil
}

func (s *syncService) saveSyncHistory(deviceID, deviceType string, results []dto.SyncItemResult, startedAt time.Time) {
	synced, conflict, failed := 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case "synced":
			synced++
		case "conflict":
			conflict++
		case "failed":
			failed++
		}
	}

	status := "success"
	if failed > 0 && synced == 0 {
		status = "failed"
	} else if conflict > 0 || failed > 0 {
		status = "partial"
	}

	if deviceType == "" {
		deviceType = "desktop"
	}

	now := time.Now()
	_ = s.repo.InsertHistory(model.SyncHistory{
		DeviceID:      deviceID,
		DeviceType:    deviceType,
		TotalItems:    len(results),
		SyncedItems:   synced,
		ConflictItems: conflict,
		FailedItems:   failed,
		DurationMs:    int(now.Sub(startedAt).Milliseconds()),
		Status:        status,
		StartedAt:     startedAt,
		FinishedAt:    &now,
	})
}

func (s *syncService) GetConflicts(filter *dto.ConflictFilter) (*dto.ConflictListResponse, error) {
	data, total, err := s.repo.GetConflicts(filter)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data konflik"}
	}
	page, limit := filter.Page, filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &dto.ConflictListResponse{Data: data, Total: total, Page: page, Limit: limit}, nil
}

func (s *syncService) CountPendingConflicts() (int, error) {
	return s.repo.CountPendingConflicts()
}

func (s *syncService) ResolveConflict(id, userID int, action string) error {
	conflict, err := s.repo.GetConflictByID(id)
	if err != nil {
		return &errors.NotFoundError{Message: "Konflik tidak ditemukan"}
	}

	if conflict.Status == "resolved" {
		return &errors.BadRequestError{Message: "Konflik sudah diselesaikan"}
	}

	switch action {
	case "approve":
		return s.applyDesktopVersion(conflict, userID)
	case "reject":
		return s.rejectDesktopVersion(conflict, userID)
	default:
		return &errors.BadRequestError{Message: "action tidak valid: gunakan 'approve' atau 'reject'"}
	}
}

func (s *syncService) rejectDesktopVersion(conflict *model.SyncConflict, resolvedBy int) error {
	if conflict.EntityType == "transaction" {
		if err := s.transactionRepo.ReturnStockForRejectSync(conflict.EntityID, resolvedBy); err != nil {
			return &errors.InternalServerError{Message: "Gagal mengembalikan stok transaksi yang ditolak"}
		}
	}
	return s.repo.MarkResolved(conflict.ID, resolvedBy, "reject")
}

func (s *syncService) applyDesktopVersion(conflict *model.SyncConflict, resolvedBy int) error {
	var desktopData map[string]interface{}
	if err := json.Unmarshal([]byte(conflict.DesktopData), &desktopData); err != nil {
		return &errors.BadRequestError{Message: "desktop_data tidak valid JSON"}
	}

	switch conflict.EntityType {
	case "transaction":
		if err := s.transactionRepo.UpdateFromSync(conflict.EntityID, desktopData); err != nil {
			return &errors.InternalServerError{Message: "Gagal menerapkan data transaksi desktop"}
		}
	case "expense":
		if err := s.expenseRepo.UpdateFromSync(conflict.EntityID, desktopData); err != nil {
			return &errors.InternalServerError{Message: "Gagal menerapkan data pengeluaran desktop"}
		}
	}

	return s.repo.MarkResolved(conflict.ID, resolvedBy, "approve")
}

func (s *syncService) GetQueue(filter *dto.QueueFilter) (*dto.QueueListResponse, error) {
	data, total, err := s.repo.GetQueue(filter)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil data antrian sync"}
	}
	return &dto.QueueListResponse{Data: data, Total: total}, nil
}

func (s *syncService) GetHistory(filter *dto.HistoryFilter) (*dto.SyncHistoryListResponse, error) {
	data, total, err := s.repo.GetHistory(filter)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil riwayat sync"}
	}
	page, limit := filter.Page, filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &dto.SyncHistoryListResponse{Data: data, Total: total, Page: page, Limit: limit}, nil
}
