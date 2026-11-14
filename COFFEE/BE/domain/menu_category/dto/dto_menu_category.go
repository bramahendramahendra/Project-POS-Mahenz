package dto

type (
	// REQUEST
	GetAllRequest struct {
		Search   string `json:"search" validate:"max=100"`
		IsActive *bool  `json:"is_active"`
	}

	IDUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	SaveRequest struct {
		ID        int    `json:"-"`
		Name      string `json:"name" validate:"required,max=100,notblank"`
		SortOrder int    `json:"sort_order" validate:"gte=0,lte=9999"`
	}

	// RESPONSE
	MenuCategoryResponse struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
		ItemCount int    `json:"item_count"`
	}

	OptionResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
