package service

import (
	repo_cash_drawer "pos_api/domain/cash_drawer/repo"
	service_cash_drawer "pos_api/domain/cash_drawer/service"
	repo_expense "pos_api/domain/expense/repo"
	repo_product "pos_api/domain/product/repo"
	"pos_api/domain/sync/dto"
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
		productRepo     repo_product.ProductRepoInterface
		cashDrawerRepo  repo_cash_drawer.CashDrawerRepoInterface
		cashDrawerSvc   service_cash_drawer.CashDrawerServiceInterface
	}
)

func NewSyncService(
	repo sync_repo.SyncRepoInterface,
	transactionRepo repo_transaction.TransactionRepoInterface,
	expenseRepo repo_expense.ExpenseRepoInterface,
	productRepo repo_product.ProductRepoInterface,
	cashDrawerRepo repo_cash_drawer.CashDrawerRepoInterface,
	cashDrawerSvc service_cash_drawer.CashDrawerServiceInterface,
) *syncService {
	return &syncService{
		repo:            repo,
		transactionRepo: transactionRepo,
		expenseRepo:     expenseRepo,
		productRepo:     productRepo,
		cashDrawerRepo:  cashDrawerRepo,
		cashDrawerSvc:   cashDrawerSvc,
	}
}
