package model

import "time"

type LogSchedulerModel struct {
	Id            string    `gorm:"primaryKey;column:id"`
	SchedulerName string    `gorm:"column:scheduler_name"`
	Status        string    `gorm:"column:status"`
	Message       *string   `gorm:"column:message"`
	DurationMs    *int64    `gorm:"column:duration_ms"`
	ExecutedAt    time.Time `gorm:"column:executed_at"`
}
