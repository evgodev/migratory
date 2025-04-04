package migratory

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/evgodev/migratory/internal/migrator"
	"github.com/evgodev/migratory/internal/migrator/executor"
)

var goMigrations []goMigration

type goMigration struct {
	fileName string
	noTx     bool // Defines what to use: up and down or upNoTx and downNoTx.

	up   GoMigrateFn
	down GoMigrateFn

	upNoTx   GoMigrateNoTxFn
	downNoTx GoMigrateNoTxFn
}

// GoMigrateFn defines a function type for performing database migrations using a context and transaction.
type GoMigrateFn func(ctx context.Context, tx *sql.Tx) error

// GoMigrateNoTxFn defines a function type for non-transactional database migrations
// using a context and a SQL database connection.
type GoMigrateNoTxFn func(ctx context.Context, db *sql.DB) error

// AddMigration registers a new migration with `up` and `down` functions for handling database schema changes.
func AddMigration(up, down GoMigrateFn) {
	_, fileName, _, _ := runtime.Caller(1) //nolint:dogsled
	goMigrations = append(goMigrations, goMigration{
		fileName: fileName,
		noTx:     false,
		up:       up,
		down:     down,
	})
}

// AddMigrationNoTx registers a database migration function pair that operates without transactions.
func AddMigrationNoTx(up, down GoMigrateNoTxFn) {
	_, fileName, _, _ := runtime.Caller(1) //nolint:dogsled
	goMigrations = append(goMigrations, goMigration{
		fileName: fileName,
		noTx:     true,
		upNoTx:   up,
		downNoTx: down,
	})
}

func convertGoMigrations() (migrator.Migrations, error) {
	result := make(migrator.Migrations, 0, len(goMigrations))
	for _, m := range goMigrations {
		id, name, err := migrator.ParseMigrationFileName(m.fileName)
		if err != nil {
			return nil, ErrIncorrectMigrationName
		}

		var converted migrator.Migration
		if m.noTx {
			converted = migrator.NewGoMigrationNoTx(id, name, convertNoTxFn(m.upNoTx), convertNoTxFn(m.downNoTx))
		} else {
			converted = migrator.NewGoMigration(id, name, convertFn(m.up), convertFn(m.down))
		}

		result = append(result, converted)
	}

	return result, nil
}

func convertFn(fn GoMigrateFn) executor.GoMigrateFn {
	return func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, tx)
	}
}

func convertNoTxFn(fn GoMigrateNoTxFn) executor.GoMigrateNoTxFn {
	return func(ctx context.Context, db *sql.DB) error {
		return fn(ctx, db)
	}
}
