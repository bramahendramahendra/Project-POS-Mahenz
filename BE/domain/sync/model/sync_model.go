package model

import "time"

type SyncHistory struct {
	ID            int64      `gorm:"column:id"`
	DeviceID      string     `gorm:"column:device_id"`
	DeviceType    string     `gorm:"column:device_type"`
	TotalItems    int        `gorm:"column:total_items"`
	SyncedItems   int        `gorm:"column:synced_items"`
	ConflictItems int        `gorm:"column:conflict_items"`
	FailedItems   int        `gorm:"column:failed_items"`
	PendingItems  int        `gorm:"column:pending_items"`
	DurationMs    int        `gorm:"column:duration_ms"`
	Status        string     `gorm:"column:status"`
	StartedAt     time.Time  `gorm:"column:started_at"`
	FinishedAt    *time.Time `gorm:"column:finished_at"`
}

type SyncConflict struct {
	ID             int        `json:"id"`
	EntityType     string     `json:"entity_type"`
	EntityID       int        `json:"entity_id"`
	LocalID        string     `json:"local_id"`
	DeviceID       string     `json:"device_id"`
	DesktopData    string     `json:"desktop_data"`
	OnlineData     string     `json:"online_data"`
	DesktopTime    time.Time  `json:"desktop_time"`
	OnlineTime     time.Time  `json:"online_time"`
	Status         string     `json:"status"`
	ResolvedBy     *int       `json:"resolved_by"`
	Resolution     *string    `json:"resolution"`
	ResolvedAction *string    `json:"resolved_action"`
	ResolvedAt     *time.Time `json:"resolved_at"`
}

type EntitySnapshot struct {
	UpdatedAt string `json:"updated_at"`
	Data      string `json:"data"`
}

type SyncQueue struct {
	ID           int        `json:"id"`
	DeviceID     string     `json:"device_id"`
	EntityType   string     `json:"entity_type"`
	EntityID     int        `json:"entity_id"`
	Action       string     `json:"action"`
	Payload      string     `json:"payload"`
	Status       string     `json:"status"`
	RetryCount   int        `json:"retry_count"`
	ErrorMessage *string    `json:"error_message"`
	SyncedAt     *time.Time `json:"synced_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
