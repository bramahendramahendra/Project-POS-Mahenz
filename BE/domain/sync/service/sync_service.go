package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cash_drawer_dto "pos_api/domain/cash_drawer/dto"
	"pos_api/domain/sync/dto"
	"pos_api/domain/sync/model"
	"pos_api/errors"
	request_helper "pos_api/helper/request"
	"pos_api/pkg/pricing"
	"pos_api/pkg/syncmap"
)

func (s *syncService) detectConflict(item *dto.SyncItem) (bool, string, error) {
	if item.ServerID == 0 {
		return false, "", nil
	}

	// cash_drawer sengaja dilewati di sini: updated_at-nya berubah terus lewat aktivitas
	// normal (tiap penjualan tunai memicu UpdateSales, lihat ApplySyncTransaction), jadi
	// perbandingan timestamp generik ini akan salah-positif mendeteksi konflik untuk setiap
	// batch offline yang berisi transaksi lalu tutup-kas. Deteksi konflik nyata untuk
	// cash_drawer ditangani terpisah di applySyncCashDrawer berdasarkan perbandingan nilai.
	if item.EntityType == "cash_drawer" {
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
	processed, conflicts, failed, pending := 0, 0, 0, 0
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
			recalculatedPayload, err := s.recalculateSyncTransactionPayload(item.Payload)
			if err != nil {
				failed++
				results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed", Message: err.Error()})
				continue
			}

			serverID, err := s.transactionRepo.ApplySyncTransaction(recalculatedPayload, req.DeviceID, item.LocalID, s.cashDrawerRepo)
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
			// Tetap dicatat ke sync_queue untuk jejak audit (sebelumnya transaksi sama
			// sekali tidak tercatat di sini), walau penerapannya sudah terjadi langsung
			// di atas (bukan ditunda) — makanya statusnya langsung 'synced', bukan 'pending'.
			_, _ = s.repo.CreateQueueItem(req.DeviceID, item, "synced")

			processed++
			results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: serverID})
			continue
		}

		if item.EntityType == "cash_drawer" {
			result, err := s.applySyncCashDrawer(req.DeviceID, item)
			if err != nil {
				failed++
				results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed", Message: err.Error()})
				continue
			}
			if result.Status == "conflict" {
				conflicts++
				results = append(results, *result)
				continue
			}
			_, _ = s.repo.CreateQueueItem(req.DeviceID, item, "synced")
			processed++
			results = append(results, *result)
			continue
		}

		if _, err := s.repo.CreateQueueItem(req.DeviceID, item, "pending"); err != nil {
			failed++
			results = append(results, dto.SyncItemResult{LocalID: item.LocalID, Status: "failed"})
			continue
		}

		// Apply-logic untuk entity non-transaksi (product/customer/stock, dst) belum
		// diimplementasikan. Item disimpan di sync_queue dengan status default 'pending'
		// (bukan ditandai 'synced') supaya client tidak menganggap perubahan sudah
		// diterapkan ke tabel target padahal belum ada kode yang menerapkannya.
		pending++
		results = append(results, dto.SyncItemResult{
			LocalID:  item.LocalID,
			Status:   "pending",
			ServerID: item.ServerID,
			Message:  "Diterima dan diantrekan, menunggu implementasi apply-logic ke tabel target",
		})
	}

	resp := &dto.PushSyncResponse{
		Processed: processed,
		Conflicts: conflicts,
		Failed:    failed,
		Pending:   pending,
		Results:   results,
	}

	s.saveSyncHistory(req.DeviceID, req.DeviceType, results, startedAt)

	return resp, nil
}

// recalculateSyncTransactionPayload menghitung ulang subtotal/total transaksi dari payload
// sync offline menggunakan harga produk asli di master data (bukan nilai mentah dari
// client), lalu mengembalikan payload JSON yang sudah diperbarui untuk diteruskan ke
// ApplySyncTransaction. Reuse logic yang sama dengan checkout langsung (pkg/pricing).
func (s *syncService) recalculateSyncTransactionPayload(payload string) (string, error) {
	var tx dto.SyncTransactionPayload
	if err := json.Unmarshal([]byte(payload), &tx); err != nil {
		return "", &errors.BadRequestError{Message: "Payload transaksi tidak valid"}
	}

	items := make([]pricing.Item, len(tx.Items))
	for i, item := range tx.Items {
		items[i] = pricing.Item{
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			UnitID:       item.UnitID,
			Quantity:     item.Quantity,
			DiscountItem: item.DiscountItem,
		}
	}

	totals, err := pricing.Recalculate(s.productRepo, items, tx.Discount, tx.Tax)
	if err != nil {
		return "", err
	}

	for i := range tx.Items {
		tx.Items[i].Price = totals.ItemPrices[i]
		tx.Items[i].Subtotal = totals.ItemSubtotals[i]
	}
	tx.Subtotal = totals.Subtotal
	tx.TotalAmount = totals.TotalAmount

	updated, err := json.Marshal(tx)
	if err != nil {
		return "", &errors.InternalServerError{Message: err.Error()}
	}
	return string(updated), nil
}

