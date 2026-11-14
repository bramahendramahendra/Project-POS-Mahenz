package repo

import (
	"fmt"
	"math"

	model "coffee_api/domain/ingredient/model"
	"coffee_api/errors"
	time_helper "coffee_api/helper/time"

	"gorm.io/gorm"
)

// Jenis entri riwayat dari modul lain (lihat migrasi 003 dan 005).
const (
	MutationPurchase     = "purchase"      // pembelian dari supplier (+)
	MutationPurchaseVoid = "purchase_void" // pembelian di-void atau dikurangi saat edit (-)
)

// ApplyStockParams: perubahan stok bahan dari modul lain (mis. pembelian), di dalam transaksi pemanggil.
type ApplyStockParams struct {
	IngredientID  int
	QtyChange     float64 // positif = stok bertambah, negatif = berkurang
	Type          string
	Notes         string
	ReferenceType string
	ReferenceID   int
	UserID        int
}

// ApplyStock mengunci baris bahan, mengubah stok, dan mencatat riwayat. Stok boleh minus
// (model lentur, keputusan 24 dan 93). Total terpakai/terbuang tidak berubah.
func ApplyStock(tx *gorm.DB, p ApplyStockParams) error {
	qty := math.Round(p.QtyChange*1000) / 1000
	if qty == 0 {
		return nil
	}
	r := &ingredientRepo{db: tx}
	ing, err := r.GetByIDForUpdate(p.IngredientID)
	if err != nil {
		return err
	}
	if ing == nil {
		return &errors.BadRequestError{Message: fmt.Sprintf("Bahan ID %d tidak ditemukan", p.IngredientID)}
	}
	if err := r.ApplyStockChange(ing.ID, qty, 0, 0); err != nil {
		return err
	}
	refID := p.ReferenceID
	_, err = r.CreateMutation(&model.IngredientMutation{
		IngredientID:  ing.ID,
		Type:          p.Type,
		QtyChange:     qty,
		StockBefore:   ing.Stock,
		StockAfter:    math.Round((ing.Stock+qty)*1000) / 1000,
		Notes:         p.Notes,
		ReferenceType: p.ReferenceType,
		ReferenceID:   &refID,
		CreatedBy:     p.UserID,
		CreatedAt:     time_helper.GetTimeNow(),
	})
	return err
}
