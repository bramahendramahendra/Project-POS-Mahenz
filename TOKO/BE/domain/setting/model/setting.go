package model

type Setting struct {
	Key   string `gorm:"column:setting_key"`
	Value string `gorm:"column:setting_value"`
}