// applySyncCashDrawer menerapkan item sync entity_type="cash_drawer". ShiftID di payload
// selalu berupa ID shift master data yang sudah pasti valid (shift tidak pernah dibuat
// offline), jadi tidak perlu resolve-ID lintas-entity — hanya dedupe/idempotency biasa
// lewat sync_id_map (pola yang sama dengan transaksi).
func (s *syncService) applySyncCashDrawer(deviceID string, item *dto.SyncItem) (*dto.SyncItemResult, error) {
	var payload dto.SyncCashDrawerPayload
	if err := json.Unmarshal([]byte(item.Payload), &payload); err != nil {
		return nil, fmt.Errorf("payload kas harian tidak valid: %w", err)
	}

	switch item.Action {
	case "create":
		if existingID, found, err := syncmap.Resolve(s.cashDrawerRepo.GetDB(), deviceID, item.LocalID, "cash_drawer"); err == nil && found {
			return &dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: existingID}, nil
		}

		resp, err := s.cashDrawerSvc.Open(payload.UserID, &cash_drawer_dto.OpenRequest{
			ShiftID:        payload.ShiftID,
			OpeningBalance: payload.OpeningBalance,
			Notes:          payload.Notes,
		})
		if err != nil {
			return nil, err
		}
		if err := syncmap.Record(s.cashDrawerRepo.GetDB(), deviceID, item.LocalID, "cash_drawer", resp.ID); err != nil {
			return nil, err
		}
		return &dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: resp.ID}, nil

	case "update":
		targetID := item.ServerID
		if targetID == 0 {
			resolvedID, found, err := syncmap.Resolve(s.cashDrawerRepo.GetDB(), deviceID, item.LocalID, "cash_drawer")
			if err != nil || !found {
				return nil, fmt.Errorf("kas harian referensi tidak ditemukan/belum tersinkron")
			}
			targetID = resolvedID
		}

		_, closeErr := s.cashDrawerSvc.Close(targetID, &cash_drawer_dto.CloseRequest{
			ClosingBalance: payload.ClosingBalance,
			Notes:          payload.Notes,
		}, payload.UserID, "")
		if closeErr == nil {
			return &dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: targetID}, nil
		}
		if !strings.Contains(closeErr.Error(), "sudah ditutup") {
			return nil, closeErr
		}

		// Kas sudah ditutup -- bisa jadi retry (idempotent, aman) ATAU kas yang sama
		// ditutup device lain dengan data berbeda (konflik nyata, bukan sekadar retry).
		// Dibedakan dengan membandingkan closing_balance yang tersimpan vs yang dikirim
		// ulang -- BUKAN pakai timestamp (lihat alasan di sync_repo.go entityTable()).
		current, err := s.cashDrawerRepo.GetByID(targetID)
		if err != nil || current == nil {
			return nil, fmt.Errorf("kas harian tidak ditemukan setelah gagal ditutup")
		}
		if current.ClosingBalance != nil && *current.ClosingBalance == payload.ClosingBalance {
			return &dto.SyncItemResult{LocalID: item.LocalID, Status: "synced", ServerID: targetID}, nil
		}

		onlineData, _ := s.repo.GetRawByEntityAndID("cash_drawer", targetID)
		conflictID, cerr := s.repo.CreateConflict(deviceID, item, onlineData)
		if cerr != nil {
			return nil, cerr
		}
		return &dto.SyncItemResult{
			LocalID:    item.LocalID,
			Status:     "conflict",
			ConflictID: conflictID,
			Message:    "Kas harian sudah ditutup dengan data berbeda",
		}, nil

	default:
		return nil, fmt.Errorf("action tidak dikenal untuk cash_drawer: %s", item.Action)
	}
}

func (s *syncService) saveSyncHistory(deviceID, deviceType string, results []dto.SyncItemResult, startedAt time.Time) {
	synced, conflict, failed, pending := 0, 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case "synced":
			synced++
		case "conflict":
			conflict++
		case "failed":
			failed++
		case "pending":
			pending++
		}
	}

	status := "success"
	if failed > 0 && synced == 0 && pending == 0 {
		status = "failed"
	} else if conflict > 0 || failed > 0 || pending > 0 {
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
		PendingItems:  pending,
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
	page, limit, _ := request_helper.NormalizePagination(filter.Page, filter.Limit, 20, 0)
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
	page, limit, _ := request_helper.NormalizePagination(filter.Page, filter.Limit, 20, 0)
	return &dto.SyncHistoryListResponse{Data: data, Total: total, Page: page, Limit: limit}, nil
}
