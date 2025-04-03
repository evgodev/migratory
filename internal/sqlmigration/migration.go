package sqlmigration

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/korfairo/migratory/internal/migrator"
	"github.com/korfairo/migratory/internal/parser"
)

type sqlExecutor struct {
	statements statements
}

type statements struct {
	up, down []string
}

func newSQLExecutor(up, down []string) sqlExecutor {
	return sqlExecutor{
		statements: statements{
			up:   up,
			down: down,
		},
	}
}

func (s sqlExecutor) UpTx(ctx context.Context, tx *sql.Tx) error {
	return execute(ctx, tx, s.statements.up)
}

func (s sqlExecutor) DownTx(ctx context.Context, tx *sql.Tx) error {
	return execute(ctx, tx, s.statements.down)
}

type sqlExecutorNoTx struct {
	statements statements
}

func newSQLExecutorNoTx(up, down []string) sqlExecutorNoTx {
	return sqlExecutorNoTx{
		statements: statements{
			up:   up,
			down: down,
		},
	}
}

func (s sqlExecutorNoTx) Up(ctx context.Context, db *sql.DB) error {
	return execute(ctx, db, s.statements.up)
}

func (s sqlExecutorNoTx) Down(ctx context.Context, db *sql.DB) error {
	return execute(ctx, db, s.statements.down)
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

type sqlPreparer struct {
	sourcePath string
	fsys       fs.FS
}

func newSQLPreparer(sourceFilePath string, fsys fs.FS) sqlPreparer {
	return sqlPreparer{
		sourcePath: sourceFilePath,
		fsys:       fsys,
	}
}

func (s sqlPreparer) Prepare() (*migrator.ExecutorContainer, error) {
	file, err := s.fsys.Open(s.sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at path %s: %w", s.sourcePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	parsed, err := parser.ParseMigration(file)
	if parsed == nil || err != nil {
		return nil, fmt.Errorf("failed to parse migration %s: %w", s.sourcePath, err)
	}

	var container *migrator.ExecutorContainer
	if parsed.DisableTransactionUp || parsed.DisableTransactionDown {
		executor := newSQLExecutorNoTx(parsed.UpStatements, parsed.DownStatements)
		container = migrator.NewExecutorDBContainer(executor)
	} else {
		executor := newSQLExecutor(parsed.UpStatements, parsed.DownStatements)
		container = migrator.NewExecutorTxContainer(executor)
	}

	return container, nil
}
