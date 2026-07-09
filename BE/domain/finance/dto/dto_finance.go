package dto

type (
	// REQUEST
	GetSummaryRequest struct {
		DateFrom string `json:"date_from"`
		DateTo   string `json:"date_to"`
	}

	GetCashflowRequest struct {
		DateFrom string `json:"date_from"`
		DateTo   string `json:"date_to"`
		Type     string `json:"type"`
		Page     int    `json:"page" validate:"required,min=1"`
		Limit    int    `json:"limit" validate:"required,min=1"`
	}

	// RESPONSE
	SummaryResponse struct {
		TotalIncome     float64 `json:"total_income"`
		TotalExpense    float64 `json:"total_expense"`
		NetProfit       float64 `json:"net_profit"`
		TotalReceivable float64 `json:"total_receivable"`
		PeriodLabel     string  `json:"period_label"`
	}

	CashflowItemResponse struct {
		ID          int     `json:"id"`
		Type        string  `json:"type"`
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Date        string  `json:"date"`
	}
)
