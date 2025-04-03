package executors

import (
	"context"
	"database/sql"
)

type SQLExecutorNoTx struct {
	statements statements
}

func NewSQLExecutorNoTx(up, down []string) SQLExecutorNoTx {
	return SQLExecutorNoTx{
		statements: statements{
			up:   up,
			down: down,
		},
	}
}

func (s SQLExecutorNoTx) Up(ctx context.Context, db *sql.DB) error {
	return execute(ctx, db, s.statements.up)
}

func (s SQLExecutorNoTx) Down(ctx context.Context, db *sql.DB) error {
	return execute(ctx, db, s.statements.down)
}
