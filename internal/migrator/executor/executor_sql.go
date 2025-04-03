package executor

import (
	"context"
	"database/sql"
)

type SQLExecutor struct {
	statements statements
}

type statements struct {
	up, down []string
}

func NewSQLExecutor(up, down []string) SQLExecutor {
	return SQLExecutor{
		statements: statements{
			up:   up,
			down: down,
		},
	}
}

func (s SQLExecutor) UpTx(ctx context.Context, tx *sql.Tx) error {
	return execute(ctx, tx, s.statements.up)
}

func (s SQLExecutor) DownTx(ctx context.Context, tx *sql.Tx) error {
	return execute(ctx, tx, s.statements.down)
}
