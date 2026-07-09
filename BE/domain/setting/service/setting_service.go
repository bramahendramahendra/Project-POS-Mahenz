package service

import (
	"fmt"
	"pos_api/domain/setting/dto"
	"pos_api/errors"
	"strconv"
)

func (s *settingService) GetAll() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, item := range settings {
		result[item.Key] = item.Value
	}
	return result, nil
}

func (s *settingService) GetByKey(key string) (string, error) {
	setting, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", &errors.NotFoundError{Message: "Setting tidak ditemukan"}
	}
	return setting.Value, nil
}

func (s *settingService) Save(data map[string]string) error {
	for key, value := range data {
		if err := s.repo.Upsert(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *settingService) Reset() error {
	return s.repo.ResetAll()
}

func (s *settingService) getSettingValue(key string) string {
	setting, err := s.repo.GetByKey(key)
	if err != nil || setting == nil {
		return ""
	}
	return setting.Value
}

func (s *settingService) GetStoreProfile() (*dto.StoreProfileResponse, error) {
	taxStr := s.getSettingValue("tax_default")
	taxDefault, _ := strconv.ParseFloat(taxStr, 64)
	return &dto.StoreProfileResponse{
		Name:       s.getSettingValue("store_name"),
		Address:    s.getSettingValue("store_address"),
		Phone:      s.getSettingValue("store_phone"),
		Email:      s.getSettingValue("store_email"),
		LogoURL:    s.getSettingValue("store_logo_url"),
		TaxDefault: taxDefault,
	}, nil
}

func (s *settingService) UpdateStoreProfile(req *dto.StoreProfileRequest) error {
	data := map[string]string{
		"store_name":     req.Name,
		"store_address":  req.Address,
		"store_phone":    req.Phone,
		"store_email":    req.Email,
		"store_logo_url": req.LogoURL,
		"tax_default":    fmt.Sprintf("%g", req.TaxDefault),
	}
	return s.Save(data)
}

func (s *settingService) GetPrinterSettings() (*dto.PrinterSettingsResponse, error) {
	showLogo := s.getSettingValue("printer_show_logo") == "true"
	autoPrint := s.getSettingValue("printer_auto_print") == "true"
	return &dto.PrinterSettingsResponse{
		PaperSize:     s.getSettingValue("printer_paper_size"),
		ReceiptHeader: s.getSettingValue("printer_receipt_header"),
		ReceiptFooter: s.getSettingValue("printer_receipt_footer"),
		ShowLogo:      showLogo,
		AutoPrint:     autoPrint,
	}, nil
}

func (s *settingService) UpdatePrinterSettings(req *dto.PrinterSettingsRequest) error {
	showLogo := "false"
	if req.ShowLogo {
		showLogo = "true"
	}
	autoPrint := "false"
	if req.AutoPrint {
		autoPrint = "true"
	}
	data := map[string]string{
		"printer_paper_size":     req.PaperSize,
		"printer_receipt_header": req.ReceiptHeader,
		"printer_receipt_footer": req.ReceiptFooter,
		"printer_show_logo":      showLogo,
		"printer_auto_print":     autoPrint,
	}
	return s.Save(data)
}
