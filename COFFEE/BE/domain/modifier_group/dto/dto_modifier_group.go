package dto

const (
	SelectionSingle   = "single"   // pilih satu
	SelectionMultiple = "multiple" // boleh banyak
)

type (
	// REQUEST
	GetAllRequest struct {
		Search   string `json:"search" validate:"max=100"`
		IsActive *bool  `json:"is_active"`
	}

	IDUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// OptionInput: id 0 = pilihan baru. Pilihan lama yang tidak dikirim lagi dihapus
	// (pesanan kelak menyimpan salinan nama dan harga, keputusan tahap 4.7).
	OptionInput struct {
		ID        int     `json:"id" validate:"gte=0"`
		Name      string  `json:"name" validate:"required,max=100,notblank"`
		Price     float64 `json:"price" validate:"gte=0,lte=999999999"`
		IsDefault bool    `json:"is_default"`
		IsActive  bool    `json:"is_active"`
	}

	SaveRequest struct {
		ID            int           `json:"-"`
		Name          string        `json:"name" validate:"required,max=100,notblank"`
		SelectionType string        `json:"selection_type" validate:"required,oneof=single multiple"`
		IsRequired    bool          `json:"is_required"`
		Options       []OptionInput `json:"options" validate:"required,min=1,max=50,dive"`
	}

	// RESPONSE
	OptionResponse struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		IsDefault bool    `json:"is_default"`
		IsActive  bool    `json:"is_active"`
	}

	GroupResponse struct {
		ID            int              `json:"id"`
		Name          string           `json:"name"`
		SelectionType string           `json:"selection_type"`
		IsRequired    bool             `json:"is_required"`
		IsActive      bool             `json:"is_active"`
		MenuCount     int              `json:"menu_count"`
		Options       []OptionResponse `json:"options"`
	}
)
