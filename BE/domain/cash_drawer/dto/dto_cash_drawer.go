package dto

import "time"

type (
	// REQUEST
	GetHistoryRequest struct {
		Page      int    `json:"page" validate:"required,min=1"`
		Limit     int    `json:"limit" validate:"required,min=1"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		UserID    *int   `json:"user_id"`
		ShiftID   *int   `json:"shift_id"`
		Status    string `json:"status" validate:"max=20"`
		SortBy    string `json:"sort_by"`
		SortOrder string `json:"sort_order"`
	}

	GetByIDRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	OpenRequest struct {
		ShiftID        *int    `json:"shift_id"`
		OpeningBalance float64 `json:"opening_balance" validate:"min=0"`
		Notes          string  `json:"notes" validate:"max=500"`
	}

	CloseUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	CloseRequest struct {
		ID             int     `json:"-"`
		ClosingBalance float64 `json:"closing_balance" validate:"min=0"`
		Notes          string  `json:"notes" validate:"max=500"`
	}

	UpdateSalesUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateSalesRequest struct {
		ID             int     `json:"-"`
		TotalSales     float64 `json:"total_sales" validate:"min=0"`
		TotalCashSales float64 `json:"total_cash_sales" validate:"min=0"`
	}

	UpdateExpensesUriRequest struct {
		ID int `uri:"id" validate:"required,min=1"`
	}

	UpdateExpensesRequest struct {
		ID            int     `json:"-"`
		TotalExpenses float64 `json:"total_expenses" validate:"min=0"`
	}

	// RESPONSE
	CurrentCashDrawerResponse struct {
		ID              int       `json:"id"`
		UserID          int       `json:"user_id"`
		UserName        string    `json:"user_name"`
		ShiftID         *int      `json:"shift_id"`
		ShiftName       *string   `json:"shift_name"`
		ShiftStart      *string   `json:"shift_start"`
		ShiftEnd        *string   `json:"shift_end"`
		OpenTime        time.Time `json:"open_time"`
		OpeningBalance  float64   `json:"opening_balance"`
		TotalSales      float64   `json:"total_sales"`
		TotalCashSales  float64   `json:"total_cash_sales"`
		TotalExpenses   float64   `json:"total_expenses"`
		ExpectedBalance float64   `json:"expected_balance"`
		Status          string    `json:"status"`
		OpenNotes       *string   `json:"open_notes"`
	}

	CashDrawerHistoryResponse struct {
		ID              int        `json:"id"`
		UserName        string     `json:"user_name"`
		ShiftName       *string    `json:"shift_name"`
		OpenTime        time.Time  `json:"open_time"`
		CloseTime       *time.Time `json:"close_time"`
		OpeningBalance  float64    `json:"opening_balance"`
		ClosingBalance  *float64   `json:"closing_balance"`
		ExpectedBalance float64    `json:"expected_balance"`
		Difference      *float64   `json:"difference"`
		TotalSales      float64    `json:"total_sales"`
		TotalCashSales  float64    `json:"total_cash_sales"`
		TotalExpenses   float64    `json:"total_expenses"`
		Status          string     `json:"status"`
	}

	NonCashSaleItem struct {
		PaymentMethod string  `json:"payment_method"`
		Label         string  `json:"label"`
		Total         float64 `json:"total"`
	}

	CashDrawerTransaction struct {
		TransactionDate time.Time `json:"transaction_date"`
		TransactionCode string    `json:"transaction_code"`
		CustomerName    string    `json:"customer_name"`
		TotalAmount     float64   `json:"total_amount"`
	}

	NonCashTransaction struct {
		TransactionDate    time.Time `json:"transaction_date"`
		TransactionCode    string    `json:"transaction_code"`
		CustomerName       string    `json:"customer_name"`
		PaymentMethodLabel string    `json:"payment_method_label"`
		TotalAmount        float64   `json:"total_amount"`
	}

	CashDrawerExpenseItem struct {
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
	}

	CashDrawerDetailResponse struct {
		ID              int                     `json:"id"`
		UserID          int                     `json:"-"`
		CashierName     string                  `json:"cashier_name"`
		ShiftName       *string                 `json:"shift_name"`
		ShiftStart      *string                 `json:"shift_start"`
		ShiftEnd        *string                 `json:"shift_end"`
		OpenTime        time.Time               `json:"open_time"`
		CloseTime       *time.Time              `json:"close_time"`
		OpeningBalance  float64                 `json:"opening_balance"`
		ClosingBalance  *float64                `json:"closing_balance"`
		ExpectedBalance float64                 `json:"expected_balance"`
		TotalCashSales  float64                 `json:"total_cash_sales"`
		TotalExpenses   float64                 `json:"total_expenses"`
		Difference      *float64                `json:"difference"`
		Status          string                  `json:"status"`
		Notes           *string                 `json:"notes"`
		OpenNotes       *string                 `json:"open_notes"`
		Transactions    []CashDrawerTransaction `json:"transactions"`
		Expenses        []CashDrawerExpenseItem `json:"expenses"`
		NonCashSales    []NonCashSaleItem       `json:"non_cash_sales"`
	}

	MyCashResponse struct {
		ID                  *int                    `json:"id"`
		Status              string                  `json:"status"`
		ShiftName           *string                 `json:"shift_name"`
		ShiftStart          *string                 `json:"shift_start"`
		ShiftEnd            *string                 `json:"shift_end"`
		OpenTime            *time.Time              `json:"open_time"`
		OpeningBalance      float64                 `json:"opening_balance"`
		TotalCashSales      float64                 `json:"total_cash_sales"`
		TotalExpenses       float64                 `json:"total_expenses"`
		ExpectedBalance     float64                 `json:"expected_balance"`
		OpenNotes           *string                 `json:"open_notes"`
		Transactions        []CashDrawerTransaction `json:"transactions"`
		Expenses            []CashDrawerExpenseItem `json:"expenses"`
		NonCashSales        []NonCashSaleItem       `json:"non_cash_sales"`
		NonCashTransactions []NonCashTransaction    `json:"non_cash_transactions"`
	}

	CashDrawerSummaryResponse struct {
		TotalOpening  float64                      `json:"total_opening"`
		TotalClosing  float64                      `json:"total_closing"`
		TotalExpenses float64                      `json:"total_expenses"`
		Net           float64                      `json:"net"`
		Records       []*CashDrawerHistoryResponse `json:"records"`
	}

	OpenResponse struct {
		ID int `json:"id"`
	}

	CloseResponse struct {
		ExpectedBalance float64 `json:"expected_balance"`
		ClosingBalance  float64 `json:"closing_balance"`
		Difference      float64 `json:"difference"`
	}
)
