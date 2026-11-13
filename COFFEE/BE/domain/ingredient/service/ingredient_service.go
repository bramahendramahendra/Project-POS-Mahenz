package service

import (
	"math"
	"strings"

	"coffee_api/config"
	dto "coffee_api/domain/ingredient/dto"
	model "coffee_api/domain/ingredient/model"
	"coffee_api/errors"
	time_helper "coffee_api/helper/time"

	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Katalog bahan
// ---------------------------------------------------------------------------

func (s *ingredientService) GetAll(req *dto.GetAllRequest) ([]dto.IngredientResponse, int64, error) {
	dataDB, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, err
	}
	data := make([]dto.IngredientResponse, 0, len(dataDB))
	for _, v := range dataDB {
		data = append(data, toIngredientResponse(v))
	}
	return data, total, nil
}

func (s *ingredientService) GetOptions() ([]dto.IngredientOptionResponse, error) {
	dataDB, err := s.repo.GetOptions()
	if err != nil {
		return nil, err
	}
	data := make([]dto.IngredientOptionResponse, 0, len(dataDB))
	for _, v := range dataDB {
		data = append(data, dto.IngredientOptionResponse{ID: v.ID, Name: v.Name, UnitAbbr: v.UnitAbbr, Stock: v.Stock})
	}
	return data, nil
}

func (s *ingredientService) GetByID(id int) (*dto.IngredientResponse, error) {
	dataDB, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if dataDB == nil {
		return nil, &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
	}
	res := toIngredientResponse(dataDB)
	return &res, nil
}

// Create: stok bahan baru selalu 0 (keputusan 83), stok bertambah lewat pembelian (tahap 4.5).
func (s *ingredientService) Create(req *dto.CreateRequest, userID int) (*dto.IngredientResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if err := s.validateNameAndUnit(req.Name, req.UnitID, 0); err != nil {
		return nil, err
	}
	id, err := s.repo.Create(req, userID)
	if err != nil {
		return nil, err
	}
	return s.GetByID(int(id))
}

// Update: satuan tidak bisa diubah bila bahan sudah punya riwayat stok, karena angka stok
// dan riwayat lama tercatat dalam satuan lama.
func (s *ingredientService) Update(req *dto.UpdateRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	current, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}
	if current == nil {
		return &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
	}
	if err := s.validateNameAndUnit(req.Name, req.UnitID, req.ID); err != nil {
		return err
	}
	if req.UnitID != current.UnitID {
		hasHistory, err := s.repo.HasMutations(req.ID)
		if err != nil {
			return err
		}
		if hasHistory {
			return &errors.BadRequestError{Message: "Satuan tidak bisa diubah karena bahan sudah punya riwayat stok"}
		}
	}
	return s.repo.Update(req)
}

func (s *ingredientService) ToggleStatus(id int) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
	}
	return s.repo.ToggleStatus(id)
}

// Delete: hanya untuk bahan tanpa riwayat stok (mis. salah buat). Bahan yang sudah dipakai
// cukup dinonaktifkan agar riwayatnya tetap utuh.
func (s *ingredientService) Delete(id int) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
	}
	hasHistory, err := s.repo.HasMutations(id)
	if err != nil {
		return err
	}
	if hasHistory {
		return &errors.BadRequestError{Message: "Bahan sudah punya riwayat stok dan tidak bisa dihapus. Nonaktifkan saja"}
	}
	return s.repo.Delete(id)
}

