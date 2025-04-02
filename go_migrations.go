package migratory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"

	"github.com/korfairo/migratory/internal/migrator"
	"github.com/korfairo/migratory/internal/sqlmigration"
)

// GoMigrateFn defines a function type for performing database migrations using a context and transaction.
type GoMigrateFn func(ctx context.Context, tx *sql.Tx) error

// AddMigration registers a new migration with `up` and `down` functions for handling database schema changes.
func AddMigration(up, down GoMigrateFn) {
	_, fileName, _, _ := runtime.Caller(1) //nolint:dogsled
	executor := newGoExecutor(up, down)
	addGoMigration(fileName, executor)
}

// GoMigrateNoTxFn defines a function type for non-transactional database migrations
// using a context and a SQL database connection.
type GoMigrateNoTxFn func(ctx context.Context, db *sql.DB) error

// AddMigrationNoTx registers a database migration function pair that operates without transactions.
func AddMigrationNoTx(up, down GoMigrateNoTxFn) {
	_, fileName, _, _ := runtime.Caller(1) //nolint:dogsled
	executorNoTx := newGoExecutorNoTx(up, down)
	addGoMigrationNoTx(fileName, executorNoTx)
}

var globalGoMigrations []goMigration

func addGoMigration(fileName string, executor goExecutor) {
	globalGoMigrations = append(globalGoMigrations, goMigration{
		sourceName: fileName,
		noTx:       false,
		executor:   executor,
	})
}

func addGoMigrationNoTx(fileName string, executorNoTx goExecutorNoTx) {
	globalGoMigrations = append(globalGoMigrations, goMigration{
		sourceName:   fileName,
		noTx:         true,
		executorNoTx: executorNoTx,
	})
}

type goMigration struct {
	sourceName string
	noTx       bool

	executor     goExecutor
	executorNoTx goExecutorNoTx
}

type goExecutor struct {
	upFn, downFn GoMigrateFn
}

func newGoExecutor(up, down GoMigrateFn) goExecutor {
	return goExecutor{
		upFn:   up,
		downFn: down,
	}
}

func (g goExecutor) UpTx(ctx context.Context, tx *sql.Tx) error {
	return g.upFn(ctx, tx)
}

func (g goExecutor) DownTx(ctx context.Context, tx *sql.Tx) error {
	return g.downFn(ctx, tx)
}

type goExecutorNoTx struct {
	upFn, downFn GoMigrateNoTxFn
}

func newGoExecutorNoTx(upNoTx, downNoTx GoMigrateNoTxFn) goExecutorNoTx {
	return goExecutorNoTx{
		upFn:   upNoTx,
		downFn: downNoTx,
	}
}

func (g goExecutorNoTx) Up(ctx context.Context, db *sql.DB) error {
	if err := g.upFn(ctx, db); err != nil {
		return err
	}
	return nil
}

func (g goExecutorNoTx) Down(ctx context.Context, db *sql.DB) error {
	return g.downFn(ctx, db)
}

func registerGoMigrations(goMigrations []goMigration) (migrator.Migrations, error) {
	goMigrationsCount := len(goMigrations)
	if goMigrationsCount == 0 {
		return nil, errors.New("no migrations were added")
	}

	var migrations migrator.Migrations
	uniqueIDMap := make(map[int64]struct{}, goMigrationsCount)
	for _, m := range goMigrations {
		id, name, err := sqlmigration.ParseMigrationFileName(m.sourceName)
		if err != nil {
			return nil, err
		}

		if _, exists := uniqueIDMap[id]; exists {
			return nil, fmt.Errorf("migration id '%d' is not unique", id)
		}
		uniqueIDMap[id] = struct{}{}

		if m.noTx {
			migrations = append(migrations, migrator.NewMigrationNoTx(id, name, m.executorNoTx))
			continue
		}

		migrations = append(migrations, migrator.NewMigration(id, name, m.executor))
	}

	return migrations, nil
}
