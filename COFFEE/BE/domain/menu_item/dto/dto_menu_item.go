package dto

const (
	StockModeRecipe = "recipe" // racikan: stok bahan lewat pemakaian manual
	StockModeDirect = "direct" // langsung / barang jadi: stok di tahap 4.5

	TemperatureHot  = "hot"
	TemperatureIced = "iced"
)

type (
	// REQUEST
	GetAllRequest struct {
		Page       int    `json:"page" validate:"required,min=1"`
		Limit      int    `json:"limit" validate:"required,min=1"`
		Search     string `json:"search" validate:"max=100"`
		CategoryID *int   `json:"category_id"`
		StockMode  string `json:"stock_mode" validate:"omitempty,oneof=recipe direct"`
		IsActive   *bool  `json:"is_active"`
		IsSoldOut  *bool  `json:"is_sold_out"`
		SortBy     string `json:"sort_by"`
		SortOrder  string `json:"sort_order"`
	}

	IDUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	RecipeInput struct {
		IngredientID int    `json:"ingredient_id" validate:"required,min=1"`
		Amount       string `json:"amount" validate:"required,max=50,notblank"`
	}

	// VariantInput: satu kombinasi. SizeID dan Temperature boleh kosong, tetapi tidak keduanya.
	VariantInput struct {
		SizeID      *int          `json:"size_id" validate:"omitempty,min=1"`
		Temperature *string       `json:"temperature" validate:"omitempty,oneof=hot iced"`
		Price       float64       `json:"price" validate:"gte=0,lte=999999999"`
		IsActive    bool          `json:"is_active"`
		Recipe      []RecipeInput `json:"recipe" validate:"max=30,dive"`
	}

	// SaveRequest: tanpa varian -> Price + Recipe. Dengan varian -> Variants (masing-masing
	// dengan harga dan resepnya sendiri).
	SaveRequest struct {
		ID               int            `json:"-"`
		CategoryID       int            `json:"category_id" validate:"required,min=1"`
		Name             string         `json:"name" validate:"required,max=100,notblank"`
		StockMode        string         `json:"stock_mode" validate:"required,oneof=recipe direct"`
		ProductID        *int           `json:"product_id" validate:"omitempty,min=1"`
		HasVariants      bool           `json:"has_variants"`
		Price            *float64       `json:"price" validate:"omitempty,gte=0,lte=999999999"`
		Recipe           []RecipeInput  `json:"recipe" validate:"max=30,dive"`
		Variants         []VariantInput `json:"variants" validate:"max=60,dive"`
		ModifierGroupIDs []int          `json:"modifier_group_ids" validate:"max=20,dive,min=1"`
	}

	SoldOutRequest struct {
		ID        int   `json:"-"`
		IsSoldOut *bool `json:"is_sold_out" validate:"required"`
	}

	// RESPONSE
	MenuItemResponse struct {
		ID           int     `json:"id"`
		CategoryID   int     `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Name         string  `json:"name"`
		StockMode    string  `json:"stock_mode"`
		ProductID    *int    `json:"product_id"`
		ProductName  string  `json:"product_name"`
		HasVariants  bool    `json:"has_variants"`
		IsActive     bool    `json:"is_active"`
		IsSoldOut    bool    `json:"is_sold_out"`
		PriceMin     float64 `json:"price_min"`
		PriceMax     float64 `json:"price_max"`
		VariantCount int     `json:"variant_count"`
	}

	RecipeLineResponse struct {
		IngredientID   int    `json:"ingredient_id"`
		IngredientName string `json:"ingredient_name"`
		Amount         string `json:"amount"`
	}

	VariantResponse struct {
		ID          int                  `json:"id"`
		SizeID      *int                 `json:"size_id"`
		SizeName    string               `json:"size_name"`
		Temperature *string              `json:"temperature"`
		Price       float64              `json:"price"`
		IsActive    bool                 `json:"is_active"`
		Recipe      []RecipeLineResponse `json:"recipe"`
	}

	ModifierGroupRef struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	MenuItemDetailResponse struct {
		MenuItemResponse
		Variants       []VariantResponse  `json:"variants"`
		ModifierGroups []ModifierGroupRef `json:"modifier_groups"`
	}
)
