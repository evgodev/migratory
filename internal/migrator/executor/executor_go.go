package executor

import (
	"context"
	"database/sql"
)

// GoMigrateFn defines a function type for performing database migrations using a context and transaction.
type GoMigrateFn func(ctx context.Context, tx *sql.Tx) error

type GoExecutor struct {
	upFn, downFn GoMigrateFn
}

func NewGoExecutor(up, down GoMigrateFn) GoExecutor {
	return GoExecutor{
		upFn:   up,
		downFn: down,
	}
}

func (g GoExecutor) UpTx(ctx context.Context, tx *sql.Tx) error {
	return g.upFn(ctx, tx)
}

func (g GoExecutor) DownTx(ctx context.Context, tx *sql.Tx) error {
	return g.downFn(ctx, tx)
}
