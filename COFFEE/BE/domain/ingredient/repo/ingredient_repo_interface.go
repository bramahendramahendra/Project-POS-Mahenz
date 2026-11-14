package repo

import (
	dto "coffee_api/domain/ingredient/dto"
	model "coffee_api/domain/ingredient/model"

	"gorm.io/gorm"
)

type (
	IngredientRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.Ingredient, int64, error)
		GetOptions() ([]*model.Ingredient, error)
		GetByID(id int) (*model.Ingredient, error)
		// GetByIDForUpdate mengunci baris bahan sampai transaksi selesai (stok dihitung ulang aman).
		GetByIDForUpdate(id int) (*model.Ingredient, error)
		NameExists(name string, excludeID int) (bool, error)
		UnitExists(unitID int) (bool, error)
		Create(req *dto.CreateRequest, userID int) (int64, error)
		Update(req *dto.UpdateRequest) error
		ToggleStatus(id int) error
		Delete(id int) error
		HasMutations(id int) (bool, error)
		UsedInRecipes(id int) (bool, error)

		// ApplyStockChange menambah stok, total terpakai, dan total terbuang sekaligus (boleh negatif).
		ApplyStockChange(id int, stockDelta, usedDelta, wastedDelta float64) error
		CreateMutation(m *model.IngredientMutation) (int64, error)
		GetMutations(req *dto.GetMutationsRequest) ([]*model.IngredientMutation, int64, error)
		GetMutationByID(id int) (*model.IngredientMutation, error)
		GetMutationByIDForUpdate(id int) (*model.IngredientMutation, error)
		MarkMutationCancelled(id, userID int, reason string) error

		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) IngredientRepoInterface
	}

	ingredientRepo struct {
		db *gorm.DB
	}
)

func NewIngredientRepo(db *gorm.DB) *ingredientRepo {
	return &ingredientRepo{db: db}
}

func (r *ingredientRepo) GetDB() *gorm.DB {
	return r.db
}

// WithTx mengembalikan repo yang terikat ke transaksi tx.
func (r *ingredientRepo) WithTx(tx *gorm.DB) IngredientRepoInterface {
	return &ingredientRepo{db: tx}
}
