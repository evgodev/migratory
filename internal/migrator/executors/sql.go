package executors

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

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func execute(ctx context.Context, executor QueryExecutor, statements []string) error {
	for _, query := range statements {
		_, err := executor.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}
