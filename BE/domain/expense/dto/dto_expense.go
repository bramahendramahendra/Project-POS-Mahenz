package dto

type (
	// REQUEST
	GetAllRequest struct {
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Category  string `json:"category" validate:"max=100"`
		UserID    *int   `json:"user_id"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CreateRequest struct {
		ExpenseDate   string  `json:"expense_date" validate:"required"`
		Category      string  `json:"category" validate:"required,max=100"`
		Description   string  `json:"description" validate:"max=255"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash transfer card qris kredit"`
		Notes         string  `json:"notes" validate:"max=500"`
		CashDrawerID  *int    `json:"cash_drawer_id"`
	}

	UpdateUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateRequest struct {
		ID            int     `json:"-"`
		ExpenseDate   string  `json:"expense_date" validate:"required"`
		Category      string  `json:"category" validate:"required,max=100"`
		Description   string  `json:"description" validate:"max=255"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash transfer card qris kredit"`
		Notes         string  `json:"notes" validate:"max=500"`
	}

	DeleteRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	// RESPONSE
	ExpenseResponse struct {
		ID            int     `json:"id"`
		ExpenseDate   string  `json:"expense_date"`
		Category      string  `json:"category"`
		Description   string  `json:"description"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		UserName      string  `json:"user_name"`
		Notes         string  `json:"notes"`
		UserID        int     `json:"-"`
	}
)
