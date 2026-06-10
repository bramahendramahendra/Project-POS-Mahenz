package dto

type PaymentMethodResponse struct {
	ID        int    `json:"id"         gorm:"column:id"`
	Code      string `json:"code"       gorm:"column:code"`
	Label     string `json:"label"      gorm:"column:label"`
	IsActive  int    `json:"is_active"  gorm:"column:is_active"`
	SortOrder int    `json:"sort_order" gorm:"column:sort_order"`
}
