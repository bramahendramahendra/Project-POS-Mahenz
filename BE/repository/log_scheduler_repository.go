package repository

import (
	"pos_api/model"
	"time"

	"gorm.io/gorm"
)

const (
	insertLogScheduler = `INSERT INTO log_schedulers (id, scheduler_name, status, message, duration_ms, executed_at) VALUES (?, ?, ?, ?, ?, ?)`
)

var (
	LogSchedulerRepo LogSchedulerRepository
)

type LogSchedulerRepository interface {
	InsertLogScheduler(data *model.LogSchedulerModel) error
}

type logSchedulerRepository struct {
	Db *gorm.DB
}

func NewLogSchedulerRepository(db *gorm.DB) *logSchedulerRepository {
	return &logSchedulerRepository{Db: db}
}

func (r *logSchedulerRepository) InsertLogScheduler(data *model.LogSchedulerModel) error {
	if data.ExecutedAt.IsZero() {
		data.ExecutedAt = time.Now()
	}
	return r.Db.Exec(insertLogScheduler,
		data.Id,
		data.SchedulerName,
		data.Status,
		data.Message,
		data.DurationMs,
		data.ExecutedAt,
	).Error
}
