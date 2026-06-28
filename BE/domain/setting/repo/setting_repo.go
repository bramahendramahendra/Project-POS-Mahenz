package repo

import "pos_api/domain/setting/model"

const (
	getAllSettingsQuery   = `SELECT setting_key, setting_value FROM settings ORDER BY setting_key`
	getSettingByKeyQuery = `SELECT setting_key, setting_value FROM settings WHERE setting_key = ?`
	upsertSettingQuery   = `INSERT INTO settings (setting_key, setting_value) VALUES (?, ?) ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_at = NOW()`
	resetSettingsQuery   = `UPDATE settings SET setting_value = CASE setting_key WHEN 'store_name' THEN 'Toko Retail' WHEN 'store_address' THEN '' WHEN 'store_phone' THEN '' WHEN 'store_email' THEN '' WHEN 'store_logo_url' THEN '' WHEN 'tax_default' THEN '0' WHEN 'tax_enabled' THEN '0' WHEN 'tax_percent' THEN '11' WHEN 'receipt_footer' THEN 'Terima kasih telah berbelanja' WHEN 'stock_notification_enabled' THEN '1' WHEN 'pagination_sizes' THEN '[10,20,50]' WHEN 'printer_paper_size' THEN '80mm' WHEN 'printer_receipt_header' THEN '' WHEN 'printer_receipt_footer' THEN 'Terima kasih telah berbelanja' WHEN 'printer_show_logo' THEN 'false' WHEN 'printer_auto_print' THEN 'false' ELSE '' END`
)

func (r *settingRepo) GetAll() ([]model.Setting, error) {
	var settings []model.Setting
	rows, err := r.db.Raw(getAllSettingsQuery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s model.Setting
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, rows.Err()
}

func (r *settingRepo) GetByKey(key string) (*model.Setting, error) {
	var s model.Setting
	rows, err := r.db.Raw(getSettingByKeyQuery, key).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		return &s, nil
	}
	return nil, nil
}

func (r *settingRepo) Upsert(key, value string) error {
	return r.db.Exec(upsertSettingQuery, key, value).Error
}

func (r *settingRepo) ResetAll() error {
	return r.db.Exec(resetSettingsQuery).Error
}
