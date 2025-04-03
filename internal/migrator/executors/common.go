package executors

import (
	"context"
	"database/sql"
)

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
