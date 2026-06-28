package dto

type SettingKeyValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StoreProfileRequest struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	LogoURL    string  `json:"logo_url"`
	TaxDefault float64 `json:"tax_default"`
}

type StoreProfileResponse struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	LogoURL    string  `json:"logo_url"`
	TaxDefault float64 `json:"tax_default"`
}

type PrinterSettingsRequest struct {
	PaperSize     string `json:"paper_size"`
	ReceiptHeader string `json:"receipt_header"`
	ReceiptFooter string `json:"receipt_footer"`
	ShowLogo      bool   `json:"show_logo"`
	AutoPrint     bool   `json:"auto_print"`
}

type PrinterSettingsResponse struct {
	PaperSize     string `json:"paper_size"`
	ReceiptHeader string `json:"receipt_header"`
	ReceiptFooter string `json:"receipt_footer"`
	ShowLogo      bool   `json:"show_logo"`
	AutoPrint     bool   `json:"auto_print"`
}
