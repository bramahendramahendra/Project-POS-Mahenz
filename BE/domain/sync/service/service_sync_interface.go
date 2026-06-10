package service

import (
	"pos_api/domain/sync/dto"
	repo_expense "pos_api/domain/expense/repo"
	sync_repo "pos_api/domain/sync/repo"
	repo_transaction "pos_api/domain/transaction/repo"
)

type (
	SyncServiceInterface interface {
		PushSync(req *dto.PushSyncRequest) (*dto.PushSyncResponse, error)
		GetConflicts(filter *dto.ConflictFilter) (*dto.ConflictListResponse, error)
		CountPendingConflicts() (int, error)
		ResolveConflict(id, userID int, resolution string) error
		GetQueue(filter *dto.QueueFilter) (*dto.QueueListResponse, error)
		GetHistory(filter *dto.HistoryFilter) (*dto.SyncHistoryListResponse, error)
	}

	syncService struct {
		repo            sync_repo.SyncRepoInterface
		transactionRepo repo_transaction.TransactionRepoInterface
		expenseRepo     repo_expense.ExpenseRepoInterface
	}
)

func NewSyncService(
	repo sync_repo.SyncRepoInterface,
	transactionRepo repo_transaction.TransactionRepoInterface,
	expenseRepo repo_expense.ExpenseRepoInterface,
) *syncService {
	return &syncService{
		repo:            repo,
		transactionRepo: transactionRepo,
		expenseRepo:     expenseRepo,
	}
}
