package repo

import (
	dto "coffee_api/domain/menu_item/dto"
	model "coffee_api/domain/menu_item/model"

	"gorm.io/gorm"
)

type (
	MenuItemRepoInterface interface {
		GetAll(req *dto.GetAllRequest) ([]*model.MenuItem, int64, error)
		GetByID(id int) (*model.MenuItem, error)
		GetVariants(itemID int) ([]*model.Variant, error)
		GetRecipeLines(variantIDs []int) ([]*model.RecipeLine, error)
		GetModifierGroups(itemID int) ([]*model.ItemModifierGroup, error)

		NameExists(name string, excludeID int) (bool, error)
		CategoryIsActive(id int) (bool, error)
		ProductIsActive(id int) (bool, error)
		ProductLinked(productID int, excludeID int) (bool, error)
		CountActiveSizes(ids []int) (int64, error)
		CountIngredients(ids []int) (int64, error)
		CountActiveModifierGroups(ids []int) (int64, error)

		Create(req *dto.SaveRequest, userID int) (int64, error)
		Update(req *dto.SaveRequest) error
		InsertVariant(itemID int, v *dto.VariantInput) (int64, error)
		UpdateVariant(variantID int, price float64, isActive bool) error
		DeleteVariantsNotIn(itemID int, keepIDs []int) error
		DeleteRecipeLines(variantID int) error
		InsertRecipeLine(variantID int, line *dto.RecipeInput, sortOrder int) error
		ReplaceModifierGroups(itemID int, groupIDs []int) error

		ToggleStatus(id int) error
		SetSoldOut(id int, soldOut bool) error
		Delete(id int) error
		CountOrderItems(menuItemID int) (int64, error)

		// Katalog kasir (tahap 4.7b).
		CatalogData() ([]model.CatalogCategory, []model.CatalogItem, []model.CatalogVariant, []model.CatalogGroup, []model.CatalogOption, error)
		AvailableStock(productIDs []int) (map[int]float64, error)

		GetDB() *gorm.DB
		WithTx(tx *gorm.DB) MenuItemRepoInterface
	}

	menuItemRepo struct {
		db *gorm.DB
	}
)

func NewMenuItemRepo(db *gorm.DB) *menuItemRepo {
	return &menuItemRepo{db: db}
}

func (r *menuItemRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *menuItemRepo) WithTx(tx *gorm.DB) MenuItemRepoInterface {
	return &menuItemRepo{db: tx}
}
