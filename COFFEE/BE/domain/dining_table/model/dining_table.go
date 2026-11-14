package model

// DiningTable: meja kedai (mis. "Meja 1", "Teras 2"), dikelola admin/owner (keputusan 56, 97, 98).
// UsageCount = jumlah pesanan yang memakai meja (mulai tahap 4.7).
type DiningTable struct {
	ID         int    `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	SortOrder  int    `gorm:"column:sort_order"`
	IsActive   bool   `gorm:"column:is_active"`
	UsageCount int    `gorm:"column:usage_count"`
	// Pesanan terbuka yang memakai meja (hanya diisi query pilihan meja, tahap 4.8).
	OpenBillID     *int   `gorm:"column:open_bill_id"`
	OpenQueueLabel string `gorm:"column:open_queue_label"`
}
