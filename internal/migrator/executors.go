package migrator

import (
	"context"
	"database/sql"
)

type ExecutorTx interface {
	UpTx(ctx context.Context, tx *sql.Tx) error
	DownTx(ctx context.Context, tx *sql.Tx) error
}

type ExecutorDB interface {
	Up(ctx context.Context, db *sql.DB) error
	Down(ctx context.Context, db *sql.DB) error
}

// executors encapsulates execution logic for database migrations,
// supporting both transactional and non-transactional modes.
// It holds either an ExecutorTx or an ExecutorDB to execute migrations based on the execution context.
type executors struct {
	useDB      bool
	executorTx ExecutorTx
	executorDB ExecutorDB
}

func newExecutorTxContainer(executorTx ExecutorTx) *executors {
	return &executors{
		useDB:      false,
		executorTx: executorTx,
	}
}

func newExecutorDBContainer(executorDB ExecutorDB) *executors {
	return &executors{
		useDB:      true,
		executorDB: executorDB,
	}
}

func (e executors) ExecutorTx() ExecutorTx {
	return e.executorTx
}

func (e executors) ExecutorDB() ExecutorDB {
	return e.executorDB
}
