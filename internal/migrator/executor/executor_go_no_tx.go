package executor

import (
	"context"
	"database/sql"
)

// GoMigrateNoTxFn defines a function type for non-transactional database migrations
// using a context and a SQL database connection.
type GoMigrateNoTxFn func(ctx context.Context, db *sql.DB) error

type GoExecutorNoTx struct {
	upFn, downFn GoMigrateNoTxFn
}

func NewGoExecutorNoTx(upNoTx, downNoTx GoMigrateNoTxFn) GoExecutorNoTx {
	return GoExecutorNoTx{
		upFn:   upNoTx,
		downFn: downNoTx,
	}
}

func (g GoExecutorNoTx) Up(ctx context.Context, db *sql.DB) error {
	if err := g.upFn(ctx, db); err != nil {
		return err
	}
	return nil
}

func (g GoExecutorNoTx) Down(ctx context.Context, db *sql.DB) error {
	return g.downFn(ctx, db)
}