func (s *ingredientService) validateNameAndUnit(name string, unitID, excludeID int) error {
	exists, err := s.repo.NameExists(name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return &errors.BadRequestError{Message: "Nama bahan sudah digunakan"}
	}
	unitOK, err := s.repo.UnitExists(unitID)
	if err != nil {
		return err
	}
	if !unitOK {
		return &errors.BadRequestError{Message: "Satuan tidak ditemukan atau tidak aktif"}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Pemakaian, barang terbuang, pembatalan
// ---------------------------------------------------------------------------

// RecordUsage mencatat bahan terpakai. Stok minus diizinkan dengan peringatan (keputusan 14.9).
func (s *ingredientService) RecordUsage(req *dto.RecordUsageRequest, userID int) (*dto.RecordResponse, error) {
	qty := roundQty(req.Qty)
	return s.record(req.IngredientID, dto.MutationUsage, qty, strings.TrimSpace(req.Notes), userID)
}

// RecordWaste mencatat bahan terbuang, alasan wajib (keputusan 83).
func (s *ingredientService) RecordWaste(req *dto.RecordWasteRequest, userID int) (*dto.RecordResponse, error) {
	qty := roundQty(req.Qty)
	return s.record(req.IngredientID, dto.MutationWaste, qty, strings.TrimSpace(req.Reason), userID)
}

func (s *ingredientService) record(ingredientID int, mutationType string, qty float64, notes string, userID int) (*dto.RecordResponse, error) {
	if qty <= 0 {
		return nil, &errors.BadRequestError{Message: "Jumlah harus lebih dari 0 (maksimal 3 angka di belakang koma)"}
	}

	var mutationID int64
	var stockAfter float64
	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		ing, err := r.GetByIDForUpdate(ingredientID)
		if err != nil {
			return err
		}
		if ing == nil {
			return &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
		}
		if !ing.IsActive {
			return &errors.BadRequestError{Message: "Bahan nonaktif tidak bisa dicatat pemakaiannya"}
		}

		used, wasted := qty, 0.0
		if mutationType == dto.MutationWaste {
			used, wasted = 0, qty
		}
		stockAfter = roundQty(ing.Stock - qty)
		if err := r.ApplyStockChange(ing.ID, -qty, used, wasted); err != nil {
			return err
		}
		mutationID, err = r.CreateMutation(&model.IngredientMutation{
			IngredientID: ing.ID,
			Type:         mutationType,
			QtyChange:    -qty,
			StockBefore:  ing.Stock,
			StockAfter:   stockAfter,
			Notes:        notes,
			CreatedBy:    userID,
			CreatedAt:    time_helper.GetTimeNow(),
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	m, err := s.repo.GetMutationByID(int(mutationID))
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil riwayat stok bahan"}
	}
	return &dto.RecordResponse{
		Mutation:      toMutationResponse(m),
		StockAfter:    stockAfter,
		NegativeStock: stockAfter < 0,
	}, nil
}

func (s *ingredientService) GetMutations(req *dto.GetMutationsRequest) ([]dto.MutationResponse, int64, error) {
	dataDB, total, err := s.repo.GetMutations(req)
	if err != nil {
		return nil, 0, err
	}
	data := make([]dto.MutationResponse, 0, len(dataDB))
	for _, v := range dataDB {
		data = append(data, toMutationResponse(v))
	}
	return data, total, nil
}

// CancelMutation membatalkan entri terpakai/terbuang (keputusan 81): stok dan total dikembalikan,
// entri asli ditandai dibatalkan, dan pembatalan dicatat sebagai entri baru (jejak utuh).
func (s *ingredientService) CancelMutation(req *dto.CancelMutationRequest, userID int) (*dto.MutationResponse, error) {
	reason := strings.TrimSpace(req.Reason)
	var cancelID int64
	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		r := s.repo.WithTx(tx)
		orig, err := r.GetMutationByIDForUpdate(req.ID)
		if err != nil {
			return err
		}
		if orig == nil {
			return &errors.NotFoundError{Message: "Riwayat stok tidak ditemukan"}
		}
		if orig.Type != dto.MutationUsage && orig.Type != dto.MutationWaste {
			return &errors.BadRequestError{Message: "Hanya entri terpakai atau terbuang yang bisa dibatalkan"}
		}
		if orig.CancelledAt != nil {
			return &errors.BadRequestError{Message: "Entri ini sudah dibatalkan"}
		}

		ing, err := r.GetByIDForUpdate(orig.IngredientID)
		if err != nil {
			return err
		}
		if ing == nil {
			return &errors.NotFoundError{Message: "Bahan tidak ditemukan"}
		}

		qty := -orig.QtyChange // entri asal bernilai negatif
		used, wasted := -qty, 0.0
		if orig.Type == dto.MutationWaste {
			used, wasted = 0, -qty
		}
		if err := r.ApplyStockChange(ing.ID, qty, used, wasted); err != nil {
			return err
		}
		if err := r.MarkMutationCancelled(orig.ID, userID, reason); err != nil {
			return err
		}
		refID := orig.ID
		cancelID, err = r.CreateMutation(&model.IngredientMutation{
			IngredientID:  ing.ID,
			Type:          dto.MutationCancel,
			QtyChange:     qty,
			StockBefore:   ing.Stock,
			StockAfter:    roundQty(ing.Stock + qty),
			Notes:         reason,
			ReferenceType: "ingredient_mutation",
			ReferenceID:   &refID,
			CreatedBy:     userID,
			CreatedAt:     time_helper.GetTimeNow(),
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	m, err := s.repo.GetMutationByID(int(cancelID))
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, &errors.InternalServerError{Message: "Gagal mengambil riwayat pembatalan"}
	}
	res := toMutationResponse(m)
	return &res, nil
}

// ---------------------------------------------------------------------------
// Pemetaan
// ---------------------------------------------------------------------------

// roundQty membulatkan ke 3 angka desimal (sama dengan kolom DECIMAL(15,3)).
func roundQty(v float64) float64 {
	return math.Round(v*1000) / 1000
}

func toIngredientResponse(v *model.Ingredient) dto.IngredientResponse {
	return dto.IngredientResponse{
		ID:          v.ID,
		Name:        v.Name,
		UnitID:      v.UnitID,
		UnitName:    v.UnitName,
		UnitAbbr:    v.UnitAbbr,
		Stock:       v.Stock,
		TotalUsed:   v.TotalUsed,
		TotalWasted: v.TotalWasted,
		MinStock:    v.MinStock,
		IsLowStock:  v.MinStock != nil && v.Stock <= *v.MinStock,
		IsActive:    v.IsActive,
		HasHistory:  v.HasHistory,
	}
}

func toMutationResponse(v *model.IngredientMutation) dto.MutationResponse {
	res := dto.MutationResponse{
		ID:             v.ID,
		IngredientID:   v.IngredientID,
		IngredientName: v.IngredientName,
		UnitAbbr:       v.UnitAbbr,
		Type:           v.Type,
		QtyChange:      v.QtyChange,
		StockBefore:    v.StockBefore,
		StockAfter:     v.StockAfter,
		Notes:          v.Notes,
		ReferenceID:    v.ReferenceID,
		CreatedByName:  v.CreatedByName,
		CreatedAt:      v.CreatedAt.Format(config.FormatTime),
	}
	if v.CancelledAt != nil {
		res.IsCancelled = true
		res.CancelledAt = v.CancelledAt.Format(config.FormatTime)
		res.CancelledByName = v.CancelledByName
		res.CancelReason = v.CancelReason
	}
	return res
}
