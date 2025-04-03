package migrator

import (
	"context"
	"database/sql"
)

// Executors encapsulates execution logic for database migrations, supporting both transactional and non-transactional modes.
// It holds either an ExecutorTx or an ExecutorDB to execute migrations based on the execution context.
type Executors struct {
	useDB      bool
	executorTx ExecutorTx
	executorDB ExecutorDB
}

type ExecutorTx interface {
	UpTx(ctx context.Context, tx *sql.Tx) error
	DownTx(ctx context.Context, tx *sql.Tx) error
}

type ExecutorDB interface {
	Up(ctx context.Context, db *sql.DB) error
	Down(ctx context.Context, db *sql.DB) error
}

func newExecutorTxContainer(executorTx ExecutorTx) *Executors {
	return &Executors{
		useDB:      false,
		executorTx: executorTx,
	}
}

func newExecutorDBContainer(executorDB ExecutorDB) *Executors {
	return &Executors{
		useDB:      true,
		executorDB: executorDB,
	}
}

func (e Executors) ExecutorTx() ExecutorTx {
	return e.executorTx
}

func (e Executors) ExecutorDB() ExecutorDB {
	return e.executorDB
}
