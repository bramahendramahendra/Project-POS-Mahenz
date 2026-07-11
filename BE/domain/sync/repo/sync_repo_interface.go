package repo

import (
	"pos_api/domain/sync/dto"
	"pos_api/domain/sync/model"

	"gorm.io/gorm"
)

type (
	SyncRepoInterface interface {
		GetEntitySnapshot(entityType string, serverID int) (*model.EntitySnapshot, error)
		GetRawByEntityAndID(entityType string, serverID int) (string, error)

		GetConflicts(filter *dto.ConflictFilter) ([]dto.ConflictResponse, int, error)
		GetConflictByID(id int) (*model.SyncConflict, error)
		CountPendingConflicts() (int, error)
		ResolveConflict(id, userID int, action string) error
		MarkResolved(id, resolvedBy int, action string) error
		CreateConflict(deviceID string, item *dto.SyncItem, onlineData string) (int, error)

		GetQueue(filter *dto.QueueFilter) ([]dto.QueueResponse, int, error)
		CreateQueueItem(deviceID string, item *dto.SyncItem, status string) (int, error)
		UpdateQueueStatus(id int, status, errMsg string) error

		InsertHistory(h model.SyncHistory) error
		GetHistory(filter *dto.HistoryFilter) ([]dto.SyncHistoryResponse, int, error)
	}

	syncRepo struct {
		db *gorm.DB
	}
)

func NewSyncRepo(db *gorm.DB) *syncRepo {
	return &syncRepo{db: db}
}
